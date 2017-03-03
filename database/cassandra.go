package database

import (
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/gemnasium/migrate/driver/cassandra"
	"github.com/gemnasium/migrate/migrate"
	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
	uuid "github.com/satori/go.uuid"
)

type emailIDMap struct {
	Email string
	ID    string
}

type Cassandra struct {
	Database

	keyspace string
	addrs    []string

	Session               *gocql.Session
	AccountsIDSimpleTable gocassa.Table
	AccountsIDTable       gocassa.MapTable
	AccountsEmailTable    gocassa.MapTable
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

	wd := os.ExpandEnv("$GOPATH/src/github.com/lileio/account_service")
	allErrors, ok := migrate.UpSync("cassandra://"+addrs[0]+"/"+keyspace+"?disable_init_host_lookup", wd+"/migrations/cassandra")
	if !ok {
		fmt.Printf("allErrors = %+v\n", allErrors)
		return errors.New("migration error")
	}

	conn := gocassa.NewConnection(gocassa.GoCQLSessionToQueryExecutor(session))

	opts := gocassa.Options{
		Consistency: Consistency,
	}

	c.AccountsIDTable = conn.KeySpace(keyspace).MapTable(
		"accounts",
		"ID",
		&Account{},
	).WithOptions(opts)

	c.AccountsIDSimpleTable = conn.KeySpace(keyspace).Table("accounts", &Account{}, gocassa.Keys{
		PartitionKeys: []string{"Id"},
	}).WithOptions(gocassa.Options{
		TableName: c.AccountsIDTable.Name(),
	}.Merge(opts))

	c.AccountsEmailTable = conn.KeySpace(keyspace).MapTable(
		"accounts",
		"Email",
		&emailIDMap{},
	).WithOptions(opts)

	return nil
}

func (c *Cassandra) Close() error {
	c.Session.Close()
	return nil
}

func (c *Cassandra) Truncate() error {
	c.Session.Query("truncate table accounts_map_id").Exec()
	c.Session.Query("truncate table accounts_map_email").Exec()
	return nil
}

func (c *Cassandra) List(count int32, token string) (accounts []*Account, next_token string, err error) {
	if token == "" {
		token = "0"
	}

	rows := c.Session.Query(
		"select id,name,email,createdat from accounts_map_id where token(id) > token(?) limit ?",
		token, count).Iter().Scanner()

	if rows == nil {
		return accounts, next_token, err
	}

	for rows.Next() {
		a := &Account{}
		err := rows.Scan(&a.ID, &a.Name, &a.Email, &a.CreatedAt)
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

	err = c.AccountsIDTable.Set(a).Run()
	if err != nil {
		return err
	}

	return c.createEmailRow(a)
}

func (c *Cassandra) ReadByEmail(email string) (*Account, error) {
	a := &Account{}
	err := c.AccountsEmailTable.Read(email, a).Run()

	switch err.(type) {
	case gocassa.RowNotFoundError:
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	a, err = c.ReadByID(a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (c *Cassandra) ReadByID(ID string) (*Account, error) {
	a := Account{}
	err := c.AccountsIDTable.Read(ID, &a).Run()

	switch err.(type) {
	case gocassa.RowNotFoundError:
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
		err = c.AccountsEmailTable.Delete(current.Email).Run()
		if err != nil {
			return nil
		}

		err = c.createEmailRow(a)
		if err != nil {
			return nil
		}
	}

	return c.AccountsIDTable.Update(a.ID, a.ToMap()).Run()
}

func (c *Cassandra) Delete(ID string) error {
	a, err := c.ReadByID(ID)
	if err != nil {
		return err
	}

	err = c.AccountsIDTable.Delete(ID).Run()
	if err != nil {
		return err
	}

	return c.AccountsEmailTable.Delete(a.Email).Run()
}

func (c *Cassandra) createEmailRow(a *Account) error {
	ea := emailIDMap{
		Email: a.Email,
		ID:    a.ID,
	}
	return c.AccountsEmailTable.Set(&ea).Run()
}
