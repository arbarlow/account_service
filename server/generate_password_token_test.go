package server

import (
	"testing"

	"github.com/lileio/account_service"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestGeneratePasswordToken(t *testing.T) {
	ctx := context.Background()

	ac := createAccount(t)
	req := &account_service.GeneratePasswordTokenRequest{Email: ac.Email}
	res, err := as.GeneratePasswordToken(ctx, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
}
