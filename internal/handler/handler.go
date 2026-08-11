package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"rkv/internal/dispatcher"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	client *dispatcher.Client
}

func New(client *dispatcher.Client) *Handler {
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

	err := h.client.Put(r.Context(), key, []byte(value))
	if err != nil {
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if len(key) > 256 || key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := h.client.Get(r.Context(), key)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Printf("unexpected status: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch st.Code() {
		case codes.NotFound:
			w.WriteHeader(http.StatusNotFound)
		case codes.Internal:
			w.WriteHeader(http.StatusInternalServerError)
		default:
			log.Printf("grpc error: code=%s, msg=%v\n", st.Code(), st.Message())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp := GetResp{Key: key, Data: string(data)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if len(key) > 256 || key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.client.Delete(r.Context(), key); err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Printf("unexpected status: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch st.Code() {
		case codes.NotFound:
			w.WriteHeader(http.StatusNotFound)
		case codes.Internal:
			w.WriteHeader(http.StatusInternalServerError)
		default:
			log.Printf("grpc error: code=%s, msg=%v\n", st.Code(), st.Message())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Exists(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if len(key) > 256 || key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := h.client.Exists(r.Context(), key)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Printf("unexpected status: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if st.Code() == codes.NotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("grpc error: code=%s, msg=%v\n", st.Code(), st.Message())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	name, err := h.client.Info(r.Context())
	if err != nil {
		log.Printf("unexpected error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &InfoResp{Name: name}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
