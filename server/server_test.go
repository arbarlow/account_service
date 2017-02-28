package server

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/gocql/gocql"
	_ "github.com/lib/pq"
	"github.com/lileio/account_service/account"
	"github.com/lileio/account_service/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var db = setupDB()
var as = AccountServer{}

func setupDB() database.Database {
	conn := database.DatabaseFromEnv()
	as.DB = conn
	return conn
}

func truncate() {
	db.Truncate()
}

var name = "Alex B"
var email = "alexb@localhost"
var pass = "password"

func createAccount(t *testing.T) *account.Account {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx := context.Background()
	req := &account.CreateAccountRequest{
		Account: &account.Account{
			Name:  name,
			Email: "alexbarlowis@localhost" + strconv.Itoa(r.Int()),
		},
		Password: pass,
	}
	account, err := as.Create(ctx, req)
	assert.Nil(t, err)
	return account
}

func TestSimpleList(t *testing.T) {
	truncate()

	for i := 0; i < 5; i++ {
		createAccount(t)
	}

	ctx := context.Background()
	req := &account.ListAccountsRequest{
		PageSize: 6,
	}

	l, err := as.List(ctx, req)
	assert.Nil(t, err)
	assert.Empty(t, l.NextPageToken)
}

func TestSimpleListToken(t *testing.T) {
	a := gocql.All
	database.Consistency = &a
	truncate()

	acc := []*account.Account{}
	for i := 0; i < 4; i++ {
		acc = append(acc, createAccount(t))
	}

	ctx := context.Background()
	req := &account.ListAccountsRequest{
		PageSize: 2,
	}

	l, err := as.List(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, l.NextPageToken)

	req = &account.ListAccountsRequest{
		PageSize:  2,
		PageToken: l.NextPageToken,
	}

	l, err = as.List(ctx, req)
	assert.Nil(t, err)
	assert.Empty(t, l.NextPageToken)

	database.Consistency = nil
}

func TestCreateSuccess(t *testing.T) {
	truncate()

	account := createAccount(t)
	assert.NotEmpty(t, account.Id)
}

func BenchmarkCreate(b *testing.B) {
	truncate()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		req := &account.CreateAccountRequest{
			Account: &account.Account{
				Name:  name,
				Email: "alexbarlowis@localhost" + strconv.Itoa(i),
			},
			Password: pass,
		}

		_, err := as.Create(ctx, req)

		if err != nil {
			panic(err)
		}
	}
}

func TestCreateUniqueness(t *testing.T) {
	truncate()

	ctx := context.Background()
	a1 := createAccount(t)

	req2 := &account.CreateAccountRequest{
		Account:  a1,
		Password: pass,
	}

	a2, err := as.Create(ctx, req2)
	assert.NotNil(t, err)
	assert.Equal(t, err, database.ErrEmailExists)
	assert.Nil(t, a2)
}

func TestCreateEmpty(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.CreateAccountRequest{}

	account, err := as.Create(ctx, req)
	assert.NotNil(t, err)
	assert.Nil(t, account)
}

func TestGetByIdSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	areq := &account.GetByIdRequest{
		Id: a.Id,
	}

	a2, err := as.GetById(ctx, areq)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.NotEmpty(t, a2.Email)
}

func TestAuthenticate(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	ar := &account.AuthenticateByEmailRequest{
		Email:    a.Email,
		Password: pass,
	}

	a, err := as.AuthenticateByEmail(ctx, ar)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Id)
	assert.NotEmpty(t, a.Email)
}

func TestAuthenticateFailure(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	ar := &account.AuthenticateByEmailRequest{
		Email:    a.Email,
		Password: "incorrect password lol",
	}

	_, err := as.AuthenticateByEmail(ctx, ar)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrAuthFail)
}

func TestGetByIdFail(t *testing.T) {
	truncate()

	u1 := uuid.NewV1()

	areq := &account.GetByIdRequest{
		Id: u1.String(),
	}

	ctx := context.Background()
	a2, err := as.GetById(ctx, areq)
	assert.NotNil(t, err)
	assert.Equal(t, err, database.ErrAccountNotFound)
	assert.Nil(t, a2)
}

func TestGetByEmailSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	areq := &account.GetByEmailRequest{
		Email: a.Email,
	}

	a2, err := as.GetByEmail(ctx, areq)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.NotEmpty(t, a2.Email)
}

func TestUpdateSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	email := "somethingnew@gmail.com"
	a.Email = email

	ar := &account.UpdateAccountRequest{
		Id:      a.Id,
		Account: a,
	}

	a2, err := as.Update(ctx, ar)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.Equal(t, a2.Email, email)

	a3, err := as.GetById(ctx, &account.GetByIdRequest{Id: a2.Id})
	assert.Nil(t, err)
	assert.Equal(t, a3.Email, email)
}

func TestUpdateNotExist(t *testing.T) {
	truncate()
	ctx := context.Background()
	u1 := uuid.NewV1()

	ar := &account.UpdateAccountRequest{
		Id: u1.String(),
	}

	a2, err := as.Update(ctx, ar)
	assert.NotNil(t, err)
	assert.Nil(t, a2)
}

func TestDeleteSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	dr := &account.DeleteAccountRequest{Id: a.Id}
	res, err := as.Delete(ctx, dr)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
}

func TestDeleteAccountNotExist(t *testing.T) {
	truncate()

	ctx := context.Background()
	u1 := uuid.NewV1()

	dr := &account.DeleteAccountRequest{Id: u1.String()}
	_, err := as.Delete(ctx, dr)
	assert.NotNil(t, err)
}
