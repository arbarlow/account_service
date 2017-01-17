package main

import (
	"flag"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/arbarlow/account_service/account"
	"github.com/arbarlow/account_service/server"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var dbConnect string

var logger = log.WithFields(log.Fields{
	"service": "accounts",
})

func main() {
	defaultDb := "host=localhost user=postgres dbname=account_service sslmode=disable"
	flag.StringVar(&dbConnect, "connect", defaultDb, "db connection string")
	flag.Parse()

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	as := server.AccountServer{}
	db, err := as.DBConnect(dbConnect)
	defer db.Close()
	if err != nil {
		logger.Fatal("failed to connect: %v", err)
	}

	grpcServer := grpc.NewServer()
	account.RegisterAccountServiceServer(grpcServer, as)

	logger.Info("Listening on :8000")
	grpcServer.Serve(lis)
}
