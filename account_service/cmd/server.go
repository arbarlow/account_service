package cmd

import (
	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	"github.com/lileio/account_service/server"
	"github.com/lileio/lile"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		db := database.DatabaseFromEnv()
		db.Migrate()
		defer db.Close()

		as := server.AccountServer{DB: db}

		impl := func(g *grpc.Server) {
			account_service.RegisterAccountServiceServer(g, as)
		}

		err := lile.NewServer(
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
		).ListenAndServe()

		logrus.Fatal(err)
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
