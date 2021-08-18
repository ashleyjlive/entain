package main

import (
	"database/sql"
	"flag"
	"net"

	"github.com/ashleyjlive/entain/sports/db"
	"github.com/ashleyjlive/entain/sports/proto/sports"
	"github.com/ashleyjlive/entain/sports/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:10000", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":10000")
	if err != nil {
		return err
	}

	eventsDB, err := sql.Open("sqlite3", "./db/events.db")
	if err != nil {
		return err
	}

	eventsRepo := db.NewEventsRepo(eventsDB)
	if err := eventsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	sports.RegisterSportsServer(
		grpcServer,
		service.NewSportsService(
			eventsRepo,
		),
	)

	log.Infof("gRPC server listening on: %s", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
