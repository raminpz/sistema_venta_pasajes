package middleware

import (
	"net/http"
	"sistema_venta_pasajes/pkg"

	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		w.Header().Set(requestIDHeader, requestID)
		r = r.WithContext(pkg.WithRequestID(r.Context(), requestID))

		next.ServeHTTP(w, r)
	})
}
