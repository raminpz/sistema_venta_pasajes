package middleware

import (
	"context"
	"log"
	"net/http"
	"sistema_venta_pasajes/pkg"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// claimsContextKey es el tipo de clave para el contexto, evita colisiones.
type claimsContextKey string

const jwtClaimsKey claimsContextKey = "jwt_claims"

// JWTClaims representa el payload del access token.
type JWTClaims struct {
	IDUsuario int    `json:"id_usuario"`
	Email     string `json:"email"`
	Rol       string `json:"rol"`
	jwt.RegisteredClaims
}

// DevBypassAuth es un middleware que omite la validacion JWT en entorno de desarrollo.
// Inyecta claims ficticios con rol PROVEEDOR (acceso total) para que los middlewares
// de roles no bloqueen las rutas.
// NUNCA activar AUTH_DISABLED=true en produccion.
func DevBypassAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[DEV] AUTH_DISABLED=true - autenticacion omitida, rol=PROVEEDOR inyectado")
		claims := &JWTClaims{
			IDUsuario: 0,
			Email:     "dev@local",
			Rol:       "PROVEEDOR",
		}
		ctx := context.WithValue(r.Context(), jwtClaimsKey, claims)
		ctx = pkg.WithActorUserID(ctx, int64(claims.IDUsuario))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JWTAuth valida el header Authorization: Bearer <token> e inyecta los claims en el contexto.
// Devuelve 401 si el token falta o es invalido/expirado.
func JWTAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				pkg.WriteError(w, r, pkg.Unauthorized(
					"token_requerido",
					"Se requiere token de autenticacion en el encabezado Authorization: Bearer <token>.",
				))
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := parseJWT(tokenStr, jwtSecret)
			if err != nil {
				pkg.WriteError(w, r, pkg.Unauthorized(
					"token_invalido",
					"Token invalido o expirado. Inicie sesion nuevamente.",
				))
				return
			}
			ctx := context.WithValue(r.Context(), jwtClaimsKey, claims)
			ctx = pkg.WithActorUserID(ctx, int64(claims.IDUsuario))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetJWTClaims extrae los claims JWT del contexto.
// Retorna nil si no hay claims (ruta publica o token no validado).
func GetJWTClaims(ctx context.Context) *JWTClaims {
	claims, _ := ctx.Value(jwtClaimsKey).(*JWTClaims)
	return claims
}

// parseJWT valida la firma y expiracion del token y retorna los claims.
func parseJWT(tokenStr, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pkg.Unauthorized("token_invalido", "Algoritmo de firma no esperado.")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, pkg.Unauthorized("token_invalido", "Token invalido o expirado.")
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, pkg.Unauthorized("token_invalido", "No se pudieron extraer los claims del token.")
	}
	return claims, nil
}
