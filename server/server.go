package server

import (
	context "golang.org/x/net/context"

	"github.com/lileio/account_service/account"
	"github.com/lileio/account_service/database"
)

type AccountServer struct {
	account.AccountServiceServer
	DB database.Database
}

func (as AccountServer) Create(ctx context.Context, r *account.AccountCreateRequest) (*account.AccountDetails, error) {
	a := database.NewAccount(r.Name, r.Email)
	err := as.DB.Create(&a)
	if err != nil {
		return nil, err
	}

	return &account.AccountDetails{
		Id:    a.Id,
		Name:  a.Name.String,
		Email: a.Email.String,
	}, nil
}

func (as AccountServer) Read(ctx context.Context, r *account.AccountRequest) (*account.AccountDetails, error) {
	a, err := as.DB.ReadByID(r.Id)
	if err != nil {
		return nil, err
	}

	return &account.AccountDetails{
		Id:    a.Id,
		Name:  a.Name.String,
		Email: a.Email.String,
	}, nil
}

func (as AccountServer) Update(ctx context.Context, r *account.AccountDetails) (*account.AccountDetails, error) {
	a := database.NewAccount(r.Name, r.Email)
	a.Id = r.Id

	err := as.DB.Update(&a)
	if err != nil {
		return nil, err
	}

	return &account.AccountDetails{
		Id:    a.Id,
		Name:  a.Name.String,
		Email: a.Email.String,
	}, nil
}

func (as AccountServer) Delete(ctx context.Context, r *account.AccountDeleteRequest) (*account.AccountDeleteResponse, error) {
	err := as.DB.Delete(r.Id)
	if err != nil {
		return nil, err
	}

	return &account.AccountDeleteResponse{Id: r.Id}, nil
}
