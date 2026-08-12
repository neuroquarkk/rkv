package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"rkv/internal/registry/server"
)

func main() {
	ctx := context.Background()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8010"
	}

	rs := server.New(15 * time.Second)
	rs.Sweeper(ctx, 5*time.Second)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /heartbeat", rs.Heartbeat)
	mux.HandleFunc("GET /members", rs.Members)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("starting registry server on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start registry server: %v\n", err)
	}
}
