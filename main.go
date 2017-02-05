package main

import (
	"os"

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
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		dburl = defaultDb
	}

	pg := database.PostgreSQL{}
	err := pg.Connect(dburl)
	// defer pg.Close()
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	as := server.AccountServer{DB: pg}

	impl := func(g *grpc.Server) {
		account.RegisterAccountServiceServer(g, as)
	}

	err = lile.NewServer(
		lile.Name("account_service"),
		lile.Implementation(impl),
	).ListenAndServe()

	log.Fatal(err)
}
