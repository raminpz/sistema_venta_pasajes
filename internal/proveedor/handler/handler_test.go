package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/pkg"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	domain "sistema_venta_pasajes/internal/proveedor/domain"
	providerinput "sistema_venta_pasajes/internal/proveedor/input"
)

type fakeService struct {
	listFn   func(ctx context.Context) ([]domain.ProveedorSistema, error)
	getFn    func(ctx context.Context, id int64) (*domain.ProveedorSistema, error)
	createFn func(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error)
	updateFn func(ctx context.Context, id int64, input providerinput.UpdateInput) (*domain.ProveedorSistema, error)
	deleteFn func(ctx context.Context, id int64) error
}

func (f fakeService) List(ctx context.Context) ([]domain.ProveedorSistema, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return []domain.ProveedorSistema{}, nil
}

func (f fakeService) GetByID(ctx context.Context, id int64) (*domain.ProveedorSistema, error) {
	if f.getFn != nil {
		return f.getFn(ctx, id)
	}
	return nil, nil
}

func (f fakeService) Create(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error) {
	if f.createFn != nil {
		return f.createFn(ctx, input)
	}
	return nil, nil
}

func (f fakeService) Update(ctx context.Context, id int64, input providerinput.UpdateInput) (*domain.ProveedorSistema, error) {
	if f.updateFn != nil {
		return f.updateFn(ctx, id, input)
	}
	return nil, nil
}

func (f fakeService) Delete(ctx context.Context, id int64) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

func TestHandlerCreateReturnsCreatedResponse(t *testing.T) {
	h := NewHandler(fakeService{
		createFn: func(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error) {
			return &domain.ProveedorSistema{
				IDProveedor: 1,
				RUC:         input.RUC,
				RazonSocial: input.RazonSocial,
			}, nil
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/api/v1/proveedor", strings.NewReader(`{"ruc":"20123456789","razon_social":"Empresa Test SAC"}`))
	response := httptest.NewRecorder()

	h.Create(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("se esperaba el estado %d, se obtuvo %d", http.StatusCreated, response.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("se esperaba una respuesta JSON válida, se obtuvo %v", err)
	}
	if payload["code"] != float64(http.StatusCreated) {
		t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusCreated, payload["code"])
	}

	if payload["message"] != "proveedor del sistema creado correctamente" {
		t.Fatalf("mensaje inesperado: %#v", payload["message"])
	}

	data, ok := payload["data"].(map[string]any)
	if !ok {
		t.Fatalf("se esperaba un objeto en data, se obtuvo %#v", payload["data"])
	}
	if data["ruc"] != "20123456789" {
		t.Fatalf("se esperaba ruc 20123456789, se obtuvo %#v", data["ruc"])
	}
}

func TestHandlerCreateRejectsInvalidJSON(t *testing.T) {
	h := NewHandler(fakeService{})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/proveedor", strings.NewReader(`{"ruc":`))
	request = request.WithContext(pkg.WithRequestID(request.Context(), "req-123"))
	response := httptest.NewRecorder()

	h.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba el estado %d, se obtuvo %d", http.StatusBadRequest, response.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("se esperaba una respuesta JSON válida, se obtuvo %v", err)
	}
	if payload["code"] != float64(http.StatusBadRequest) {
		t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusBadRequest, payload["code"])
	}
	if payload["error"] != "invalid_json" {
		t.Fatalf("se esperaba invalid_json, se obtuvo %#v", payload["error"])
	}
}

func TestHandlerGetByIDRejectsInvalidID(t *testing.T) {
	h := NewHandler(fakeService{})
	request := httptest.NewRequest(http.MethodGet, "/api/v1/proveedor/abc", nil)
	request = mux.SetURLVars(request, map[string]string{"id": "abc"})
	request = request.WithContext(pkg.WithRequestID(request.Context(), "req-456"))
	response := httptest.NewRecorder()

	h.GetByID(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba el estado %d, se obtuvo %d", http.StatusBadRequest, response.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("se esperaba una respuesta JSON válida, se obtuvo %v", err)
	}
	if payload["code"] != float64(http.StatusBadRequest) {
		t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusBadRequest, payload["code"])
	}
	if payload["error"] != "invalid_provider_id" {
		t.Fatalf("se esperaba invalid_provider_id, se obtuvo %#v", payload["error"])
	}
}

func TestHandlerDeleteReturnsSuccessResponse(t *testing.T) {
	h := NewHandler(fakeService{
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	})

	request := httptest.NewRequest(http.MethodDelete, "/api/v1/proveedor/1", nil)
	request = mux.SetURLVars(request, map[string]string{"id": "1"})
	response := httptest.NewRecorder()

	h.Delete(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("se esperaba el estado %d, se obtuvo %d", http.StatusOK, response.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("se esperaba una respuesta JSON válida, se obtuvo %v", err)
	}
	if payload["code"] != float64(http.StatusOK) {
		t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusOK, payload["code"])
	}
	if payload["message"] != "proveedor del sistema eliminado correctamente" {
		t.Fatalf("mensaje inesperado: %#v", payload["message"])
	}
}
