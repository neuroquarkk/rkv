package main

import (
	"context"
	"log"
	"net"

	pb "rkv/gen"
	"rkv/internal/config"
	"rkv/internal/constants"
	"rkv/internal/interceptor"
	"rkv/internal/registry/reporter"
	"rkv/internal/service"
	"rkv/internal/store"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	cfg := config.NewShard()

	lis, err := net.Listen("tcp", cfg.PORT)
	if err != nil {
		log.Fatalf("failed to open tcp port: %v\n", err)
	}

	st := store.New()
	reporter.Start(
		ctx, cfg.REGISTRY_URL, cfg.SHARD_ADDR, constants.ReporterTick,
	)

	service := service.New(st, cfg.SHARD_ADDR)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingInterceptor),
	)
	pb.RegisterShardServiceServer(server, service)

	log.Printf("service starting on PORT: %v\n", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v\n", err)
	}
}
