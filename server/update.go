package server

import (
	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	context "golang.org/x/net/context"
)

func (as AccountServer) Update(ctx context.Context, r *account_service.UpdateAccountRequest) (*account_service.Account, error) {
	if r.Account == nil {
		return nil, ErrNoAccount
	}

	a := database.Account{
		ID:    r.Id,
		Name:  r.Account.Name,
		Email: r.Account.Email,
	}

	if r.Image != nil {
		ca, err := as.DB.ReadByID(a.ID)
		if err != nil {
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
		return nil, err
	}

	return accountDetailsFromAccount(&a), nil
}
