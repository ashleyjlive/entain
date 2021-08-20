package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/ashleyjlive/entain/api/proto/racing"
	"github.com/ashleyjlive/entain/api/proto/sports"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	apiEndpoint = flag.String("api-endpoint", "localhost:8000",
		"API endpoint")
	racingGrpcEndpoint = flag.String("racing-grpc-endpoint", "localhost:9000",
		"gRPC server endpoint for racing service")
	sportsGrpcEndpoint = flag.String("sports-grpc-endpoint", "localhost:10000",
		"gRPC server endpoint for sports service")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running api server: %s", err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	if err := racing.RegisterRacingHandlerFromEndpoint(
		ctx,
		mux,
		*racingGrpcEndpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}

	if err := sports.RegisterSportsHandlerFromEndpoint(
		ctx,
		mux,
		*sportsGrpcEndpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}

	log.Infof("API server listening on: %s", *apiEndpoint)

	return http.ListenAndServe(*apiEndpoint, mux)
}
