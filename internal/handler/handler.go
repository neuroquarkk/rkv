package handler

import (
	"encoding/json"
	"net/http"
	"rkv/internal/dispatcher"
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
		w.WriteHeader(statusFor(err, putCodes))
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
		w.WriteHeader(statusFor(err, getCodes))
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

	err := h.client.Delete(r.Context(), key)
	if err != nil {
		w.WriteHeader(statusFor(err, deleteCodes))
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
		w.WriteHeader(statusFor(err, existsCodes))
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	names, err := h.client.Info(r.Context())
	if err != nil {
		w.WriteHeader(statusFor(err, infoCodes))
		return
	}

	resp := &InfoResp{Names: names}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
