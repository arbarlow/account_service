package main

import (
	"net"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/arbarlow/account_service/account"
	"github.com/arbarlow/account_service/server"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var defaultDb = "postgres://postgres@localhost/account_service?sslmode=disable"

var logger = log.WithFields(log.Fields{
	"service": "accounts",
})

func main() {
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		dburl = defaultDb
	}

	log.SetLevel(log.InfoLevel)

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	as := server.AccountServer{}
	db, err := as.DBConnect(dburl)
	defer db.Close()
	if err != nil {
		logger.Fatalf("failed to connect: %v", err)
	}

	grpcServer := grpc.NewServer()
	account.RegisterAccountServiceServer(grpcServer, as)

	logger.Info("Listening on :8000")
	grpcServer.Serve(lis)
}
