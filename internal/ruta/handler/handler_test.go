package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sistema_venta_pasajes/internal/ruta/domain"
	"sistema_venta_pasajes/internal/ruta/input"
	"sistema_venta_pasajes/internal/ruta/service"
)

func setupHandler() (*Handler, *service.ServiceMock) {
	mockService := &service.ServiceMock{}
	h := New(mockService)
	return h, mockService
}

func TestHandler_List(t *testing.T) {
	h, mockService := setupHandler()
	rutas := []domain.Ruta{{IDRuta: 1, IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5}}
	mockService.On("List", mock.Anything).Return(rutas, nil)
	req := httptest.NewRequest(http.MethodGet, "/rutas", nil)
	w := httptest.NewRecorder()

	h.List(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByID(t *testing.T) {
	h, mockService := setupHandler()
	ruta := &domain.Ruta{IDRuta: 1, IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5}
	mockService.On("GetByID", mock.Anything, 1).Return(ruta, nil)

	req := httptest.NewRequest(http.MethodGet, "/rutas/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	h.GetByID(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Create(t *testing.T) {
	h, mockService := setupHandler()
	in := input.CreateRutaInput{IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5}
	ruta := &domain.Ruta{IDRuta: 1, IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5}
	mockService.On("Create", mock.Anything, in).Return(ruta, nil)

	body, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPost, "/rutas", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Create(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandler_Update(t *testing.T) {
	h, mockService := setupHandler()
	in := input.UpdateRutaInput{DuracionHoras: ptrFloat64(7.0)}
	ruta := &domain.Ruta{IDRuta: 1, IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 7.0}
	mockService.On("Update", mock.Anything, 1, in).Return(ruta, nil)

	body, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPut, "/rutas/1", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	h.Update(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Delete(t *testing.T) {
	h, mockService := setupHandler()
	mockService.On("Delete", mock.Anything, 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/rutas/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	h.Delete(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func ptrFloat64(f float64) *float64 {
	return &f
}

func TestHandler_ErrorBranches(t *testing.T) {
	h, mockService := setupHandler()

	t.Run("list service error", func(t *testing.T) {
		mockService.On("List", mock.Anything).Return([]domain.Ruta(nil), errors.New("db")).Once()
		req := httptest.NewRequest(http.MethodGet, "/rutas", nil)
		w := httptest.NewRecorder()
		h.List(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("get invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/rutas/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.GetByID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("get service error", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, 99).Return((*domain.Ruta)(nil), errors.New("db")).Once()
		req := httptest.NewRequest(http.MethodGet, "/rutas/99", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "99"})
		w := httptest.NewRecorder()
		h.GetByID(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("create invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/rutas", bytes.NewBufferString("{"))
		w := httptest.NewRecorder()
		h.Create(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("create service error", func(t *testing.T) {
		in := input.CreateRutaInput{IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5}
		mockService.On("Create", mock.Anything, in).Return((*domain.Ruta)(nil), errors.New("db")).Once()
		body, _ := json.Marshal(in)
		req := httptest.NewRequest(http.MethodPost, "/rutas", bytes.NewReader(body))
		w := httptest.NewRecorder()
		h.Create(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("update invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/rutas/abc", bytes.NewBufferString(`{"duracion_horas":7}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("update invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/rutas/1", bytes.NewBufferString("{"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("update service error", func(t *testing.T) {
		in := input.UpdateRutaInput{DuracionHoras: ptrFloat64(7)}
		mockService.On("Update", mock.Anything, 1, in).Return((*domain.Ruta)(nil), errors.New("db")).Once()
		body, _ := json.Marshal(in)
		req := httptest.NewRequest(http.MethodPut, "/rutas/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.Update(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("delete invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/rutas/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		h.Delete(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("delete service error", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, 1).Return(errors.New("db")).Once()
		req := httptest.NewRequest(http.MethodDelete, "/rutas/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		h.Delete(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
