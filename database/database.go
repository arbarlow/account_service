package database

import (
	"errors"
	"os"
	"strings"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate           = validator.New()
	ErrAccountNotFound = errors.New("account not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrNoDatabase      = errors.New("no database connection details")
)

type Database interface {
	ReadByID(ID string) (*Account, error)
	ReadByEmail(email string) (*Account, error)
	Create(a *Account) error
	Update(a *Account) error
	Delete(ID string) error
	Truncate() error
	Close() error
}

type Account struct {
	ID        string    `db:"id"`
	Name      string    `validate:"required"`
	Email     string    `validate:"required"`
	CreatedAt time.Time `db:"created_at"`
}

func NewAccount(name, email string) Account {
	return Account{
		Name:  name,
		Email: email,
	}
}

func (a *Account) Valid() error {
	return validate.Struct(a)
}

func (a *Account) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Name":  a.Name,
		"Email": a.Email,
	}
}

func DatabaseFromEnv() Database {
	var conn Database
	var err error

	pg := os.Getenv("POSTGRESQL_URL")
	if pg != "" {
		pgConn := &PostgreSQL{}
		err = pgConn.Connect(pg)
		conn = pgConn
	}

	cass := os.Getenv("CASSANDRA_DB_NAME")
	if cass != "" {
		hosts := strings.Split(os.Getenv("CASSANDRA_HOSTS"), ",")

		cassConn := &Cassandra{}
		err = cassConn.Connect(cass, hosts)
		conn = cassConn
	}

	if conn == nil {
		panic(ErrNoDatabase)
	}

	if err != nil {
		panic(err)
	}

	return conn
}
