package database

import (
	"errors"
	"fmt"
	"strings"

	_ "github.com/gemnasium/migrate/driver/postgres"
	"github.com/gemnasium/migrate/migrate"
	"github.com/jmoiron/sqlx"
)

type PostgreSQL struct {
	Database
	DB *sqlx.DB
}

func (p *PostgreSQL) Connect(conn string) error {
	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return err
	}

	allErrors, ok := migrate.UpSync(conn, "../migrations/pg")
	if !ok {
		fmt.Printf("allErrors = %+v\n", allErrors)
		return errors.New("migration error")
	}

	p.DB = db
	return nil
}

func (p *PostgreSQL) Close() error {
	return p.DB.Close()
}

func (p *PostgreSQL) Truncate() error {
	p.DB.MustExec("TRUNCATE accounts;")
	return nil
}

func (p *PostgreSQL) ReadByID(ID string) (*Account, error) {
	var a Account
	err := p.DB.Get(&a, "SELECT * FROM accounts WHERE id = $1", ID)
	if err != nil && notFoundError(err) {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (p *PostgreSQL) ReadByEmail(email string) (*Account, error) {
	var a Account
	err := p.DB.Get(&a, "SELECT * FROM accounts WHERE email = $1", email)
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

	sql := `
	INSERT INTO accounts (name, email, hashed_password)
	VALUES (:name, :email, :hashed_password) RETURNING id`

	stmt, err := p.DB.PrepareNamed(sql)
	if err != nil {
		return err
	}

	var id string
	err = stmt.Get(&id, a)
	if err != nil && uniqueEmailError(err) {
		return ErrEmailExists
	}

	if err != nil {
		return err
	}

	a.ID = id

	return nil
}

func (p *PostgreSQL) Update(a *Account) error {
	err := a.Valid()
	if err != nil {
		return err
	}

	sql := `UPDATE accounts SET name = :name, email = :email WHERE id = :id`
	res, err := p.DB.NamedExec(sql, a)
	if err != nil && uniqueEmailError(err) {
		return ErrEmailExists
	}

	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return ErrAccountNotFound
	}

	return nil
}

func (p *PostgreSQL) Delete(ID string) error {
	sql := `DELETE from accounts WHERE id = $1`
	res, err := p.DB.Exec(sql, ID)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if aff == 0 {
		return ErrAccountNotFound
	}

	return nil
}

func uniqueEmailError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint") && strings.Contains(err.Error(), "email")
}

func notFoundError(err error) bool {
	return strings.Contains(err.Error(), "no rows in result")
}
