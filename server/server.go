package server

import (
	"errors"
	"os"
	"time"

	"google.golang.org/grpc"

	context "golang.org/x/net/context"

	"github.com/golang/protobuf/ptypes/empty"
	account "github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	is "github.com/lileio/image_service/image_service"
)

type AccountServer struct {
	account.AccountServiceServer
	DB database.Database
}

var (
	ErrAuthFail  = errors.New("email or password incorrect")
	ErrNoAccount = errors.New("no account details provided")

	image_service is.ImageServiceClient
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

	a := database.Account{
		Name:  r.Account.Name,
		Email: r.Account.Email,
	}

	if r.Image != nil {
		as.storeImage(ctx, r.Image, &a)
	}

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

	a := database.Account{
		ID:    r.Id,
		Name:  r.Account.Name,
		Email: r.Account.Email,
	}

	if r.Image != nil {
		as.storeImage(ctx, r.Image, &a)
	}

	err := as.DB.Update(&a)
	if err != nil {
		return nil, err
	}

	return accountDetailsFromAccount(&a), nil
}

func (as AccountServer) Delete(ctx context.Context, r *account.DeleteAccountRequest) (*empty.Empty, error) {
	a := database.Account{ID: r.Id}
	err := as.deleteImages(ctx, &a)
	if err != nil {
		return nil, err
	}

	err = as.DB.Delete(r.Id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (as AccountServer) storeImage(
	ctx context.Context,
	img *is.ImageStoreRequest,
	a *database.Account) error {

	if a.ID != "" {
		as.deleteImages(ctx, a)
	}

	c, err := imageService()
	if err != nil {
		return err
	}

	res, err := c.StoreSync(ctx, img)
	if err != nil {
		return err
	}

	a.Images = res.Images
	return nil
}

func (as AccountServer) deleteImages(ctx context.Context, a *database.Account) error {
	c, err := imageService()
	if err != nil {
		return err
	}

	acc, err := as.DB.ReadByID(a.ID)
	if err != nil {
		return err
	}

	for _, i := range acc.Images {
		dr := is.DeleteRequest{Filename: i.Filename}
		_, err := c.Delete(ctx, &dr)
		if err != nil {
			return err
		}
	}

	return nil
}

func accountDetailsFromAccount(a *database.Account) *account.Account {
	imgs := map[string]*is.Image{}
	for _, i := range a.Images {
		imgs[i.VersionName] = i
	}

	return &account.Account{
		Id:     a.ID,
		Name:   a.Name,
		Email:  a.Email,
		Images: imgs,
	}
}

func imageService() (is.ImageServiceClient, error) {
	if image_service != nil {
		return image_service, nil
	}

	conn, err := grpc.Dial(
		os.Getenv("IMAGE_SERVICE_ADDR"),
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second),
	)
	if err != nil {
		return nil, err
	}

	image_service = is.NewImageServiceClient(conn)
	return image_service, nil
}
