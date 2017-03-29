package server

import (
	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (as AccountServer) GetByEmail(ctx context.Context, r *account_service.GetByEmailRequest) (*account_service.Account, error) {
	a, err := as.DB.ReadByEmail(r.Email)
	if err != nil {
		if err == database.ErrAccountNotFound {
			return nil, grpc.Errorf(codes.NotFound, "account not found")
		}
		return nil, err
	}

	return accountDetailsFromAccount(a), nil
}
