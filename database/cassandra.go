package database

import (
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/gemnasium/migrate/driver/cassandra"
	"github.com/gemnasium/migrate/migrate"
	"github.com/gocql/gocql"
	"github.com/relops/cqlr"
	uuid "github.com/satori/go.uuid"
)

type emailIDMap struct {
	Email string `cql:"email"`
	ID    string `cql:"id"`
}

type Cassandra struct {
	Database

	keyspace string
	addrs    []string

	Session *gocql.Session
}

var Consistency *gocql.Consistency

func (c *Cassandra) Connect(keyspace string, addrs []string) error {
	c.keyspace = keyspace
	c.addrs = addrs

	cluster := gocql.NewCluster(addrs...)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	c.Session = session

	return nil
}

func (c *Cassandra) Migrate() error {
	wd := os.ExpandEnv("$GOPATH/src/github.com/lileio/account_service")
	errs, ok := migrate.UpSync("cassandra://"+c.addrs[0]+"/"+c.keyspace+"?disable_init_host_lookup", wd+"/migrations/cassandra")
	if !ok {
		fmt.Printf("migrations failed: %+v\n", errs)
		return errors.New("migration error")
	}

	return nil
}

func (c *Cassandra) Close() error {
	c.Session.Close()
	return nil
}

func (c *Cassandra) Truncate() error {
	return nil
}

func (c *Cassandra) List(
	count int32, token string) (
	accounts []*Account, next_token string, err error) {
	if token == "" {
		token = "0"
	}

	rows := c.Session.Query(
		`select id, name, email, images, createdat from accounts_map_id
		 where token(id) > token(?) limit ?`,
		token, count,
	).Iter().Scanner()

	if rows == nil {
		return accounts, next_token, err
	}

	for rows.Next() {
		a := &Account{}
		err := rows.Scan(&a.ID, &a.Name, &a.Email, &a.Images, &a.CreatedAt)
		if err != nil {
			return accounts, next_token, err
		}
		accounts = append(accounts, a)
	}

	if err = rows.Err(); err != nil {
		return accounts, next_token, err
	}

	if len(accounts) > 0 && len(accounts) == int(count) {
		next_token = accounts[count-1].ID
	}

	return accounts, next_token, err
}

func (c *Cassandra) Create(a *Account, password string) error {
	// Check the email table to see if this account exists
	ae, err := c.ReadByEmail(a.Email)
	if err != nil && err != ErrAccountNotFound {
		return err
	}

	if ae != nil {
		return ErrEmailExists
	}

	a.ID = uuid.NewV1().String()
	a.CreatedAt = time.Now()

	err = a.Valid()
	if err != nil {
		return err
	}

	if password == "" {
		return ErrNoPasswordGiven
	}

	err = a.hashPassword(password)
	if err != nil {
		return err
	}

	q := `INSERT INTO accounts_map_id
	(id, name, email, hashedpassword, images, createdat)
	VALUES (?, ?, ?, ?, ?, ?)`
	err = c.Session.Query(
		q,
		a.ID,
		a.Name,
		a.Email,
		a.HashedPassword,
		a.Images,
		a.CreatedAt).Exec()

	if err != nil {
		return err
	}

	return c.createEmailRow(a)
}

func (c *Cassandra) ReadByEmail(email string) (*Account, error) {
	e := &emailIDMap{}
	cql := `select email,id from accounts_map_email where email = ?`
	b := cqlr.BindQuery(c.Session.Query(cql, email))
	b.Scan(e)
	err := b.Close()

	if e.ID == "" {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	a, err := c.ReadByID(e.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (c *Cassandra) ReadByID(ID string) (*Account, error) {
	a := Account{}
	q := c.Session.Query(`select * from accounts_map_id where id = ?`, ID)
	b := cqlr.BindQuery(q)
	b.Scan(&a)
	err := b.Close()

	if a.ID == "" {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (c *Cassandra) Update(a *Account) error {
	current, err := c.ReadByID(a.ID)
	if err != nil {
		return err
	}

	if current.Email != a.Email {
		err = c.deleteEmailRow(a)
		if err != nil {
			return err
		}

		err = c.createEmailRow(a)
		if err != nil {
			return nil
		}
	}

	b := cqlr.Bind(`
		update accounts_map_id set
		name = ?, email = ?, images = ? where id = ?`, a)
	return b.Exec(c.Session)
}

func (c *Cassandra) Delete(ID string) error {
	a, err := c.ReadByID(ID)
	if err != nil {
		return err
	}

	b := cqlr.Bind(`delete from accounts_map_id where id = ?`, a)
	err = b.Exec(c.Session)
	if err != nil {
		return err
	}

	err = c.deleteEmailRow(a)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cassandra) createEmailRow(a *Account) error {
	cql := `insert into accounts_map_email (email, id) values (?, ?)`
	b := cqlr.Bind(cql, a)
	return b.Exec(c.Session)
}

func (c *Cassandra) deleteEmailRow(a *Account) error {
	b := cqlr.Bind(`delete from accounts_map_email where email = ?`, a)
	return b.Exec(c.Session)
}
