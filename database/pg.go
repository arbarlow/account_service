package database

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type PostgreSQL struct {
	Database
	DB *sqlx.DB
}

const schema = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS accounts (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	name text NULL,
	email text NOT NULL,
	created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE UNIQUE INDEX IF NOT EXISTS accounts_email ON accounts ((lower(email)));
`

func (p *PostgreSQL) Connect(conn string) error {
	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return err
	}

	db.MustExec(schema)

	p.DB = db
	return nil
}

func (p PostgreSQL) Close() error {
	return p.DB.Close()
}

func (p PostgreSQL) ReadByID(ID string) (*Account, error) {
	var a Account
	err := p.DB.Get(&a, "SELECT * FROM accounts WHERE id = $1", ID)
	if err != nil && notFoundError(err) {
		return nil, ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	if a.Id == "" {
		return nil, ErrAccountNotFound
	}

	return &a, nil
}

func (p PostgreSQL) Create(a *Account) error {
	sql := "INSERT INTO accounts (name, email) VALUES (:name, :email) RETURNING id"
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

	a.Id = id

	return nil
}

func (p PostgreSQL) Update(a *Account) error {
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

func (p PostgreSQL) Delete(ID string) error {
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
