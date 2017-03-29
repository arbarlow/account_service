package server

import (
	"os"
	"testing"

	"github.com/lileio/account_service"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestSimpleList(t *testing.T) {
	truncate()

	for i := 0; i < 5; i++ {
		createAccount(t)
	}

	ctx := context.Background()
	req := &account_service.ListAccountsRequest{
		PageSize: 6,
	}

	l, err := as.List(ctx, req)
	assert.Nil(t, err)
	assert.Empty(t, l.NextPageToken)
}

func TestSimpleListToken(t *testing.T) {
	if os.Getenv("CASSANDRA_DB_NAME") != "" {
		t.Skip()
	}

	truncate()

	acc := []*account_service.Account{}
	for i := 0; i < 4; i++ {
		acc = append(acc, createAccount(t))
	}

	ctx := context.Background()
	req := &account_service.ListAccountsRequest{
		PageSize: 2,
	}

	l, err := as.List(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, l.NextPageToken)

	req = &account_service.ListAccountsRequest{
		PageSize:  2,
		PageToken: l.NextPageToken,
	}

	l, err = as.List(ctx, req)
	assert.Nil(t, err)
	assert.Empty(t, l.NextPageToken)
}
