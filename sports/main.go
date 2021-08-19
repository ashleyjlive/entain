package main

import (
	"database/sql"
	"flag"
	"net"
	"os"
	"path/filepath"

	"github.com/ashleyjlive/entain/sports/db"
	"github.com/ashleyjlive/entain/sports/proto/sports"
	"github.com/ashleyjlive/entain/sports/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:10000", "gRPC server endpoint")
	dflt_db_path = filepath.Join(homeDir(), "entain", "sports", "data.db")
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
	conn, err := net.Listen("tcp", ":10000")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(*db_path), os.ModeDir)
	if err != nil {
		panic(err)
	}
	eventsDB, err := sql.Open("sqlite3", *db_path)
	if err != nil {
		return err
	}

	eventsRepo := db.NewEventsRepo(eventsDB)
	if err := eventsRepo.Init(*seed); err != nil {
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
