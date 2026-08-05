package main

import (
	"log"
	"net/http"
	"os"
	"rkv/internal/client"
	"rkv/internal/handler"
	"rkv/internal/middleware"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8081"
	}

	client, err := client.New("http://localhost:8080")
	if err != nil {
		log.Fatalf("could not connect to shard: %v\n", err)
	}
	defer client.Close()
	log.Println("shard client connected successfully")

	handler := handler.New(client.Client)

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", handler.Ping)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: middleware.Logger(mux),
	}

	log.Printf("router starting on PORT: %v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
