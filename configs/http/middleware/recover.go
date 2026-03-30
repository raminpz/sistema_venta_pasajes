package middleware

import (
	"net/http"
	"sistema_venta_pasajes/pkg"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				pkg.WriteError(w, r, pkg.Internal("ocurrió un error inesperado"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
