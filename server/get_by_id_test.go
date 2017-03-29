package server

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/lileio/account_service"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestGetByIdSuccess(t *testing.T) {
	truncate()

	ctx := context.Background()
	a := createAccount(t)

	areq := &account_service.GetByIdRequest{
		Id: a.Id,
	}

	a2, err := as.GetById(ctx, areq)
	assert.Nil(t, err)
	assert.NotEmpty(t, a2.Id)
	assert.NotEmpty(t, a2.Email)
}

func TestGetByIdFail(t *testing.T) {
	truncate()

	u1 := uuid.NewV1()

	areq := &account_service.GetByIdRequest{
		Id: u1.String(),
	}

	ctx := context.Background()
	a2, err := as.GetById(ctx, areq)
	assert.NotNil(t, err)
	assert.Equal(t, grpc.Code(err), codes.NotFound)
	assert.Nil(t, a2)
}
