package server

import (
	"context"
	"flag"
	"testing"

	"github.com/arbarlow/account_service/account"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var dbConnect string
var db = setupDB()
var as = AccountServer{}

func setupDB() *sqlx.DB {
	conn := "host=localhost user=postgres dbname=account_service_test sslmode=disable"
	flag.StringVar(&dbConnect, "connect", conn, "db connection string")
	flag.Parse()

	db, err := as.DBConnect(dbConnect)
	if err != nil {
		panic(err)
	}

	return db
}

func truncate() {
	db.MustExec("TRUNCATE accounts;")
}

func TestCreateSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.AccountCreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}
	account, err := as.Create(ctx, req)

	assert.NotEmpty(t, account.Id)
	assert.Nil(t, err)
}

func TestCreateUniqueness(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.AccountCreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a, err := as.Create(ctx, req)
	assert.NotEmpty(t, a.Id)
	assert.Nil(t, err)

	req2 := &account.AccountCreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a2, err := as.Create(ctx, req2)
	assert.Empty(t, a2.Id)
	assert.NotNil(t, err)
}

func TestCreateEmpty(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.AccountCreateRequest{
		Name:  "Alex B",
		Email: "",
	}

	account, err := as.Create(ctx, req)
	assert.Empty(t, account.Id)
	assert.NotNil(t, err)
}

func TestReadSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.AccountCreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a, err := as.Create(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Id)

	areq := &account.AccountRequest{
		Id: a.Id,
	}

	a2, err := as.Read(ctx, areq)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.NotEmpty(t, a2.Email)
}

func TestReadFail(t *testing.T) {
	truncate()

	areq := &account.AccountRequest{
		Id: "somefalseid",
	}

	ctx := context.Background()
	a2, err := as.Read(ctx, areq)
	assert.NotNil(t, err)
	assert.Nil(t, a2)
}

func TestUpdateSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.AccountCreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a, err := as.Create(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Id)

	email := "somethingnew@gmail.com"
	a.Email = email

	a2, err := as.Update(ctx, a)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.Equal(t, a2.Email, email)

	a3, err := as.Read(ctx, &account.AccountRequest{Id: a2.Id})
	assert.Nil(t, err)
	assert.Equal(t, a3.Email, email)
}

func TestUpdateNotExist(t *testing.T) {
	truncate()
	ctx := context.Background()
	u1 := uuid.NewV1()

	email := "somethingnew@gmail.com"

	a := &account.AccountDetails{
		Id:    u1.String(),
		Email: email,
	}

	a2, err := as.Update(ctx, a)
	assert.NotNil(t, err)
	assert.Nil(t, a2)
}

func TestDeleteSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.AccountCreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a, err := as.Create(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Id)

	dr := &account.AccountDeleteRequest{Id: a.Id}
	res, err := as.Delete(ctx, dr)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.Id)
}

func TestDeleteAccountNotExist(t *testing.T) {
	truncate()

	ctx := context.Background()
	u1 := uuid.NewV1()

	dr := &account.AccountDeleteRequest{Id: u1.String()}
	_, err := as.Delete(ctx, dr)
	assert.NotNil(t, err)
}
