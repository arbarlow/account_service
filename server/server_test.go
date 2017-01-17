package server

import (
	"context"
	"flag"
	"log"
	"testing"

	"github.com/arbarlow/account_service/account"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var dbConnect string
var db = setupDB()
var as = AccountServer{}

func setupDB() *gorm.DB {
	conn := "host=localhost user=postgres dbname=account_service_test sslmode=disable"
	flag.StringVar(&dbConnect, "connect", conn, "db connection string")
	flag.Parse()

	db, err := as.DBConnect(dbConnect)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db.LogMode(false)

	return db
}

func truncate() {
	db.Exec("TRUNCATE accounts;")
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

	account, err := as.Create(ctx, req)
	assert.NotEmpty(t, account.Id)
	assert.Nil(t, err)

	account, err = as.Create(ctx, req)
	assert.Empty(t, account.Id)
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
