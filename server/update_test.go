package server

import (
	"testing"

	"github.com/lileio/account_service"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestUpdateSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	email := "somethingnew@gmail.com"
	a.Email = email

	ar := &account_service.UpdateAccountRequest{
		Id:      a.Id,
		Account: a,
	}

	a2, err := as.Update(ctx, ar)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.Equal(t, a2.Email, email)

	a3, err := as.GetById(ctx, &account_service.GetByIdRequest{Id: a2.Id})
	assert.Nil(t, err)
	assert.Equal(t, a3.Email, email)
}

func TestUpdateNotExist(t *testing.T) {
	truncate()
	ctx := context.Background()
	u1 := uuid.NewV1()

	ar := &account_service.UpdateAccountRequest{
		Id: u1.String(),
	}

	a2, err := as.Update(ctx, ar)
	assert.NotNil(t, err)
	assert.Nil(t, a2)
}
