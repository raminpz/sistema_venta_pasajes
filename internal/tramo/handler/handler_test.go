package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/tramo/input"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// ── mock service ─────────────────────────────────────────────────────────────

type mockTramoService struct {
	createFn     func(input.CreateTramoInput) (*input.TramoOutput, error)
	updateFn     func(input.UpdateTramoInput) (*input.TramoOutput, error)
	deleteFn     func(int64) error
	getByIDFn    func(int64) (*input.TramoOutput, error)
	listFn       func(int, int) ([]input.TramoOutput, int, error)
	listByRutaFn func(int64) ([]input.TramoOutput, error)
}

func (m *mockTramoService) Create(in input.CreateTramoInput) (*input.TramoOutput, error) {
	return m.createFn(in)
}
func (m *mockTramoService) Update(in input.UpdateTramoInput) (*input.TramoOutput, error) {
	return m.updateFn(in)
}
func (m *mockTramoService) Delete(id int64) error { return m.deleteFn(id) }
func (m *mockTramoService) GetByID(id int64) (*input.TramoOutput, error) {
	return m.getByIDFn(id)
}
func (m *mockTramoService) List(page, size int) ([]input.TramoOutput, int, error) {
	return m.listFn(page, size)
}
func (m *mockTramoService) ListByRuta(idRuta int64) ([]input.TramoOutput, error) {
	return m.listByRutaFn(idRuta)
}

// ── tests ─────────────────────────────────────────────────────────────────────

func TestHandler_Create_OK(t *testing.T) {
	svc := &mockTramoService{
		createFn: func(in input.CreateTramoInput) (*input.TramoOutput, error) {
			return &input.TramoOutput{IDTramo: 1, IDRuta: in.IDRuta}, nil
		},
	}
	h := NewTramoHandler(svc)
	body, _ := json.Marshal(map[string]any{
		"id_ruta": 1, "id_parada_origen": 1, "id_parada_destino": 2,
	})
	req := httptest.NewRequest(http.MethodPost, "/tramo", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestHandler_GetByID_OK(t *testing.T) {
	svc := &mockTramoService{
		getByIDFn: func(id int64) (*input.TramoOutput, error) {
			return &input.TramoOutput{IDTramo: id, IDRuta: 1}, nil
		},
	}
	h := NewTramoHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/tramo/1", nil)
	rec := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tramo/{id}", h.GetByID)
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_Delete_OK(t *testing.T) {
	svc := &mockTramoService{
		deleteFn: func(id int64) error { return nil },
	}
	h := NewTramoHandler(svc)
	req := httptest.NewRequest(http.MethodDelete, "/tramo/1", nil)
	rec := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tramo/{id}", h.Delete)
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_List_OK(t *testing.T) {
	svc := &mockTramoService{
		listFn: func(page, size int) ([]input.TramoOutput, int, error) {
			return []input.TramoOutput{{IDTramo: 1}}, 1, nil
		},
	}
	h := NewTramoHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/tramos", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_ListByRuta_OK(t *testing.T) {
	svc := &mockTramoService{
		listByRutaFn: func(idRuta int64) ([]input.TramoOutput, error) {
			return []input.TramoOutput{{IDTramo: 1, IDRuta: idRuta}}, nil
		},
	}
	h := NewTramoHandler(svc)
	req := httptest.NewRequest(http.MethodGet, "/tramos/ruta/1", nil)
	rec := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tramos/ruta/{id_ruta}", h.ListByRuta)
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_ErrorBranches(t *testing.T) {
	h := NewTramoHandler(&mockTramoService{
		createFn:     func(input.CreateTramoInput) (*input.TramoOutput, error) { return nil, errors.New("boom") },
		updateFn:     func(input.UpdateTramoInput) (*input.TramoOutput, error) { return nil, errors.New("boom") },
		deleteFn:     func(int64) error { return errors.New("boom") },
		getByIDFn:    func(int64) (*input.TramoOutput, error) { return nil, errors.New("boom") },
		listFn:       func(int, int) ([]input.TramoOutput, int, error) { return nil, 0, errors.New("boom") },
		listByRutaFn: func(int64) ([]input.TramoOutput, error) { return nil, errors.New("boom") },
	})

	t.Run("create invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tramo", bytes.NewBufferString("{"))
		rec := httptest.NewRecorder()
		h.Create(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("create service error", func(t *testing.T) {
		body := []byte(`{"id_ruta":1,"id_parada_origen":1,"id_parada_destino":2}`)
		req := httptest.NewRequest(http.MethodPost, "/tramo", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		h.Create(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("update invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/tramo/abc", bytes.NewBufferString(`{"id_ruta":1}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()
		h.Update(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("update invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/tramo/1", bytes.NewBufferString("{"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()
		h.Update(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("delete invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/tramo/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()
		h.Delete(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("delete service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/tramo/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()
		h.Delete(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("get invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tramo/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()
		h.GetByID(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("list invalid pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tramos?page=x", nil)
		rec := httptest.NewRecorder()
		h.List(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("list service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tramos?page=1&size=10", nil)
		rec := httptest.NewRecorder()
		h.List(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("list by ruta invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tramos/ruta/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id_ruta": "abc"})
		rec := httptest.NewRecorder()
		h.ListByRuta(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("list by ruta service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tramos/ruta/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id_ruta": "1"})
		rec := httptest.NewRecorder()
		h.ListByRuta(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
