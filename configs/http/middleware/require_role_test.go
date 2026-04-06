package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"sistema_venta_pasajes/configs/http/middleware"
)

// aplicarConRol inyecta claims en el contexto y aplica RequireRole.
func aplicarConRol(rol string, rolesPermitidos ...string) *httptest.ResponseRecorder {
	// Generar token real y pasar por JWTAuth para tener claims en el contexto
	tokenStr := generarToken(jwtTestSecret, rol, 15*time.Minute)

	var w *httptest.ResponseRecorder
	handler := middleware.RequireRole(rolesPermitidos...)(okHandlerJWT)
	jwtMw := middleware.JWTAuth(jwtTestSecret)

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	w = httptest.NewRecorder()
	jwtMw(handler).ServeHTTP(w, req)
	return w
}

// ─── RequireRole ─────────────────────────────────────────────────────────────

func TestRequireRole_RolPermitido(t *testing.T) {
	w := aplicarConRol("ADMIN", "ADMIN", "PROVEEDOR")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_ProveedorAccedeTodo(t *testing.T) {
	w := aplicarConRol("PROVEEDOR", "ADMIN", "VENDEDOR", "PROVEEDOR")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_VendedorAccedeOperaciones(t *testing.T) {
	w := aplicarConRol("VENDEDOR", "ADMIN", "VENDEDOR", "PROVEEDOR")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_VendedorNoAccedeDatosMaestros(t *testing.T) {
	w := aplicarConRol("VENDEDOR", "ADMIN", "PROVEEDOR") // VENDEDOR no está
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "acceso_denegado")
}

func TestRequireRole_SinClaims_Unauthorized(t *testing.T) {
	// Sin pasar por JWTAuth → no hay claims en el contexto
	mw := middleware.RequireRole("ADMIN")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "no_autenticado")
}

func TestRequireRole_RolDesconocido(t *testing.T) {
	w := aplicarConRol("ROL_INEXISTENTE", "ADMIN", "PROVEEDOR")
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── Simulación de ataque: escalada de privilegios ───────────────────────────

func TestAtaque_EscaladaPrivilegios_TokenManipulado(t *testing.T) {
	// Un VENDEDOR intenta forjar un token con rol ADMIN cambiando el payload
	// (simulado con secreto diferente → firma inválida)
	tokenFalso := generarToken("secreto-atacante", "ADMIN", 15*time.Minute)

	mw := middleware.JWTAuth(jwtTestSecret)(middleware.RequireRole("ADMIN")(okHandlerJWT))
	req := httptest.NewRequest(http.MethodPost, "/admin/ruta", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenFalso))
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Token con firma inválida debe ser rechazado")
}

func TestAtaque_TokenSinAutorizacion(t *testing.T) {
	// Intento de acceso sin token a ruta protegida
	mw := middleware.JWTAuth(jwtTestSecret)(middleware.RequireRole("ADMIN")(okHandlerJWT))
	req := httptest.NewRequest(http.MethodDelete, "/admin/usuario/1", nil)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAtaque_SoloProveedorAccedeControlAcceso(t *testing.T) {
	// ADMIN intenta acceder a rutas de control_acceso que son solo PROVEEDOR
	w := aplicarConRol("ADMIN", "PROVEEDOR") // solo PROVEEDOR permitido
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── GetJWTClaims con contexto vacío ─────────────────────────────────────────

func TestGetJWTClaims_ContextoVacio(t *testing.T) {
	claims := middleware.GetJWTClaims(context.Background())
	assert.Nil(t, claims)
}
