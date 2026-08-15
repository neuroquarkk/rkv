package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
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
