package server

import (
	"github.com/lileio/account_service"
	context "golang.org/x/net/context"
)

func (as AccountServer) AuthenticateByEmail(ctx context.Context, r *account_service.AuthenticateByEmailRequest) (*account_service.Account, error) {
	a, err := as.DB.ReadByEmail(r.Email)
	if err != nil {
		return nil, err
	}

	if a == nil {
		return nil, ErrAuthFail
	}

	err = a.ComparePasswordToHash(r.Password)
	if err != nil {
		return nil, ErrAuthFail
	}

	return accountDetailsFromAccount(a), nil
}
