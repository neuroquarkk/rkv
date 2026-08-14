package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"rkv/internal/config"
	"rkv/internal/constants"
	"rkv/internal/registry/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	cfg := config.NewRegistry()

	rs := server.New(constants.StaleInterval)
	rs.Sweeper(ctx, constants.SweeperTick)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /heartbeat", rs.Heartbeat)
	mux.HandleFunc("POST /leave", rs.Leave)
	mux.HandleFunc("GET /members", rs.Members)

	server := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: mux,
	}

	log.Printf("starting registry server on port %s...\n", cfg.PORT)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("failed to start registry server: %v\n", err)
			cancel()
		}
	}()

	<-ctx.Done()

	server.Shutdown(context.TODO())
}
