package server

import (
	"context"
	"strconv"
	"testing"

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

func TestCreateSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.CreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}
	account, err := as.Create(ctx, req)

	assert.Nil(t, err)
	assert.NotEmpty(t, account.Id)
}

func BenchmarkCreate(b *testing.B) {
	truncate()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		req := &account.CreateRequest{
			Name:  "Alex B",
			Email: "alexbarlowis@localhost" + strconv.Itoa(i),
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
	req := &account.CreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a, err := as.Create(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Id)

	req2 := &account.CreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a2, err := as.Create(ctx, req2)
	assert.NotNil(t, err)
	assert.Equal(t, err, database.ErrEmailExists)
	assert.Nil(t, a2)
}

func TestCreateEmpty(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.CreateRequest{
		Name:  "Alex B",
		Email: "",
	}

	account, err := as.Create(ctx, req)
	assert.NotNil(t, err)
	assert.Nil(t, account)
}

func TestGetByIdSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.CreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a, err := as.Create(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Id)

	areq := &account.GetByIdRequest{
		Id: a.Id,
	}

	a2, err := as.GetById(ctx, areq)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.NotEmpty(t, a2.Email)
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
	req := &account.CreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a, err := as.Create(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Id)

	areq := &account.GetByEmailRequest{
		Email: "alexbarlowis@localhost",
	}

	a2, err := as.GetByEmail(ctx, areq)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.NotEmpty(t, a2.Email)
}

func TestUpdateSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account.CreateRequest{
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

	a3, err := as.GetById(ctx, &account.GetByIdRequest{Id: a2.Id})
	assert.Nil(t, err)
	assert.Equal(t, a3.Email, email)
}

func TestUpdateNotExist(t *testing.T) {
	truncate()
	ctx := context.Background()
	u1 := uuid.NewV1()

	email := "somethingnew@gmail.com"

	a := &account.Account{
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
	req := &account.CreateRequest{
		Name:  "Alex B",
		Email: "alexbarlowis@localhost",
	}

	a, err := as.Create(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Id)

	dr := &account.DeleteRequest{Id: a.Id}
	res, err := as.Delete(ctx, dr)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.Id)
}

func TestDeleteAccountNotExist(t *testing.T) {
	truncate()

	ctx := context.Background()
	u1 := uuid.NewV1()

	dr := &account.DeleteRequest{Id: u1.String()}
	_, err := as.Delete(ctx, dr)
	assert.NotNil(t, err)
}
