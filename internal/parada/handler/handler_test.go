package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/parada/input"
	"testing"

	"github.com/gorilla/mux"
)

type fakeParadaService struct {
	createFn     func(input.CreateParadaInput) (*input.ParadaOutput, error)
	updateFn     func(input.UpdateParadaInput) (*input.ParadaOutput, error)
	deleteFn     func(int64) error
	getByIDFn    func(int64) (*input.ParadaOutput, error)
	listByRutaFn func(int64) ([]input.ParadaOutput, error)
}

func (f *fakeParadaService) Create(in input.CreateParadaInput) (*input.ParadaOutput, error) {
	return f.createFn(in)
}

func (f *fakeParadaService) Update(in input.UpdateParadaInput) (*input.ParadaOutput, error) {
	return f.updateFn(in)
}

func (f *fakeParadaService) Delete(id int64) error {
	return f.deleteFn(id)
}

func (f *fakeParadaService) GetByID(id int64) (*input.ParadaOutput, error) {
	return f.getByIDFn(id)
}

func (f *fakeParadaService) ListByRuta(idRuta int64) ([]input.ParadaOutput, error) {
	return f.listByRutaFn(idRuta)
}

func TestParadaHandler_CRUDAndList(t *testing.T) {
	h := NewParadaHandler(&fakeParadaService{
		createFn: func(in input.CreateParadaInput) (*input.ParadaOutput, error) {
			return &input.ParadaOutput{IDParada: 1, IDRuta: in.IDRuta, NombreParada: in.NombreParada, Orden: in.Orden}, nil
		},
		updateFn: func(in input.UpdateParadaInput) (*input.ParadaOutput, error) {
			name := "Nueva"
			if in.NombreParada != nil {
				name = *in.NombreParada
			}
			return &input.ParadaOutput{IDParada: in.IDParada, NombreParada: name}, nil
		},
		deleteFn: func(int64) error { return nil },
		getByIDFn: func(id int64) (*input.ParadaOutput, error) {
			if id == 404 {
				return nil, errors.New("not found")
			}
			return &input.ParadaOutput{IDParada: id, IDRuta: 1, NombreParada: "Huanta", Orden: 1}, nil
		},
		listByRutaFn: func(idRuta int64) ([]input.ParadaOutput, error) {
			return []input.ParadaOutput{{IDParada: 1, IDRuta: idRuta, NombreParada: "Huanta", Orden: 1}}, nil
		},
	})

	t.Run("create json invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/paradas", bytes.NewReader([]byte("{")))
		rw := httptest.NewRecorder()
		h.Create(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("create ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/paradas", bytes.NewReader([]byte(`{"id_ruta":1,"nombre_parada":"Huanta","orden":1}`)))
		rw := httptest.NewRecorder()
		h.Create(rw, req)
		if rw.Code != http.StatusCreated {
			t.Fatalf("status esperado 201, obtenido %d", rw.Code)
		}
	})

	t.Run("update id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/paradas/abc", bytes.NewReader([]byte(`{"nombre_parada":"Nueva"}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("update ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/paradas/1", bytes.NewReader([]byte(`{"nombre_parada":"Nueva"}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("status esperado 200, obtenido %d", rw.Code)
		}
	})

	t.Run("get id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/paradas/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("get error servicio", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/paradas/404", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "404"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusInternalServerError {
			t.Fatalf("status esperado 500, obtenido %d", rw.Code)
		}
	})

	t.Run("delete ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/paradas/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Delete(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("status esperado 200, obtenido %d", rw.Code)
		}
	})

	t.Run("list by ruta id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/paradas/ruta/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id_ruta": "abc"})
		rw := httptest.NewRecorder()
		h.ListByRuta(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", rw.Code)
		}
	})

	t.Run("list by ruta ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/paradas/ruta/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id_ruta": "1"})
		rw := httptest.NewRecorder()
		h.ListByRuta(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("status esperado 200, obtenido %d", rw.Code)
		}
	})
}
