package cmd

import "github.com/spf13/cobra"

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "a basic cmd line client for account service",
}

func init() {
	RootCmd.AddCommand(clientCmd)
}
