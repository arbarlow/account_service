package server

import (
	"testing"

	"github.com/lileio/account_service"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestConfirmAccount(t *testing.T) {
	ctx := context.Background()

	ac := createAccount(t)
	req := &account_service.ConfirmAccountRequest{Token: ac.ConfirmToken}

	res, err := as.ConfirmAccount(ctx, req)
	assert.Nil(t, err)
	assert.Empty(t, res.ConfirmToken)
	assert.NotNil(t, res)
}
