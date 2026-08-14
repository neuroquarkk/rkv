package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"rkv/internal/config"
	"rkv/internal/dispatcher"
	"rkv/internal/handler"
	"rkv/internal/middleware"
	"rkv/internal/registry/poller"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	cfg := config.NewRouter()

	addrsChan := make(chan []string, 1)

	poller := poller.New(cfg.REGISTRY_URL)
	poller.Start(ctx, addrsChan)

	disp := dispatcher.New()
	defer disp.Close()

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
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("failed to serve: %v\n", err)
			cancel()
		}
	}()

	<-ctx.Done()

	server.Shutdown(context.TODO())
}
