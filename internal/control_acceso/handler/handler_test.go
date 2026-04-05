package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
// Tests de acceso al módulo — solo rol PROVEEDOR (X-Provider-Key)
// =============================================================================

func newTestRouterConAuth(h *handler.Handler, apiKey string) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/control-acceso/status", h.GetStatus).Methods(http.MethodGet)

	admin := r.PathPrefix("/api/v1/control-acceso").Subrouter()
	admin.Use(handler.ProviderAuthMiddleware(apiKey))
	admin.HandleFunc("", h.GetLatest).Methods(http.MethodGet)
	admin.HandleFunc("", h.Create).Methods(http.MethodPost)
	admin.HandleFunc("/{id:[0-9]+}/activar", h.Activar).Methods(http.MethodPut)
	admin.HandleFunc("/{id:[0-9]+}/bloquear", h.Bloquear).Methods(http.MethodPut)
	admin.HandleFunc("/{id:[0-9]+}/renovar", h.Renovar).Methods(http.MethodPut)
	return r
}

func TestProviderAuth_ClaveValida_AccedeAdminOK(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h, "clave-proveedor-123")

	svc.On("GetLatest").Return(&input.ControlAccesoOutput{
		IDAcceso:       1,
		EstadoDB:       "OPERATIVO",
		EstadoEfectivo: "OPERATIVO",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/control-acceso", nil)
	req.Header.Set("X-Provider-Key", "clave-proveedor-123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestProviderAuth_SinClave_Rechaza401(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h, "clave-proveedor-123")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/control-acceso", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetLatest")
}

func TestProviderAuth_ClaveIncorrecta_Rechaza401(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h, "clave-proveedor-123")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/control-acceso", nil)
	req.Header.Set("X-Provider-Key", "clave-equivocada")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetLatest")
}

func TestProviderAuth_ClaveNoConfigurada_503(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h, "")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/control-acceso", nil)
	req.Header.Set("X-Provider-Key", "cualquier-clave")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	svc.AssertNotCalled(t, "GetLatest")
}

func TestProviderAuth_StatusEsPublico_SinKey(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h, "clave-proveedor-123")

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

func TestProviderAuth_OtroUsuario_NoPuedeBloquear(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h, "clave-proveedor-123")

	req := httptest.NewRequest(http.MethodPut, "/api/v1/control-acceso/1/bloquear", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Bloquear")
}

func TestProviderAuth_OtroUsuario_NoPuedeActivar(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h, "clave-proveedor-123")

	req := httptest.NewRequest(http.MethodPut, "/api/v1/control-acceso/1/activar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Activar")
}

func TestProviderAuth_OtroUsuario_NoPuedeRenovar(t *testing.T) {
	svc := new(mockService)
	h := handler.New(svc)
	r := newTestRouterConAuth(h, "clave-proveedor-123")

	req := httptest.NewRequest(http.MethodPut, "/api/v1/control-acceso/1/renovar", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "Renovar")
}
