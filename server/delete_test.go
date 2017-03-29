package server

import (
	"testing"

	"github.com/lileio/account_service"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestDeleteSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	dr := &account_service.DeleteAccountRequest{Id: a.Id}
	res, err := as.Delete(ctx, dr)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
}

func TestDeleteAccountNotExist(t *testing.T) {
	truncate()

	ctx := context.Background()
	u1 := uuid.NewV1()

	dr := &account_service.DeleteAccountRequest{Id: u1.String()}
	_, err := as.Delete(ctx, dr)
	assert.NotNil(t, err)
}
