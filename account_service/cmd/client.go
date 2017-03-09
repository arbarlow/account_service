package cmd

import (
	"log"
	"time"

	"github.com/lileio/account_service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var addr string

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "a basic cmd line client for account service",
}

var client = func() account_service.AccountServiceClient {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second),
	)

	if err != nil {
		log.Fatal(err)
	}

	return account_service.NewAccountServiceClient(conn)
}

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().StringVarP(&addr, "addr", "a", "localhost:8000", "address for service. i.e localhost:8001")
}
