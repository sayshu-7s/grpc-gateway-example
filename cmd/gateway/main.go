package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sayshu-7s/grpc-gateway-example/gen/go/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const grpcServerAddress = "grpc-server:50051"
const docsServerAddress = "http://docs-server:8080"

func main() {
	grpcGateway := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := example.RegisterExampleApiHandlerFromEndpoint(context.Background(), grpcGateway, grpcServerAddress, opts); err != nil {
		log.Fatal("failed to register grpc-server")
	}

	docsURL, err := url.Parse(docsServerAddress)
	if err != nil {
		log.Fatalf("failed to parse docsServerURL=%v", docsServerAddress)
	}
	docsProxy := httputil.NewSingleHostReverseProxy(docsURL)

	mux := http.NewServeMux()
	mux.Handle("/docs/", docsProxy)
	mux.Handle("/", grpcGateway)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("err")
	}
}
