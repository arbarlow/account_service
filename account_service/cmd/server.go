package cmd

import (
	"github.com/lileio/account_service/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Fatal(server.NewAccountServer().
			ListenAndServe())
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
