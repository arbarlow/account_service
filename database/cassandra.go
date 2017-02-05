package database

import (
	"time"

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
	Session            *gocql.Session
	AccountsIDTable    gocassa.MapTable
	AccountsEmailTable gocassa.MapTable
}

func (c *Cassandra) Connect(keyspace string, addrs []string) error {
	cluster := gocql.NewCluster(addrs...)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	c.Session = session

	conn := gocassa.NewConnection(gocassa.GoCQLSessionToQueryExecutor(session))

	c.AccountsIDTable = conn.KeySpace(keyspace).MapTable(
		"accounts",
		"ID",
		&Account{},
	)
	err = c.AccountsIDTable.CreateIfNotExist()
	if err != nil {
		return err
	}

	c.AccountsEmailTable = conn.KeySpace(keyspace).MapTable(
		"accounts",
		"Email",
		&emailIDMap{},
	)
	err = c.AccountsEmailTable.CreateIfNotExist()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cassandra) Close() error {

	c.Session.Close()
	return nil
}

func (c *Cassandra) Truncate() error {
	return c.Session.Query("truncate account_service.accounts_map_email").Exec()
}

func (p *Cassandra) Create(a *Account) error {
	// Check the email table to see if this account exists
	ae, err := p.ReadByEmail(a.Email)
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

	err = p.AccountsIDTable.Set(a).Run()
	if err != nil {
		return err
	}

	return p.createEmailRow(a)
}

func (p *Cassandra) ReadByEmail(email string) (*Account, error) {
	a := Account{}
	err := p.AccountsEmailTable.Read(email, &a).Run()

	switch err.(type) {
	case gocassa.RowNotFoundError:
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (p *Cassandra) ReadByID(ID string) (*Account, error) {
	a := Account{}
	err := p.AccountsIDTable.Read(ID, &a).Run()

	switch err.(type) {
	case gocassa.RowNotFoundError:
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (p *Cassandra) Update(a *Account) error {
	current, err := p.ReadByID(a.ID)
	if err != nil {
		return err
	}

	if current.Email != a.Email {
		err = p.AccountsEmailTable.Delete(current.Email).Run()
		if err != nil {
			return nil
		}

		err = p.createEmailRow(a)
		if err != nil {
			return nil
		}
	}

	return p.AccountsIDTable.Update(a.ID, a.ToMap()).Run()
}

func (p *Cassandra) Delete(ID string) error {
	a, err := p.ReadByID(ID)
	if err != nil {
		return err
	}

	err = p.AccountsIDTable.Delete(ID).Run()
	if err != nil {
		return err
	}

	return p.AccountsEmailTable.Delete(a.Email).Run()
}

func (p *Cassandra) createEmailRow(a *Account) error {
	ea := emailIDMap{
		Email: a.Email,
		ID:    a.ID,
	}
	return p.AccountsEmailTable.Set(&ea).Run()
}