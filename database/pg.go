package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	pg "gopkg.in/pg.v5"

	_ "github.com/gemnasium/migrate/driver/postgres"
	"github.com/gemnasium/migrate/migrate"
)

type PostgreSQL struct {
	Database
	db *pg.DB
}

func (p *PostgreSQL) Connect(conn string) error {
	opts, err := pg.ParseURL(conn)
	if err != nil {
		return err
	}

	p.db = pg.Connect(opts)

	wd := os.ExpandEnv("$GOPATH/src/github.com/lileio/account_service")
	allErrors, ok := migrate.UpSync(conn, wd+"/migrations/pg")
	if !ok {
		fmt.Printf("allErrors = %+v\n", allErrors)
		return errors.New("migration error")
	}

	return nil
}

func (p *PostgreSQL) Close() error {
	return p.db.Close()
}

func (p *PostgreSQL) Truncate(reconnect bool) error {
	p.db.Exec("TRUNCATE accounts;")
	return nil
}

func (p *PostgreSQL) List(count32 int32, token string) (accounts []*Account, next_token string, err error) {
	count := int(count32)
	if token == "" {
		token = "0"
	}

	offset, err := strconv.Atoi(token)
	if err != nil {
		return accounts, next_token, err
	}

	err = p.db.Model(&Account{}).Column("account.*").
		Limit(count).
		Offset(offset).
		Select(&accounts)

	if err != nil {
		return accounts, next_token, err
	}

	if len(accounts) == int(count) {
		next_token = strconv.FormatInt(int64(offset+count+1), 10)
	}

	return accounts, next_token, err
}

func (p *PostgreSQL) ReadByID(ID string) (*Account, error) {
	a := Account{ID: ID}
	err := p.db.Select(&a)
	if err != nil && notFoundError(err) {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (p *PostgreSQL) ReadByEmail(email string) (*Account, error) {
	a := Account{}
	err := p.db.Model(&a).Where("email = ?", email).Select()
	if err != nil && notFoundError(err) {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (p *PostgreSQL) Create(a *Account, password string) error {
	err := a.Valid()
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

	err = p.db.Insert(&a)
	if err != nil && uniqueEmailError(err) {
		return ErrEmailExists
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) Update(a *Account) error {
	err := a.Valid()
	if err != nil {
		return err
	}

	_, err = p.db.Model(&a).
		Column("name", "email", "images").
		Returning("*").
		Update()
	if err != nil && uniqueEmailError(err) {
		return ErrEmailExists
	}

	if err != nil && notFoundError(err) {
		return ErrAccountNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) Delete(ID string) error {
	a := Account{ID: ID}
	err := p.db.Delete(&a)
	if err != nil {
		return err
	}

	return nil
}

func uniqueEmailError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint") && strings.Contains(err.Error(), "email")
}

func notFoundError(err error) bool {
	return strings.Contains(err.Error(), "no rows in result")
}
