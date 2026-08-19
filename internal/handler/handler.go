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
	if err := validateKey(key); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	var data ValReq
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	value := data.Value
	if err := validateValue(value); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.client.Put(r.Context(), key, []byte(value))
	if err != nil {
		writeJSONError(w, statusFor(err, putCodes), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.client.Get(r.Context(), key)
	if err != nil {
		writeJSONError(w, statusFor(err, getCodes), err.Error())
		return
	}

	resp := GetResp{Key: key, Data: string(data)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.client.Delete(r.Context(), key)
	if err != nil {
		writeJSONError(w, statusFor(err, deleteCodes), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Exists(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := validateKey(key); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	exists, err := h.client.Exists(r.Context(), key)
	if err != nil {
		writeJSONError(w, statusFor(err, existsCodes), err.Error())
		return
	}

	code := http.StatusNoContent
	if !exists {
		code = http.StatusNotFound
	}
	w.WriteHeader(code)
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	names, err := h.client.Info(r.Context())
	if err != nil {
		writeJSONError(w, statusFor(err, infoCodes), err.Error())
		return
	}

	resp := &InfoResp{Names: names}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
