package server

import (
	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
)

func (as AccountServer) Create(ctx context.Context, r *account_service.CreateAccountRequest) (*account_service.Account, error) {
	if r.Account == nil {
		return nil, ErrNoAccount
	}

	a := database.Account{
		Name:  r.Account.Name,
		Email: r.Account.Email,
	}

	err := a.Valid()
	if err != nil {
		return nil, err
	}

	err = database.EmailExists(as.DB, &a)
	if err != nil {
		return nil, err
	}

	err = a.HashPassword(r.Password)
	if err != nil {
		return nil, err
	}

	if r.Image != nil {
		err := as.storeImage(ctx, r.Image, &a)
		if err != nil {
			return nil, errors.Wrap(err, "image service:")
		}
	}

	err = as.DB.Create(&a, r.Password)
	if err != nil {
		as.deleteImages(ctx, &a)
		return nil, err
	}

	return accountDetailsFromAccount(&a), nil
}
