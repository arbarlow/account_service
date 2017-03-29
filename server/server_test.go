package server

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	account "github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	"github.com/lileio/image_service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var db = setupDB()
var as = AccountServer{}

type MockImageService struct {
	mock.Mock
}

func (m MockImageService) Store(ctx context.Context, in *image_service.ImageStoreRequest, opts ...grpc.CallOption) (image_service.ImageService_StoreClient, error) {
	m.Called()
	return nil, nil
}

func (m MockImageService) StoreSync(ctx context.Context, in *image_service.ImageStoreRequest, opts ...grpc.CallOption) (*image_service.ImageSyncResponse, error) {
	m.Called()
	return &image_service.ImageSyncResponse{Images: []*image_service.Image{
		{Filename: in.Filename},
	}}, nil
}

func (m MockImageService) Delete(ctx context.Context, in *image_service.DeleteRequest, opts ...grpc.CallOption) (*image_service.DeleteResponse, error) {
	m.Called()
	return nil, nil
}

func setupDB() database.Database {
	conn := database.DatabaseFromEnv()
	as.DB = conn
	as.DB.Migrate()
	return conn
}

func truncate() {
	db.Truncate()
}

var name = "Alex B"
var email = "alexb@localhost"
var pass = "password"

func createAccount(t *testing.T) *account.Account {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx := context.Background()
	req := &account.CreateAccountRequest{
		Account: &account.Account{
			Name:     name,
			Email:    "alexbarlowimage_service.localhost" + strconv.Itoa(r.Int()),
			Metadata: map[string]string{"test": "test"},
		},
		Password: pass,
	}
	account, err := as.Create(ctx, req)
	assert.Nil(t, err)
	return account
}
