package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "rkv/gen"
	"rkv/internal/interceptor"
	"rkv/internal/registry/reporter"
	"rkv/internal/service"
	"rkv/internal/store"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	registryUrl, ok := os.LookupEnv("REGISTRY_URL")
	if !ok {
		registryUrl = "localhost:8010"
	}

	shardAddr, ok := os.LookupEnv("SHARD_ADDR")
	if !ok {
		shardAddr = "localhost:" + port
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to open tcp port: %v\n", err)
	}

	st := store.New()
	reporter.Start(ctx, registryUrl, shardAddr, 5*time.Second)

	service := service.New(st, shardAddr)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingInterceptor),
	)
	pb.RegisterShardServiceServer(server, service)

	log.Printf("service starting on PORT: %v\n", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v\n", err)
	}
}
