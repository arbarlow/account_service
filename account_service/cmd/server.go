package cmd

import (
	"log"

	"github.com/lileio/account_service"
	"github.com/lileio/account_service/database"
	"github.com/lileio/account_service/server"
	"github.com/lileio/lile"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		db := database.DatabaseFromEnv()
		defer db.Close()

		as := server.AccountServer{DB: db}

		impl := func(g *grpc.Server) {
			account_service.RegisterAccountServiceServer(g, as)
		}

		err := lile.NewServer(
			lile.Name("account_service"),
			lile.Implementation(impl),
		).ListenAndServe()

		log.Fatal(err)
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
