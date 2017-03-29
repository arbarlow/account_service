package server

import (
	"testing"

	"github.com/lileio/account_service"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestResetPassword(t *testing.T) {
	ctx := context.Background()

	ac := createAccount(t)
	req := &account_service.GeneratePasswordTokenRequest{Email: ac.Email}
	res, err := as.GeneratePasswordToken(ctx, req)
	assert.Nil(t, err)

	newPass := "somenewpassword"

	resetReq := &account_service.ResetPasswordRequest{
		Token:    res.Token,
		Password: newPass,
	}

	r, err := as.ResetPassword(ctx, resetReq)
	assert.Nil(t, err)
	assert.NotNil(t, r)

	ar := &account_service.AuthenticateByEmailRequest{
		Email:    ac.Email,
		Password: newPass,
	}

	_, err = as.AuthenticateByEmail(ctx, ar)
	assert.Nil(t, err)
}
