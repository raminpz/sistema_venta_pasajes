package middleware

import (
	"net/http"

	"sistema_venta_pasajes/pkg"
)

// RequireRole devuelve un middleware que verifica que el usuario autenticado
// tenga uno de los roles permitidos. Debe aplicarse después de JWTAuth.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetJWTClaims(r.Context())
			if claims == nil {
				pkg.WriteError(w, r, pkg.Unauthorized(
					"no_autenticado",
					"No autenticado. Por favor inicie sesión.",
				))
				return
			}

			if _, ok := allowed[claims.Rol]; !ok {
				pkg.WriteError(w, r, pkg.Forbidden(
					"acceso_denegado",
					"No tienes permisos para realizar esta acción.",
				))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
