package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/pkg"
	"testing"

	"github.com/gorilla/mux"
	"sistema_venta_pasajes/internal/conductor/domain"
	"sistema_venta_pasajes/internal/conductor/input"
)

type mockService struct {
	listErr   error
	getErr    error
	createErr error
	updateErr error
	deleteErr error
}

func (m *mockService) List(ctx context.Context) ([]domain.Conductor, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return []domain.Conductor{{IDConductor: 1, Nombres: "Juan", Apellidos: "Perez", NumeroLicencia: "ABC123456", Telefono: "987654321"}}, nil
}
func (m *mockService) GetByID(ctx context.Context, id int64) (*domain.Conductor, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return &domain.Conductor{IDConductor: id, Nombres: "Juan", Apellidos: "Perez", NumeroLicencia: "ABC123456", Telefono: "987654321"}, nil
}
func (m *mockService) Create(ctx context.Context, in input.CreateConductorInput) (*domain.Conductor, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &domain.Conductor{IDConductor: 2, Nombres: in.Nombres, Apellidos: in.Apellidos, NumeroLicencia: in.NumeroLicencia, Telefono: in.Telefono}, nil
}
func (m *mockService) Update(ctx context.Context, id int64, in input.UpdateConductorInput) (*domain.Conductor, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return &domain.Conductor{IDConductor: id, Nombres: "Mod", Apellidos: "Mod", NumeroLicencia: "MOD123456", Telefono: "987654321"}, nil
}
func (m *mockService) Delete(ctx context.Context, id int64) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	return nil
}

func TestHandler_Create(t *testing.T) {
	h := New(&mockService{})
	r := mux.NewRouter()
	r.HandleFunc("/conductor", h.Create).Methods("POST")
	body := input.CreateConductorInput{
		Nombres:           "Juan",
		Apellidos:         "Perez",
		DNI:               "12345678",
		NumeroLicencia:    "ABC123456",
		Telefono:          "987654321",
		FechaVencLicencia: "2027-01-01",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/conductor", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestHandler_List_ErrorsAndSuccess(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/conductor", New(&mockService{}).List).Methods(http.MethodGet)
	r.HandleFunc("/conductor-error", New(&mockService{listErr: errors.New("boom")}).List).Methods(http.MethodGet)

	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/conductor", nil)
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("status esperado 200, obtenido %d", w1.Code)
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/conductor-error", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusInternalServerError {
		t.Fatalf("status esperado 500, obtenido %d", w2.Code)
	}
}

func TestHandler_GetByID_Update_Delete(t *testing.T) {
	h := New(&mockService{})

	t.Run("get id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/conductor/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.GetByID(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", w.Code)
		}
	})

	t.Run("get servicio not found", func(t *testing.T) {
		h2 := New(&mockService{getErr: pkg.NotFound("conductor_not_found", "No encontrado")})
		req := httptest.NewRequest(http.MethodGet, "/conductor/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h2.GetByID(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("status esperado 404, obtenido %d", w.Code)
		}
	})

	t.Run("update json invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/conductor/1", bytes.NewBufferString("{"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.Update(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", w.Code)
		}
	})

	t.Run("delete id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/conductor/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.Delete(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", w.Code)
		}
	})

	t.Run("delete error", func(t *testing.T) {
		h2 := New(&mockService{deleteErr: pkg.NotFound("conductor_not_found", "No encontrado")})
		req := httptest.NewRequest(http.MethodDelete, "/conductor/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h2.Delete(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("status esperado 404, obtenido %d", w.Code)
		}
	})
}
