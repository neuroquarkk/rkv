package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rkv/internal/constants"
	"rkv/internal/dispatcher"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// maps endpoint specific grpc code to HTTP status
// code not present here fallbacks to 500 after logging
type codeMap map[codes.Code]int

var (
	putCodes = codeMap{
		codes.InvalidArgument: http.StatusBadRequest,
	}
	getCodes = codeMap{
		codes.InvalidArgument: http.StatusBadRequest,
		codes.NotFound:        http.StatusNotFound,
		codes.Internal:        http.StatusInternalServerError,
	}
	deleteCodes = codeMap{
		codes.InvalidArgument: http.StatusBadRequest,
		codes.NotFound:        http.StatusNotFound,
		codes.Internal:        http.StatusInternalServerError,
	}
	existsCodes = codeMap{
		codes.InvalidArgument: http.StatusBadRequest,
	}
	infoCodes = codeMap{}
)

func statusFor(err error, m codeMap) int {
	// order: no-members -> 503, deadline -> 408
	// endpoint specific grpc code -> mapped
	// anything else -> 500 + logged
	if errors.Is(err, dispatcher.ErrNoMembers) {
		return http.StatusServiceUnavailable
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusRequestTimeout
	}

	st, ok := status.FromError(err)
	if !ok {
		log.Printf("unexpected error: %v\n", err)
		return http.StatusInternalServerError
	}

	if code, mapped := m[st.Code()]; mapped {
		return code
	}

	log.Printf("grpc error: code=%s, msg=%v\n", st.Code(), st.Message())
	return http.StatusInternalServerError
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func validateKey(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if len(key) > constants.MaxKeySize {
		return fmt.Errorf("key exceeds max size of %d bytes",
			constants.MaxKeySize)
	}
	return nil
}

func validateValue(value string) error {
	if value == "" {
		return errors.New("value cannot be empty")
	}
	if len(value) > constants.MaxValueSize {
		return fmt.Errorf("value exceeds max size of %d bytes",
			constants.MaxValueSize)
	}
	return nil
}
