package server

import (
	context "golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
)

func (as AccountServer) ResetPassword(ctx context.Context, r *account_service.ResetPasswordRequest) (*account_service.Account, error) {
	a := database.Account{}

	err := a.HashPassword(r.Password)
	if err != nil {
		return nil, err
	}

	ac, err := as.DB.UpdatePassword(r.Token, a.HashedPassword)
	if err != nil {
		if err == database.ErrAccountNotFound {
			return nil, grpc.Errorf(codes.NotFound, "account not found")
		}
		return nil, err
	}

	return accountDetailsFromAccount(ac), nil
}
