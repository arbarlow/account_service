package server

import (
	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (as AccountServer) GeneratePasswordToken(ctx context.Context, r *account_service.GeneratePasswordTokenRequest) (*account_service.GeneratePasswordTokenResponse, error) {
	a, err := as.DB.GeneratePasswordToken(r.Email)
	if err != nil {
		if err == database.ErrAccountNotFound {
			return nil, grpc.Errorf(codes.NotFound, "account not found")
		}
		return nil, err
	}

	return &account_service.GeneratePasswordTokenResponse{
		Token: a.PasswordResetToken,
	}, nil
}
