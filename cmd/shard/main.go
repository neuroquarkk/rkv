package main

import (
	"log"
	"net"
	"os"

	pb "rkv/gen"
	"rkv/internal/service"
	"rkv/internal/store"

	"google.golang.org/grpc"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	lis, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("failed to open tcp port: %v\n", err)
	}

	st := store.New()
	service := service.New(st)

	server := grpc.NewServer()
	pb.RegisterShardServiceServer(server, service)

	log.Printf("service starting on PORT: %v\n", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v\n", err)
	}
}
