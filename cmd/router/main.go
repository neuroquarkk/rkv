package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"rkv/internal/dispatcher"
	"rkv/internal/handler"
	"rkv/internal/middleware"
	"rkv/internal/registry"
)

func main() {
	ctx := context.Background()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8081"
	}

	memberAddrs, ok := os.LookupEnv("MEMBER_ADDR")
	if !ok {
		memberAddrs = "localhost:8080,localhost:8082"
	}

	reg, err := registry.New(ctx, memberAddrs)
	if err != nil {
		log.Fatalf("failed to create new pool: %v\n", err)
	}
	defer reg.Close()
	log.Printf("%d members connected\n", len(reg.Clients))

	disp := dispatcher.New(reg)

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
