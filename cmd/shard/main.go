package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"

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
	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	cfg := config.NewShard()

	lis, err := net.Listen("tcp", cfg.PORT)
	if err != nil {
		log.Fatalf("failed to open tcp port: %v\n", err)
	}

	st := store.New()
	rt := reporter.New(cfg.REGISTRY_URL, cfg.SHARD_ADDR)

	rt.Start(ctx, constants.ReporterTick)

	service := service.New(st, cfg.SHARD_ADDR)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingInterceptor),
	)
	pb.RegisterShardServiceServer(server, service)

	log.Printf("service starting on PORT: %v\n", lis.Addr())
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("failed to serve: %v\n", err)
			cancel()
		}
	}()

	<-ctx.Done()

	rt.Leave()
	server.GracefulStop()
}
