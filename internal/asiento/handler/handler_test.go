package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/asiento/domain"
	"sistema_venta_pasajes/internal/asiento/input"
	"testing"

	"github.com/gorilla/mux"
)

type mockService struct {
	CreateFn         func(input.CreateAsientoInput) (*domain.Asiento, error)
	GetByIDFn        func(int64) (*domain.Asiento, error)
	ListByVehiculoFn func(int64) ([]*domain.Asiento, error)
	UpdateFn         func(int64, input.UpdateAsientoInput) error
	DeleteFn         func(int64) error
	CambiarEstadoFn  func(int64, string) error
}

func (m *mockService) Create(in input.CreateAsientoInput) (*domain.Asiento, error) {
	return m.CreateFn(in)
}
func (m *mockService) GetByID(id int64) (*domain.Asiento, error) { return m.GetByIDFn(id) }
func (m *mockService) ListByVehiculo(id int64) ([]*domain.Asiento, error) {
	return m.ListByVehiculoFn(id)
}
func (m *mockService) Update(id int64, in input.UpdateAsientoInput) error { return m.UpdateFn(id, in) }
func (m *mockService) Delete(id int64) error                              { return m.DeleteFn(id) }
func (m *mockService) CambiarEstado(id int64, estado string) error {
	return m.CambiarEstadoFn(id, estado)
}

func TestAsientoHandler_Create(t *testing.T) {
	service := &mockService{
		CreateFn: func(in input.CreateAsientoInput) (*domain.Asiento, error) {
			return &domain.Asiento{IDAsiento: 1, IDVehiculo: in.IDVehiculo, NumeroAsiento: in.NumeroAsiento}, nil
		},
	}
	h := New(service)
	body, _ := json.Marshal(input.CreateAsientoInput{IDVehiculo: 2, NumeroAsiento: "A1"})
	req := httptest.NewRequest(http.MethodPost, "/asiento", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rw.Code)
	}
}

func TestAsientoHandler_GetByID(t *testing.T) {
	service := &mockService{
		GetByIDFn: func(id int64) (*domain.Asiento, error) {
			if id == 1 {
				return &domain.Asiento{IDAsiento: 1, IDVehiculo: 2, NumeroAsiento: "A1"}, nil
			}
			return nil, nil
		},
	}
	h := New(service)
	req := httptest.NewRequest(http.MethodGet, "/asiento/1", nil)
	rw := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "1"})
	h.GetByID(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
}

func TestAsientoHandler_ListByVehiculo(t *testing.T) {
	service := &mockService{
		ListByVehiculoFn: func(id int64) ([]*domain.Asiento, error) {
			if id == 2 {
				return []*domain.Asiento{{IDAsiento: 1, IDVehiculo: 2, NumeroAsiento: "A1"}}, nil
			}
			return nil, nil
		},
	}
	h := New(service)
	req := httptest.NewRequest(http.MethodGet, "/vehiculo/2/asientos", nil)
	rw := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id_vehiculo": "2"})
	h.ListByVehiculo(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
}

func TestAsientoHandler_Update(t *testing.T) {
	updated := false
	service := &mockService{
		UpdateFn: func(id int64, in input.UpdateAsientoInput) error {
			if id == 1 && in.NumeroAsiento == "B2" && in.Estado == "OCUPADO" {
				updated = true
			}
			return nil
		},
	}
	h := New(service)
	body, _ := json.Marshal(input.UpdateAsientoInput{NumeroAsiento: "B2", Estado: "OCUPADO"})
	req := httptest.NewRequest(http.MethodPut, "/asiento/1", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "1"})
	h.Update(rw, req)
	if rw.Code != http.StatusOK || !updated {
		t.Errorf("expected status 200 and updated, got %d, updated=%v", rw.Code, updated)
	}
}

func TestAsientoHandler_Delete(t *testing.T) {
	called := false
	service := &mockService{
		DeleteFn: func(id int64) error {
			called = true
			return nil
		},
	}
	h := New(service)
	req := httptest.NewRequest(http.MethodDelete, "/asiento/1", nil)
	rw := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "1"})
	h.Delete(rw, req)
	if rw.Code != http.StatusOK || !called {
		t.Errorf("expected status 200 and called, got %d, called=%v", rw.Code, called)
	}
}

func TestAsientoHandler_CambiarEstado(t *testing.T) {
	called := false
	service := &mockService{
		CambiarEstadoFn: func(id int64, estado string) error {
			if id == 1 && estado == "OCUPADO" {
				called = true
			}
			return nil
		},
	}
	h := New(service)
	body, _ := json.Marshal(input.CambiarEstadoAsientoInput{Estado: "OCUPADO"})
	req := httptest.NewRequest(http.MethodPatch, "/asiento/1/estado", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "1"})
	h.CambiarEstado(rw, req)
	if rw.Code != http.StatusOK || !called {
		t.Errorf("expected status 200 and called, got %d, called=%v", rw.Code, called)
	}
}

// muxSetVars simula la inyección de variables de ruta en mux
func muxSetVars(r *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(r, vars)
}

func TestAsientoHandler_ErrorBranches(t *testing.T) {
	h := New(&mockService{
		CreateFn:         func(input.CreateAsientoInput) (*domain.Asiento, error) { return nil, errors.New("boom") },
		GetByIDFn:        func(int64) (*domain.Asiento, error) { return nil, errors.New("boom") },
		ListByVehiculoFn: func(int64) ([]*domain.Asiento, error) { return nil, errors.New("boom") },
		UpdateFn:         func(int64, input.UpdateAsientoInput) error { return errors.New("boom") },
		DeleteFn:         func(int64) error { return errors.New("boom") },
		CambiarEstadoFn:  func(int64, string) error { return errors.New("boom") },
	})

	t.Run("create invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/asiento", bytes.NewBufferString("{"))
		rw := httptest.NewRecorder()
		h.Create(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("create service error", func(t *testing.T) {
		body := []byte(`{"id_vehiculo":1,"numero_asiento":"A1","estado":"DISPONIBLE"}`)
		req := httptest.NewRequest(http.MethodPost, "/asiento", bytes.NewReader(body))
		rw := httptest.NewRecorder()
		h.Create(rw, req)
		if rw.Code != http.StatusInternalServerError {
			t.Fatalf("status esperado 500, obtenido %d", rw.Code)
		}
	})

	t.Run("get invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/asiento/abc", nil)
		req = muxSetVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("list vehiculo invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/vehiculo/abc/asientos", nil)
		req = muxSetVars(req, map[string]string{"id_vehiculo": "abc"})
		rw := httptest.NewRecorder()
		h.ListByVehiculo(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("update invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/asiento/abc", bytes.NewBufferString(`{"estado":"OCUPADO"}`))
		req = muxSetVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("update invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/asiento/1", bytes.NewBufferString("{"))
		req = muxSetVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("delete invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/asiento/abc", nil)
		req = muxSetVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.Delete(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("cambiar estado invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/asiento/abc/estado", bytes.NewBufferString(`{"estado":"OCUPADO"}`))
		req = muxSetVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.CambiarEstado(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("cambiar estado invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/asiento/1/estado", bytes.NewBufferString(`{"estado":"X"}`))
		req = muxSetVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.CambiarEstado(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})
}
