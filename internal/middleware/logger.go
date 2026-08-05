package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type responseWrite struct {
	http.ResponseWriter
	code int
}

func (rw *responseWrite) WriteHeader(code int) {
	rw.code = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWrite{ResponseWriter: w, code: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		timestamp := time.Now().Format("2006/01/02 15:04:05")

		fmt.Printf("[LOG] %s | %3d | %10v | %-5s | %s\n",
			timestamp,
			rw.code,
			duration,
			r.Method,
			r.URL.Path,
		)
	})
}
