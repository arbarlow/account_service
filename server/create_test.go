package server

import (
	"context"
	"strconv"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/lileio/account_service"
	"github.com/stretchr/testify/assert"
)

func TestCreateSuccess(t *testing.T) {
	truncate()

	account := createAccount(t)
	assert.NotEmpty(t, account.Id)
}

func BenchmarkCreate(b *testing.B) {
	truncate()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		req := &account_service.CreateAccountRequest{
			Account: &account_service.Account{
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

	req2 := &account_service.CreateAccountRequest{
		Account:  a1,
		Password: pass,
	}

	a2, err := as.Create(ctx, req2)
	assert.NotNil(t, err)
	assert.Equal(t, grpc.Code(err), codes.AlreadyExists)
	assert.Nil(t, a2)
}

func TestCreateEmpty(t *testing.T) {
	truncate()

	ctx := context.Background()
	req := &account_service.CreateAccountRequest{}

	account, err := as.Create(ctx, req)
	assert.NotNil(t, err)
	assert.Nil(t, account)
}
