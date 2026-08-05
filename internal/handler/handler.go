package handler

import (
	"net/http"
	pb "rkv/gen"
)

type Handler struct {
	client pb.ShardServiceClient
}

func New(client pb.ShardServiceClient) *Handler {
	return &Handler{client}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}
