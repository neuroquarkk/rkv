package main

import (
	"context"
	"log"
	"net/http"

	"rkv/internal/config"
	"rkv/internal/constants"
	"rkv/internal/registry/server"
)

func main() {
	ctx := context.Background()
	cfg := config.NewRegistry()

	rs := server.New(constants.StaleInterval)
	rs.Sweeper(ctx, constants.SweeperTick)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /heartbeat", rs.Heartbeat)
	mux.HandleFunc("GET /members", rs.Members)

	server := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: mux,
	}

	log.Printf("starting registry server on port %s...\n", cfg.PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start registry server: %v\n", err)
	}
}
