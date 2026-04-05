package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"sistema_venta_pasajes/configs/http/middleware"
)

// okHandler es un handler auxiliar que siempre responde 200 OK.
var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

// aplicarConEstado aplica ControlAccesoConFetcher con un estado fijo y ejecuta la petición.
func aplicarConEstado(estado, metodo string) *httptest.ResponseRecorder {
	mw := middleware.ControlAccesoConFetcher(func() string { return estado })
	req := httptest.NewRequest(metodo, "/cualquier-ruta", nil)
	w := httptest.NewRecorder()
	mw(okHandler).ServeHTTP(w, req)
	return w
}

// =============================================================================
// ESTADO: OPERATIVO — todas las peticiones pasan
// =============================================================================

func TestControlAcceso_Operativo_PermiteGET(t *testing.T) {
	w := aplicarConEstado("OPERATIVO", http.MethodGet)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestControlAcceso_Operativo_PermitePOST(t *testing.T) {
	w := aplicarConEstado("OPERATIVO", http.MethodPost)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestControlAcceso_Operativo_PermitePUT(t *testing.T) {
	w := aplicarConEstado("OPERATIVO", http.MethodPut)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestControlAcceso_Operativo_PermiteDELETE(t *testing.T) {
	w := aplicarConEstado("OPERATIVO", http.MethodDelete)
	assert.Equal(t, http.StatusOK, w.Code)
}

// =============================================================================
// ESTADO: SOLO_LECTURA — GET pasa, escrituras bloqueadas (403)
// =============================================================================

func TestControlAcceso_SoloLectura_PermiteGET(t *testing.T) {
	w := aplicarConEstado("SOLO_LECTURA", http.MethodGet)
	assert.Equal(t, http.StatusOK, w.Code, "GET debe pasar en modo solo lectura")
}

func TestControlAcceso_SoloLectura_BloqueaPOST(t *testing.T) {
	w := aplicarConEstado("SOLO_LECTURA", http.MethodPost)
	assert.Equal(t, http.StatusForbidden, w.Code, "POST debe estar bloqueado en modo solo lectura")
}

func TestControlAcceso_SoloLectura_BloqueaPUT(t *testing.T) {
	w := aplicarConEstado("SOLO_LECTURA", http.MethodPut)
	assert.Equal(t, http.StatusForbidden, w.Code, "PUT debe estar bloqueado en modo solo lectura")
}

func TestControlAcceso_SoloLectura_BloqueaDELETE(t *testing.T) {
	w := aplicarConEstado("SOLO_LECTURA", http.MethodDelete)
	assert.Equal(t, http.StatusForbidden, w.Code, "DELETE debe estar bloqueado en modo solo lectura")
}

func TestControlAcceso_SoloLectura_MensajeError(t *testing.T) {
	w := aplicarConEstado("SOLO_LECTURA", http.MethodPost)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "sistema_solo_lectura")
}

// =============================================================================
// ESTADO: BLOQUEADO — todas las peticiones bloqueadas (403)
// =============================================================================

func TestControlAcceso_Bloqueado_BloqueaGET(t *testing.T) {
	w := aplicarConEstado("BLOQUEADO", http.MethodGet)
	assert.Equal(t, http.StatusForbidden, w.Code, "GET debe estar bloqueado cuando el sistema está BLOQUEADO")
}

func TestControlAcceso_Bloqueado_BloqueaPOST(t *testing.T) {
	w := aplicarConEstado("BLOQUEADO", http.MethodPost)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestControlAcceso_Bloqueado_BloqueaPUT(t *testing.T) {
	w := aplicarConEstado("BLOQUEADO", http.MethodPut)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestControlAcceso_Bloqueado_BloqueaDELETE(t *testing.T) {
	w := aplicarConEstado("BLOQUEADO", http.MethodDelete)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestControlAcceso_Bloqueado_MensajeError(t *testing.T) {
	w := aplicarConEstado("BLOQUEADO", http.MethodGet)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "sistema_bloqueado")
}

// =============================================================================
// DB nil — el middleware pasa directamente (guard para tests)
// =============================================================================

func TestControlAcceso_DBNil_PermiteGET(t *testing.T) {
	mw := middleware.ControlAccesoSistema(nil)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	mw(okHandler).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestControlAcceso_DBNil_PermitePOST(t *testing.T) {
	mw := middleware.ControlAccesoSistema(nil)
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	mw(okHandler).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
