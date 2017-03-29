package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/gemnasium/migrate/driver/postgres"
	"github.com/gemnasium/migrate/migrate"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/sirupsen/logrus"
)

var TOKEN_LENGTH = 32

func (a *Account) BeforeInsert(db orm.DB) error {
	if a.ConfirmationToken == "" {
		t, err := GenerateRandomString(TOKEN_LENGTH)
		if err != nil {
			logrus.Errorf("confirm token generation error %v", err)
			return err
		}

		a.ConfirmationToken = t
	}
	return nil
}

type PostgreSQL struct {
	Database
	db   *pg.DB
	conn string
}

func (p *PostgreSQL) Connect(conn string) error {
	p.conn = conn
	opts, err := pg.ParseURL(conn)
	if err != nil {
		return err
	}

	p.db = pg.Connect(opts)

	p.db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		logrus.Debugf("SQL[%s]: %s", time.Since(event.StartTime), query)
	})

	return nil
}

func (p *PostgreSQL) Migrate() error {
	wd := os.ExpandEnv("$GOPATH/src/github.com/lileio/account_service")
	allErrors, ok := migrate.UpSync(p.conn, wd+"/migrations/pg")
	if !ok {
		fmt.Printf("migration failed: %+v\n", allErrors)
		return errors.New("migration error")
	}

	return nil
}

func (p *PostgreSQL) Close() error {
	return p.db.Close()
}

func (p *PostgreSQL) Truncate() error {
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

	err = p.db.Model(&Account{}).
		Column("account.*").
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
	err := p.db.Insert(&a)
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

func (p *PostgreSQL) GeneratePasswordToken(email string) (*Account, error) {
	a, err := p.ReadByEmail(email)
	if err != nil {
		return nil, err
	}

	t, err := GenerateRandomString(TOKEN_LENGTH)
	if err != nil {
		logrus.Errorf("password token generation error %v", err)
		return nil, err
	}

	a.PasswordResetToken = t
	err = p.db.Update(a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (p *PostgreSQL) UpdatePassword(token, hashed_password string) (*Account, error) {
	var a Account
	err := p.db.Model(&a).
		Where("password_reset_token = ?", token).
		Select()
	if err != nil && notFoundError(err) {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	a.HashedPassword = hashed_password
	err = p.db.Update(&a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (p *PostgreSQL) Confirm(token string) (*Account, error) {
	var a Account
	err := p.db.Model(&a).
		Where("confirmation_token = ?", token).
		Select()
	if err != nil && notFoundError(err) {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	a.ConfirmationToken = ""
	err = p.db.Update(&a)
	if err != nil {
		return nil, err
	}

	return &a, nil
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
