package server

import (
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
)

func (as AccountServer) GetById(ctx context.Context, r *account_service.GetByIdRequest) (*account_service.Account, error) {
	a, err := as.DB.ReadByID(r.Id)
	if err != nil {
		if err == database.ErrAccountNotFound {
			return nil, grpc.Errorf(codes.NotFound, "account not found")
		}
		return nil, err
	}

	return accountDetailsFromAccount(a), nil
}
