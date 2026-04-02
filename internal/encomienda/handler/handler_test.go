package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/encomienda/input"
	"testing"

	"github.com/gorilla/mux"
)

type fakeEncomiendaService struct {
	createFn  func(input.CreateEncomiendaInput) (*input.EncomiendaOutput, error)
	updateFn  func(int64, input.UpdateEncomiendaInput) (*input.EncomiendaOutput, error)
	deleteFn  func(int64) error
	getByIDFn func(int64) (*input.EncomiendaOutput, error)
	listFn    func(int, int) ([]input.EncomiendaOutput, int, error)
}

func (f *fakeEncomiendaService) Create(in input.CreateEncomiendaInput) (*input.EncomiendaOutput, error) {
	return f.createFn(in)
}

func (f *fakeEncomiendaService) Update(id int64, in input.UpdateEncomiendaInput) (*input.EncomiendaOutput, error) {
	return f.updateFn(id, in)
}

func (f *fakeEncomiendaService) Delete(id int64) error { return f.deleteFn(id) }

func (f *fakeEncomiendaService) GetByID(id int64) (*input.EncomiendaOutput, error) {
	return f.getByIDFn(id)
}

func (f *fakeEncomiendaService) List(page, size int) ([]input.EncomiendaOutput, int, error) {
	return f.listFn(page, size)
}

func TestHandlerCreateOK(t *testing.T) {
	h := &EncomiendaHandler{service: &fakeEncomiendaService{createFn: func(in input.CreateEncomiendaInput) (*input.EncomiendaOutput, error) {
		return &input.EncomiendaOutput{IDEncomienda: 1, IDVenta: in.IDVenta}, nil
	}}}

	body := `{"id_venta":1,"id_programacion":1,"costo":30,"remitente_nombre":"Juan","destinatario_nombre":"Maria"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/encomienda", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)

	if rw.Code != http.StatusCreated {
		t.Fatalf("esperado 201, obtuvo %d", rw.Code)
	}
}

func TestHandlerCreateInvalidBody(t *testing.T) {
	h := &EncomiendaHandler{service: &fakeEncomiendaService{}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/encomienda", bytes.NewBufferString("{invalid"))
	rw := httptest.NewRecorder()
	h.Create(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtuvo %d", rw.Code)
	}
}

func TestHandlerListOK(t *testing.T) {
	h := &EncomiendaHandler{service: &fakeEncomiendaService{listFn: func(page, size int) ([]input.EncomiendaOutput, int, error) {
		if page != 2 || size != 5 {
			t.Fatalf("paginacion esperada page=2,size=5; obtuvo page=%d,size=%d", page, size)
		}
		return []input.EncomiendaOutput{{IDEncomienda: 1}}, 7, nil
	}}}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/encomienda?page=2&size=5", nil)
	rw := httptest.NewRecorder()
	h.List(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtuvo %d", rw.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("json invalido: %v", err)
	}
	if _, ok := resp["meta"]; !ok {
		t.Fatal("meta no presente")
	}
}

func TestHandlerGetByIDInvalidID(t *testing.T) {
	h := &EncomiendaHandler{service: &fakeEncomiendaService{}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/encomienda/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rw := httptest.NewRecorder()
	h.GetByID(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtuvo %d", rw.Code)
	}
}

func TestHandlerDeleteServiceError(t *testing.T) {
	h := &EncomiendaHandler{service: &fakeEncomiendaService{deleteFn: func(int64) error {
		return errors.New("error")
	}}}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/encomienda/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rw := httptest.NewRecorder()
	h.Delete(rw, req)

	if rw.Code == http.StatusOK {
		t.Fatal("se esperaba error")
	}
}
