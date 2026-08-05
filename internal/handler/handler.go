package handler

import (
	"net/http"
	"rkv/internal/client"
)

type Handler struct {
	client *client.Client
}

func New(client *client.Client) *Handler {
	return &Handler{client}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}
