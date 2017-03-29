package server

import (
	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (as AccountServer) Update(ctx context.Context, r *account_service.UpdateAccountRequest) (*account_service.Account, error) {
	if r.Account == nil {
		return nil, ErrNoAccount
	}

	a := database.Account{
		ID:       r.Id,
		Name:     r.Account.Name,
		Email:    r.Account.Email,
		Metadata: r.Account.Metadata,
	}

	if r.Image != nil {
		ca, err := as.DB.ReadByID(a.ID)
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

		err = as.storeImage(ctx, r.Image, &a)
		if err != nil {
			return nil, err
		}
	}

	err := as.DB.Update(&a)
	if err != nil {
		if err == database.ErrAccountNotFound {
			return nil, grpc.Errorf(codes.NotFound, "account not found")
		}
		return nil, err
	}

	return accountDetailsFromAccount(&a), nil
}
