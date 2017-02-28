package server

import (
	"errors"

	context "golang.org/x/net/context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lileio/account_service/account"
	"github.com/lileio/account_service/database"
)

type AccountServer struct {
	account.AccountServiceServer
	DB database.Database
}

var (
	ErrAuthFail  = errors.New("email or password incorrect")
	ErrNoAccount = errors.New("no account details provided")
)

func (as AccountServer) List(ctx context.Context, l *account.ListAccountsRequest) (*account.ListAccountsResponse, error) {
	accounts, next_token, err := as.DB.List(l.PageSize, l.PageToken)
	if err != nil {
		return nil, err
	}

	accs := make([]*account.Account, len(accounts))

	for i, acc := range accounts {
		accs[i] = accountDetailsFromAccount(acc)
	}

	return &account.ListAccountsResponse{
		Accounts:      accs,
		NextPageToken: next_token,
	}, err
}

func (as AccountServer) Create(ctx context.Context, r *account.CreateAccountRequest) (*account.Account, error) {
	if r.Account == nil {
		return nil, ErrNoAccount
	}

	a := database.NewAccount(r.Account.Name, r.Account.Email)
	err := as.DB.Create(&a, r.Password)
	if err != nil {
		return nil, err
	}

	return accountDetailsFromAccount(&a), nil
}

func (as AccountServer) GetById(ctx context.Context, r *account.GetByIdRequest) (*account.Account, error) {
	a, err := as.DB.ReadByID(r.Id)
	if err != nil {
		return nil, err
	}

	return accountDetailsFromAccount(a), nil
}

func (as AccountServer) GetByEmail(ctx context.Context, r *account.GetByEmailRequest) (*account.Account, error) {
	a, err := as.DB.ReadByEmail(r.Email)
	if err != nil {
		return nil, err
	}

	return accountDetailsFromAccount(a), nil
}

func (as AccountServer) AuthenticateByEmail(ctx context.Context, r *account.AuthenticateByEmailRequest) (*account.Account, error) {
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

func (as AccountServer) Update(ctx context.Context, r *account.UpdateAccountRequest) (*account.Account, error) {
	if r.Account == nil {
		return nil, ErrNoAccount
	}

	a := database.NewAccount(r.Account.Name, r.Account.Email)
	a.ID = r.Id

	err := as.DB.Update(&a)
	if err != nil {
		return nil, err
	}

	return accountDetailsFromAccount(&a), nil
}

func (as AccountServer) Delete(ctx context.Context, r *account.DeleteAccountRequest) (*empty.Empty, error) {
	err := as.DB.Delete(r.Id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func accountDetailsFromAccount(a *database.Account) *account.Account {
	return &account.Account{
		Id:    a.ID,
		Name:  a.Name,
		Email: a.Email,
	}
}
