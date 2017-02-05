package main

import (
	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/lileio/account_service/account"
	"github.com/lileio/account_service/database"
	"github.com/lileio/account_service/server"
	"github.com/lileio/lile"
)

var defaultDb = "postgres://postgres@localhost/account_service?sslmode=disable"

func main() {
	db := database.DatabaseFromEnv()
	defer db.Close()
	as := server.AccountServer{DB: db}

	impl := func(g *grpc.Server) {
		account.RegisterAccountServiceServer(g, as)
	}

	err := lile.NewServer(
		lile.Name("account_service"),
		lile.Implementation(impl),
	).ListenAndServe()

	log.Fatal(err)
}
