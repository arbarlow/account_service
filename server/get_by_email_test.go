package server

import (
	"testing"

	"github.com/lileio/account_service"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestGetByEmailSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	areq := &account_service.GetByEmailRequest{
		Email: a.Email,
	}

	a2, err := as.GetByEmail(ctx, areq)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.NotEmpty(t, a2.Email)
}
