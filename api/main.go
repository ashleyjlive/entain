package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"git.neds.sh/matty/entain/api/proto/racing"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	apiEndpoint  = flag.String("api-endpoint", "localhost:8000", "API endpoint")
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9000", "gRPC server endpoint")
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
		*grpcEndpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}

	log.Infof("API server listening on: %s", *apiEndpoint)

	return http.ListenAndServe(*apiEndpoint, mux)
}
