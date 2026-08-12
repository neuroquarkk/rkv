package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"rkv/internal/dispatcher"
	"rkv/internal/handler"
	"rkv/internal/middleware"
	"rkv/internal/registry/poller"
)

func main() {
	ctx := context.Background()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8081"
	}

	registryUrl, ok := os.LookupEnv("REGISTRY_URL")
	if !ok {
		registryUrl = "localhost:8010"
	}

	addrsChan := make(chan []string, 1)
	disp := dispatcher.New()

	poller.Start(ctx, registryUrl, addrsChan)
	disp.Start(ctx, addrsChan)

	handler := handler.New(disp)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", handler.Ping)
	mux.HandleFunc("PUT /put/{key}", handler.Put)
	mux.HandleFunc("GET /get/{key}", handler.Get)
	mux.HandleFunc("DELETE /delete/{key}", handler.Delete)
	mux.HandleFunc("HEAD /key/{key}", handler.Exists)
	mux.HandleFunc("GET /info", handler.Info)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: middleware.Logger(mux),
	}

	log.Printf("router starting on PORT: %v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
