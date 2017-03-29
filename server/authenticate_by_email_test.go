package server

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/lileio/account_service"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestAuthenticate(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	ar := &account_service.AuthenticateByEmailRequest{
		Email:    a.Email,
		Password: pass,
	}

	a, err := as.AuthenticateByEmail(ctx, ar)
	assert.Nil(t, err)
	assert.NotEmpty(t, a.Id)
	assert.NotEmpty(t, a.Email)
}

func TestAuthenticateFailure(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	ar := &account_service.AuthenticateByEmailRequest{
		Email:    a.Email,
		Password: "incorrect password lol",
	}

	_, err := as.AuthenticateByEmail(ctx, ar)
	assert.NotNil(t, err)
	assert.Equal(t, grpc.Code(err), codes.PermissionDenied)
}
