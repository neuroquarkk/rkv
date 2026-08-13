package main

import (
	"context"
	"log"
	"net/http"

	"rkv/internal/config"
	"rkv/internal/dispatcher"
	"rkv/internal/handler"
	"rkv/internal/middleware"
	"rkv/internal/registry/poller"
)

func main() {
	ctx := context.Background()
	cfg := config.NewRouter()

	addrsChan := make(chan []string, 1)
	disp := dispatcher.New()
	poller := poller.New(cfg.REGISTRY_URL)

	poller.Start(ctx, addrsChan)
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
		Addr:    ":" + cfg.PORT,
		Handler: middleware.Logger(mux),
	}

	log.Printf("router starting on PORT: %v\n", cfg.PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
