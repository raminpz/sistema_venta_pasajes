package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"sistema_venta_pasajes/configs/http/middleware"
)

const jwtTestSecret = "secreto_de_prueba_para_tests_jwt_256bits"

var okHandlerJWT = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

// generarToken crea un JWT válido para tests.
func generarToken(secret, rol string, expiracion time.Duration) string {
	claims := &middleware.JWTClaims{
		IDUsuario: 1,
		Email:     "test@test.com",
		Rol:       rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiracion)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

// ─── JWTAuth ─────────────────────────────────────────────────────────────────

func TestJWTAuth_TokenValido(t *testing.T) {
	tokenStr := generarToken(jwtTestSecret, "ADMIN", 15*time.Minute)
	mw := middleware.JWTAuth(jwtTestSecret)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	w := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTAuth_SinHeader(t *testing.T) {
	mw := middleware.JWTAuth(jwtTestSecret)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "token_requerido")
}

func TestJWTAuth_HeaderSinBearer(t *testing.T) {
	mw := middleware.JWTAuth(jwtTestSecret)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "token-sin-prefijo")
	w := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_TokenExpirado(t *testing.T) {
	tokenStr := generarToken(jwtTestSecret, "ADMIN", -1*time.Hour) // ya expiró
	mw := middleware.JWTAuth(jwtTestSecret)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	w := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "token_invalido")
}

func TestJWTAuth_TokenFirmaDistinta(t *testing.T) {
	// Token firmado con secreto diferente
	tokenStr := generarToken("otro-secreto-diferente", "ADMIN", 15*time.Minute)
	mw := middleware.JWTAuth(jwtTestSecret)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	w := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuth_TokenMalformado(t *testing.T) {
	mw := middleware.JWTAuth(jwtTestSecret)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer esto.no.es.un.jwt.valido")
	w := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── GetJWTClaims ────────────────────────────────────────────────────────────

func TestGetJWTClaims_ExtracteClaims(t *testing.T) {
	tokenStr := generarToken(jwtTestSecret, "VENDEDOR", 15*time.Minute)

	var capturedClaims *middleware.JWTClaims
	handlerCaptura := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedClaims = middleware.GetJWTClaims(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.JWTAuth(jwtTestSecret)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	w := httptest.NewRecorder()
	mw(handlerCaptura).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedClaims)
	assert.Equal(t, 1, capturedClaims.IDUsuario)
	assert.Equal(t, "VENDEDOR", capturedClaims.Rol)
}
