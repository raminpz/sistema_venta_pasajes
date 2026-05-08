package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/pkg"
	"testing"

	"sistema_venta_pasajes/internal/terminal/domain"
	"sistema_venta_pasajes/internal/terminal/input"

	"github.com/gorilla/mux"
)

type mockService struct {
	createFn  func(input.CreateTerminalInput) (*domain.Terminal, error)
	getByIDFn func(int64) (*domain.Terminal, error)
	updateFn  func(int64, input.UpdateTerminalInput) (*domain.Terminal, error)
	deleteFn  func(int64) error
	listFn    func() ([]domain.Terminal, error)
}

func (m *mockService) Create(in input.CreateTerminalInput) (*domain.Terminal, error) {
	return m.createFn(in)
}
func (m *mockService) GetByID(id int64) (*domain.Terminal, error) {
	return m.getByIDFn(id)
}
func (m *mockService) Update(id int64, in input.UpdateTerminalInput) (*domain.Terminal, error) {
	return m.updateFn(id, in)
}
func (m *mockService) Delete(id int64) error {
	return m.deleteFn(id)
}
func (m *mockService) List() ([]domain.Terminal, error) {
	return m.listFn()
}

func TestTerminalHandler_CRUDAndList(t *testing.T) {
	ms := &mockService{
		createFn: func(in input.CreateTerminalInput) (*domain.Terminal, error) {
			return &domain.Terminal{IDTerminal: 1, NOMBRE: in.Nombre, CIUDAD: in.Ciudad, ESTADO: in.Estado}, nil
		},
		getByIDFn: func(id int64) (*domain.Terminal, error) {
			if id == 404 {
				return nil, pkg.NotFound("terminal_not_found", "No encontrado")
			}
			return &domain.Terminal{IDTerminal: id, NOMBRE: "Terminal Test", CIUDAD: "Ayacucho", ESTADO: "ACTIVO"}, nil
		},
		updateFn: func(id int64, in input.UpdateTerminalInput) (*domain.Terminal, error) {
			name := "Terminal"
			if in.Nombre != "" {
				name = in.Nombre
			}
			return &domain.Terminal{IDTerminal: id, NOMBRE: name, ESTADO: "ACTIVO"}, nil
		},
		deleteFn: func(id int64) error {
			if id == 404 {
				return errors.New("no existe")
			}
			return nil
		},
		listFn: func() ([]domain.Terminal, error) {
			return []domain.Terminal{{IDTerminal: 1, NOMBRE: "A", CIUDAD: "B", ESTADO: "ACTIVO"}}, nil
		},
	}
	h := NewTerminalHandler(ms)
	r := mux.NewRouter()
	r.HandleFunc("/terminal", h.Create).Methods("POST")
	r.HandleFunc("/terminal", h.List).Methods("GET")
	r.HandleFunc("/terminal/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/terminal/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/terminal/{id}", h.Delete).Methods("DELETE")
	t.Run("create ok", func(t *testing.T) {
		body := []byte(`{"nombre":"Terminal Test","ciudad":"Ciudad Test","departamento":"Depto","direccion":"Dir","estado":"ACTIVO"}`)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/terminal", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("status esperado 201, obtenido %d", w.Code)
		}
	})

	t.Run("create json invalido", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/terminal", bytes.NewBufferString("{"))
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", w.Code)
		}
	})

	t.Run("get id invalido", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/terminal/abc", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", w.Code)
		}
	})

	t.Run("get ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/terminal/1", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status esperado 200, obtenido %d", w.Code)
		}
	})

	t.Run("update ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/terminal/1", bytes.NewBufferString(`{"nombre":"Terminal Edit"}`))
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status esperado 200, obtenido %d", w.Code)
		}
	})

	t.Run("delete error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/terminal/404", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status esperado 500, obtenido %d", w.Code)
		}
	})

	t.Run("list ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/terminal", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status esperado 200, obtenido %d", w.Code)
		}
	})
}
