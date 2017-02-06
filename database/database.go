package database

import (
	"database/sql"
	"errors"
	"time"
)

type Account struct {
	Id        string `db:"id"`
	Name      sql.NullString
	Email     sql.NullString
	CreatedAt time.Time `db:"created_at"`
}

type Database interface {
	ReadByID(ID string) (*Account, error)
	Create(a *Account) error
	Update(a *Account) error
	Delete(ID string) error
	Close() error
}

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrEmailExists     = errors.New("email already exists")
)

func NewAccount(name, email string) Account {
	return Account{
		Name:  nullString(name),
		Email: nullString(email),
	}
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
