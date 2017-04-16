package cmd

import (
	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	"github.com/lileio/account_service/server"
	"github.com/lileio/account_service/subscribers"
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
		defer db.Close()

		asub := subscribers.AccountServiceSubscribers{DB: db}

		as := server.AccountServer{DB: db}
		impl := func(g *grpc.Server) {
			account_service.RegisterAccountServiceServer(g, as)
		}

		s := lile.NewServer(
			lile.Name("account_service"),
			lile.Implementation(impl),
			lile.Publishers(map[string]string{
				"Create": "account_service.created",
			}),
			lile.Subscriber(&asub),
		)

		logrus.Fatal(s.ListenAndServe())
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
