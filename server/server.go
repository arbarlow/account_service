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

func (as AccountServer) Create(ctx context.Context, r *account.CreateRequest) (*account.Account, error) {
	a := database.NewAccount(r.Name, r.Email)
	err := as.DB.Create(&a)
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

func (as AccountServer) Update(ctx context.Context, r *account.Account) (*account.Account, error) {
	a := database.NewAccount(r.Name, r.Email)
	a.ID = r.Id

	err := as.DB.Update(&a)
	if err != nil {
		return nil, err
	}

	return accountDetailsFromAccount(&a), nil
}

func (as AccountServer) Delete(ctx context.Context, r *account.DeleteRequest) (*account.DeleteResponse, error) {
	err := as.DB.Delete(r.Id)
	if err != nil {
		return nil, err
	}

	return &account.DeleteResponse{Id: r.Id}, nil
}

func accountDetailsFromAccount(a *database.Account) *account.Account {
	return &account.Account{
		Id:    a.ID,
		Name:  a.Name,
		Email: a.Email,
	}
}
