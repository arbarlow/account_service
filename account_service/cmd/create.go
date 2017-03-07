package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lileio/account_service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var addr string
var name string
var email string
var password string
var image string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create an account",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(
			addr,
			grpc.WithInsecure(),
			grpc.WithTimeout(1*time.Second),
		)

		if err != nil {
			log.Fatal(err)
		}

		client := account_service.NewAccountServiceClient(conn)

		ar := &account_service.CreateAccountRequest{
			Account: &account_service.Account{
				Name:  name,
				Email: email,
			},
			Password: password,
		}

		ctx := context.Background()
		res, err := client.Create(ctx, ar)
		if err != nil {
			log.Fatal(err)
		}

		js, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(js))
	},
}

func init() {
	clientCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&addr, "addr", "a", "localhost:8000", "address for service. i.e localhost:8001")
	createCmd.Flags().StringVarP(&name, "name", "n", "", "name for account")
	createCmd.Flags().StringVarP(&email, "email", "e", "", "email address")
	createCmd.Flags().StringVarP(&password, "password", "p", "", "password for account")
	createCmd.Flags().StringVarP(&image, "image", "i", "", "path to image file")
}
