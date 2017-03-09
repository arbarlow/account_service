package server

import (
	"github.com/lileio/account_service"
	context "golang.org/x/net/context"
)

func (as AccountServer) GetById(ctx context.Context, r *account_service.GetByIdRequest) (*account_service.Account, error) {
	a, err := as.DB.ReadByID(r.Id)
	if err != nil {
		return nil, err
	}

	return accountDetailsFromAccount(a), nil
}

func (as AccountServer) GetByEmail(ctx context.Context, r *account_service.GetByEmailRequest) (*account_service.Account, error) {
	a, err := as.DB.ReadByEmail(r.Email)
	if err != nil {
		return nil, err
	}

	return accountDetailsFromAccount(a), nil
}
