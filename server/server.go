package server

import (
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	context "golang.org/x/net/context"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	account "github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	"github.com/lileio/image_service"
	"github.com/lileio/lile"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type AccountServer struct {
	account.AccountServiceServer
	DB database.Database
}

var (
	is image_service.ImageServiceClient

	ErrNoAccount = grpc.Errorf(codes.InvalidArgument, "account is nil")
)

func NewAccountServer() *lile.Server {
	db := database.DatabaseFromEnv()
	db.Migrate()
	defer db.Close()

	as := AccountServer{DB: db}

	impl := func(g *grpc.Server) {
		account.RegisterAccountServiceServer(g, as)
	}

	return lile.NewServer(
		lile.Name("account_service"),
		lile.Implementation(impl),
		lile.Publishers(map[string]string{
			"Create":                "account_service.created",
			"Update":                "account_service.updated",
			"Delete":                "account_service.deleted",
			"GeneratePasswordToken": "account_service.password_token_generated",
			"ResetPassword":         "account_service.password_reset",
			"ConfirmAccount":        "account_service.account_confirmed",
		}),
	)
}

func accountDetailsFromAccount(a *database.Account) *account.Account {
	imgs := map[string]*image_service.Image{}
	for _, i := range a.Images {
		imgs[i.VersionName] = i
	}

	return &account.Account{
		Id:                 a.ID,
		Name:               a.Name,
		Email:              a.Email,
		Images:             imgs,
		Metadata:           a.Metadata,
		ConfirmToken:       a.ConfirmationToken,
		PasswordResetToken: a.PasswordResetToken,
	}
}

func (as AccountServer) storeImage(
	ctx context.Context,
	img *image_service.ImageStoreRequest,
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
		dr := image_service.DeleteRequest{Filename: i.Filename}
		_, err := imageService().Delete(ctx, &dr)
		if err != nil {
			return err
		}
	}

	return nil
}

func imageService() image_service.ImageServiceClient {
	if is != nil {
		return is
	}

	addr := os.Getenv("IMAGE_SERVICE_ADDR")
	if addr == "" {
		addr = "image_service"
	}

	t := opentracing.GlobalTracer()

	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(t)),
	)
	if err != nil {
		logrus.Warnf("image service connection error: %s", err)
	}

	is = image_service.NewImageServiceClient(conn)
	return is

}
