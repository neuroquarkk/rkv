package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"rkv/internal/client"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if len(key) > 256 || key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var data ValReq
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	value := data.Value
	if len(value) > 1*1024 || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.client.Put(r.Context(), key, []byte(value)); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Printf("unexpected status: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if st.Code() == codes.InvalidArgument {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Printf("grpc error: code=%s, msg=%v\n", st.Code(), st.Message())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
