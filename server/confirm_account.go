package server

import (
	"github.com/lileio/account_service"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (as AccountServer) ConfirmAccount(ctx context.Context, r *account_service.ConfirmAccountRequest) (*account_service.Account, error) {
	a, err := as.DB.Confirm(r.Token)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "token incorrect %s", err)
	}

	return accountDetailsFromAccount(a), nil
}
