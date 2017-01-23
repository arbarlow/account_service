package server

import (
	"database/sql"
	"errors"
	"time"

	context "golang.org/x/net/context"

	"github.com/arbarlow/account_service/account"
	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS accounts (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	name text NULL,
	email text NOT NULL,
	created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE UNIQUE INDEX IF NOT EXISTS accounts_email ON accounts ((lower(email)));
`

type Account struct {
	Id        string `db:"id"`
	Name      sql.NullString
	Email     sql.NullString
	CreatedAt time.Time `db:"created_at"`
}

type AccountServer struct {
	account.AccountServiceServer
	db *sqlx.DB
}

var (
	ErrAccountNotFound = errors.New("Account not found")
)

func (as *AccountServer) DBConnect(conn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	db.MustExec(schema)

	as.db = db
	return db, nil
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func (as AccountServer) Create(ctx context.Context, r *account.AccountCreateRequest) (*account.AccountDetails, error) {
	a := Account{
		Name:  nullString(r.Name),
		Email: nullString(r.Email),
	}

	sql := "INSERT INTO accounts (name, email) VALUES (:name, :email) RETURNING id"
	stmt, err := as.db.PrepareNamed(sql)
	if err != nil {
		return nil, err
	}

	var id string
	err = stmt.Get(&id, &a)

	return &account.AccountDetails{
		Id:    id,
		Name:  a.Name.String,
		Email: a.Email.String,
	}, err
}

func (as AccountServer) Read(ctx context.Context, r *account.AccountRequest) (*account.AccountDetails, error) {
	var a Account
	err := as.db.Get(&a, "SELECT * FROM accounts WHERE id = $1", r.Id)
	if err != nil {
		return nil, err
	}

	if a.Id == "" {
		return nil, ErrAccountNotFound
	}

	return &account.AccountDetails{
		Id:    a.Id,
		Name:  a.Name.String,
		Email: a.Email.String,
	}, nil
}

func (as AccountServer) Update(ctx context.Context, r *account.AccountDetails) (*account.AccountDetails, error) {
	a := Account{
		Id:    r.Id,
		Name:  nullString(r.Name),
		Email: nullString(r.Email),
	}

	sql := `UPDATE accounts SET name = :name, email = :email WHERE id = :id`
	res, err := as.db.NamedExec(sql, &a)
	if err != nil {
		return nil, err
	}

	aff, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if aff == 0 {
		return nil, ErrAccountNotFound
	}

	return &account.AccountDetails{
		Id:    a.Id,
		Name:  a.Name.String,
		Email: a.Email.String,
	}, nil
}

func (as AccountServer) Delete(ctx context.Context, r *account.AccountDeleteRequest) (*account.AccountDeleteResponse, error) {
	sql := `DELETE from accounts WHERE id = :id`
	res, err := as.db.NamedExec(sql, &r)
	if err != nil {
		return nil, err
	}

	aff, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if aff == 0 {
		return nil, ErrAccountNotFound
	}

	return &account.AccountDeleteResponse{Id: r.Id}, nil
}
