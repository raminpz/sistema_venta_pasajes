package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mw "sistema_venta_pasajes/configs/http/middleware"
	"sistema_venta_pasajes/internal/control_acceso/handler"
	"sistema_venta_pasajes/internal/control_acceso/input"
)

// mockService simula el service
type mockService struct {
	mock.Mock
}

func (m *mockService) GetStatus() (*input.ControlAccesoStatusOutput, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*input.ControlAccesoStatusOutput), args.Error(1)
}

func (m *mockService) GetLatest() (*input.ControlAccesoOutput, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*input.ControlAccesoOutput), args.Error(1)
}

func (m *mockService) Create(in input.ActivarControlAccesoInput) (*input.ControlAccesoOutput, error) {
	args := m.Called(in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*input.ControlAccesoOutput), args.Error(1)
}

func (m *mockService) Activar(id int64) error {
	return m.Called(id).Error(0)
}

func (m *mockService) Bloquear(id int64) error {
	return m.Called(id).Error(0)
}

func (m *mockService) Renovar(id int64, in input.RenovarControlAccesoInput) (*input.ControlAccesoOutput, error) {
	args := m.Called(id, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*input.ControlAccesoOutput), args.Error(1)
}

// helper para registrar rutas en un router de test (sin middleware de auth)
func newTestRouter(h *handler.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/control-acceso/status", h.GetStatus).Methods("GET")
	r.HandleFunc("/api/v1/control-acceso", h.GetLatest).Methods("GET")
	r.HandleFunc("/api/v1/control-acceso", h.Create).Methods("POST")
	r.HandleFunc("/api/v1/control-acceso/{id:[0-9]+}/activar", h.Activar).Methods("PUT")
	r.HandleFunc("/api/v1/control-acceso/{id:[0-9]+}/bloquear", h.Bloquear).Methods("PUT")
	r.HandleFunc("/api/v1/control-acceso/{id:[0-9]+}/renovar", h.Renovar).Methods("PUT")
	return r
}

func TestHandler_GetStatus_OK(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouter(h)

	svc.On("GetStatus").Return(&input.ControlAccesoStatusOutput{
		EstadoEfectivo: "OPERATIVO",
		Mensaje:        "Sistema operativo.",
		DiasParaVencer: 180,
		EnGracia:       false,
	}, nil)

	req := httptest.NewRequest("GET", "/api/v1/control-acceso/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_GetLatest_OK(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouter(h)

	svc.On("GetLatest").Return(&input.ControlAccesoOutput{
		IDAcceso:       1,
		EstadoDB:       "OPERATIVO",
		EstadoEfectivo: "OPERATIVO",
	}, nil)

	req := httptest.NewRequest("GET", "/api/v1/control-acceso", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_Create_OK(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouter(h)

	in := input.ActivarControlAccesoInput{
		FechaActivacion: "2026-01-01",
		FechaExpiracion: "2027-01-01",
	}
	svc.On("Create", in).Return(&input.ControlAccesoOutput{
		IDAcceso:        1,
		FechaActivacion: "2026-01-01",
		FechaExpiracion: "2027-01-01",
		EstadoDB:        "OPERATIVO",
		EstadoEfectivo:  "OPERATIVO",
	}, nil)

	body, _ := json.Marshal(in)
	req := httptest.NewRequest("POST", "/api/v1/control-acceso", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_Activar_OK(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouter(h)

	svc.On("Activar", int64(1)).Return(nil)

	req := httptest.NewRequest("PUT", "/api/v1/control-acceso/1/activar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_Bloquear_OK(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouter(h)

	svc.On("Bloquear", int64(1)).Return(nil)

	req := httptest.NewRequest("PUT", "/api/v1/control-acceso/1/bloquear", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_Renovar_OK(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouter(h)

	in := input.RenovarControlAccesoInput{FechaExpiracion: "2028-01-01"}
	svc.On("Renovar", int64(1), in).Return(&input.ControlAccesoOutput{
		IDAcceso:       1,
		EstadoEfectivo: "OPERATIVO",
	}, nil)

	body, _ := json.Marshal(in)
	req := httptest.NewRequest("PUT", "/api/v1/control-acceso/1/renovar", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHandler_Create_BodyVacio(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouter(h)

	req := httptest.NewRequest("POST", "/api/v1/control-acceso", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// =============================================================================
// Tests de acceso al módulo — solo rol PROVEEDOR (JWT + RequireRole)
// =============================================================================

// generarTokenRol genera un JWT firmado con el rol indicado.
func generarTokenRol(rol string) string {
	type claims struct {
		IDUsuario int    `json:"id_usuario"`
		Email     string `json:"email"`
		Rol       string `json:"rol"`
		jwt.RegisteredClaims
	}
	c := &claims{
		IDUsuario: 1,
		Email:     "test@test.com",
		Rol:       rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, _ := t.SignedString([]byte(jwtTestSecret))
	return signed
}

const jwtTestSecret = "secreto_jwt_test_control_acceso"

func newTestRouterConAuth(h *handler.Handler) *mux.Router {
	r := mux.NewRouter()

	// Público: status
	r.HandleFunc("/api/v1/control-acceso/status", h.GetStatus).Methods(http.MethodGet)

	// Protegido: JWT + PROVEEDOR
	admin := r.PathPrefix("/api/v1/control-acceso").Subrouter()
	admin.Use(mw.JWTAuth(jwtTestSecret))
	admin.Use(mw.RequireRole("PROVEEDOR"))
	admin.HandleFunc("", h.GetLatest).Methods(http.MethodGet)
	admin.HandleFunc("", h.Create).Methods(http.MethodPost)
	admin.HandleFunc("/{id:[0-9]+}/activar", h.Activar).Methods(http.MethodPut)
	admin.HandleFunc("/{id:[0-9]+}/bloquear", h.Bloquear).Methods(http.MethodPut)
	admin.HandleFunc("/{id:[0-9]+}/renovar", h.Renovar).Methods(http.MethodPut)
	return r
}

func TestProviderAuth_RolProveedor_AccedeAdminOK(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h)

	svc.On("GetLatest").Return(&input.ControlAccesoOutput{
		IDAcceso:       1,
		EstadoDB:       "OPERATIVO",
		EstadoEfectivo: "OPERATIVO",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/control-acceso", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generarTokenRol("PROVEEDOR")))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestProviderAuth_SinToken_Rechaza401(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/control-acceso", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetLatest")
}

func TestProviderAuth_RolAdmin_Rechaza403(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/control-acceso", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generarTokenRol("ADMIN")))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	svc.AssertNotCalled(t, "GetLatest")
}

func TestProviderAuth_RolVendedor_Rechaza403(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/control-acceso", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", generarTokenRol("VENDEDOR")))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProviderAuth_StatusEsPublico_SinToken(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h)

	svc.On("GetStatus").Return(&input.ControlAccesoStatusOutput{
		EstadoEfectivo: "OPERATIVO",
		Mensaje:        "Sistema operativo.",
		DiasParaVencer: 365,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/control-acceso/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "el endpoint de status debe ser público")
	svc.AssertExpectations(t)
}

func TestProviderAuth_SinToken_NoPuedeBloquear(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/control-acceso/1/bloquear", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Bloquear")
}

func TestProviderAuth_SinToken_NoPuedeActivar(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/control-acceso/1/activar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Activar")
}

func TestProviderAuth_SinToken_NoPuedeRenovar(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/control-acceso/1/renovar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Renovar")
}
