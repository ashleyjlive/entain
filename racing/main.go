package main

import (
	"database/sql"
	"flag"
	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"git.neds.sh/matty/entain/racing/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9000", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9000")
	if err != nil {
		return err
	}

	racingDB, err := sql.Open("sqlite3", "./db/racing.db")
	if err != nil {
		return err
	}

	racesRepo := db.NewRacesRepo(racingDB)
	if err := racesRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	racing.RegisterRacingServer(
		grpcServer,
		service.NewRacingService(
			racesRepo,
		),
	)

	log.Infof("gRPC server listening on: %s", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
