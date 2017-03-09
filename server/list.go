package server

import (
	"github.com/lileio/account_service"
	context "golang.org/x/net/context"
)

func (as AccountServer) List(
	ctx context.Context, l *account_service.ListAccountsRequest) (
	*account_service.ListAccountsResponse, error) {

	accounts, next_token, err := as.DB.List(l.PageSize, l.PageToken)
	if err != nil {
		return nil, err
	}

	accs := make([]*account_service.Account, len(accounts))
	for i, acc := range accounts {
		accs[i] = accountDetailsFromAccount(acc)
	}

	return &account_service.ListAccountsResponse{
		Accounts:      accs,
		NextPageToken: next_token,
	}, err
}
