package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (as AccountServer) Delete(ctx context.Context, r *account_service.DeleteAccountRequest) (*empty.Empty, error) {
	ca, err := as.DB.ReadByID(r.Id)
	if err != nil {
		if err == database.ErrAccountNotFound {
			return nil, grpc.Errorf(codes.NotFound, "account not found")
		}
		return nil, err
	}

	err = as.deleteImages(ctx, ca)
	if err != nil {
		return nil, err
	}

	a := database.Account{ID: r.Id}
	err = as.deleteImages(ctx, &a)
	if err != nil {
		return nil, err
	}

	err = as.DB.Delete(r.Id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
