package server

import (
	"errors"
	"os"
	"time"

	"google.golang.org/grpc"

	context "golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	account "github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	is "github.com/lileio/image_service/image_service"
	opentracing "github.com/opentracing/opentracing-go"
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

func (as AccountServer) storeImage(
	ctx context.Context,
	img *is.ImageStoreRequest,
	a *database.Account) error {

	// Delete previous images if present
	if len(a.Images) > 0 {
		as.deleteImages(ctx, a)
	}

	res, err := imageService().StoreSync(ctx, img)
	if err != nil {
		return err
	}

	a.Images = res.Images
	return nil
}

func (as AccountServer) deleteImages(ctx context.Context, a *database.Account) error {
	for _, i := range a.Images {
		dr := is.DeleteRequest{Filename: i.Filename}
		_, err := imageService().Delete(ctx, &dr)
		if err != nil {
			return err
		}
	}

	return nil
}

func imageService() is.ImageServiceClient {
	if image_service != nil {
		return image_service
	}

	t := opentracing.GlobalTracer()

	conn, err := grpc.Dial(
		os.Getenv("IMAGE_SERVICE_ADDR"),
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(t)),
	)
	if err != nil {
		logrus.Warnf("Image service error: %s", err)
	}

	return is.NewImageServiceClient(conn)
}
