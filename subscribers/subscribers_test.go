package subscribers

import (
	"context"
	"testing"

	account "github.com/lileio/account_service"
	"github.com/stretchr/testify/assert"
)

var as = AccountServiceSubscribers{}

func TestAccountCreated(t *testing.T) {
	ctx := context.Background()
	a := account.Account{
		Name:  "alex",
		Email: "alexbarlowis@localhost",
	}

	err := as.AccountCreated(ctx, a)
	assert.Nil(t, err)
}
