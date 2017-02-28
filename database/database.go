package database

import (
	"errors"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate           = validator.New()
	ErrAccountNotFound = errors.New("account not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrNoDatabase      = errors.New("no database connection details")
	ErrNoPasswordGiven = errors.New("a password is required")
)

type Database interface {
	List(count int32, token string) ([]*Account, string, error)
	ReadByID(ID string) (*Account, error)
	ReadByEmail(email string) (*Account, error)
	Create(a *Account, password string) error
	Update(a *Account) error
	Delete(ID string) error
	Truncate(reconnect bool) error
	Close() error
}

type Account struct {
	ID             string    `db:"id"`
	Name           string    `validate:"required"`
	Email          string    `validate:"required"`
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at"`
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

func (a *Account) hashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	a.HashedPassword = string(hash[:])

	return nil
}

func (a *Account) ComparePasswordToHash(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.HashedPassword), []byte(password))
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
