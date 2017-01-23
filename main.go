package main

import (
	"net"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/arbarlow/account_service/account"
	"github.com/arbarlow/account_service/server"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
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

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	account.RegisterAccountServiceServer(grpcServer, as)
	grpc_prometheus.Register(grpcServer)

	logger.Info("Listening: gRPC:8000, Prometheus:8080")

	http.Handle("/metrics", prometheus.Handler())
	go http.ListenAndServe(":8080", nil)
	grpcServer.Serve(lis)
}
