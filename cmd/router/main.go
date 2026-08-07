package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"rkv/internal/client"
	"rkv/internal/handler"
	"rkv/internal/middleware"
)

func main() {
	ctx := context.Background()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8081"
	}

	memberAddr, ok := os.LookupEnv("MEMBER_ADDR")
	if !ok {
		memberAddr = "localhost:8080"
	}

	client, err := client.New(ctx, memberAddr)
	if err != nil {
		log.Fatalf("could not connect to shard: %v\n", err)
	}
	defer client.Close()
	log.Println("shard client connected successfully")

	handler := handler.New(client)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", handler.Ping)
	mux.HandleFunc("PUT /put/{key}", handler.Put)
	mux.HandleFunc("GET /get/{key}", handler.Get)
	mux.HandleFunc("DELETE /delete/{key}", handler.Delete)
	mux.HandleFunc("HEAD /key/{key}", handler.Exists)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: middleware.Logger(mux),
	}

	log.Printf("router starting on PORT: %v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
