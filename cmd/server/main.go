package main

import (
	"log"
	"net"

	"github.com/sayshu-7s/grpc-gateway-example/gen/go/example"
	"github.com/sayshu-7s/grpc-gateway-example/server"
	"google.golang.org/grpc"
)

func main() {

	srv := grpc.NewServer()

	api, err := server.NewExampleAPIServer()
	if err != nil {
		log.Fatal("failed to new ExampleAPIServer")
	}
	example.RegisterExampleApiServer(srv, api)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen(tcp, :50051)")
	}
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("err has occured while serving: %v", err)
	}
}
