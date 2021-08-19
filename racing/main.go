package main

import (
	"database/sql"
	"flag"
	"net"
	"os"
	"path/filepath"

	"github.com/ashleyjlive/entain/racing/db"
	"github.com/ashleyjlive/entain/racing/proto/racing"
	"github.com/ashleyjlive/entain/racing/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9000", "gRPC server endpoint")
	dflt_db_path = filepath.Join(homeDir(), "entain", "racing", "data.db")
	db_path      = flag.String("db_path", dflt_db_path, "The path of the database.")
	seed         = flag.Bool("seed", false, "Determines if sample data is to be inserted into the database")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s", err)
	}
}

func homeDir() string {
	osPath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return osPath
}

func run() error {
	conn, err := net.Listen("tcp", ":9000")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(*db_path), os.ModeDir)
	if err != nil {
		panic(err)
	}
	racingDB, err := sql.Open("sqlite3", *db_path)
	if err != nil {
		return err
	}

	racesRepo := db.NewRacesRepo(racingDB)
	if err := racesRepo.Init(*seed); err != nil {
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
