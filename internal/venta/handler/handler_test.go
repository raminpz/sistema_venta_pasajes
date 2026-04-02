package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sistema_venta_pasajes/internal/venta/input"
)

type fakeVentaService struct {
	CreateFn  func(input.VentaCreateInput) (*input.VentaOutput, error)
	UpdateFn  func(int64, input.VentaUpdateInput) (*input.VentaOutput, error)
	DeleteFn  func(int64) error
	GetByIDFn func(int64) (*input.VentaOutput, error)
	ListFn    func(int, int) ([]input.VentaOutput, int, error)
}

func (f *fakeVentaService) Create(in input.VentaCreateInput) (*input.VentaOutput, error) {
	return f.CreateFn(in)
}
func (f *fakeVentaService) Update(id int64, in input.VentaUpdateInput) (*input.VentaOutput, error) {
	return f.UpdateFn(id, in)
}
func (f *fakeVentaService) Delete(id int64) error { return f.DeleteFn(id) }
func (f *fakeVentaService) GetByID(id int64) (*input.VentaOutput, error) {
	return f.GetByIDFn(id)
}
func (f *fakeVentaService) List(page, size int) ([]input.VentaOutput, int, error) {
	return f.ListFn(page, size)
}

// ---------------------------------------------------------------------------
// Tests Create
// ---------------------------------------------------------------------------

func TestHandler_Create_OK(t *testing.T) {
	h := &VentaHandler{service: &fakeVentaService{
		CreateFn: func(in input.VentaCreateInput) (*input.VentaOutput, error) {
			return &input.VentaOutput{IDVenta: 1, Serie: "F001", NumeroComprobante: "F001-000001"}, nil
		},
	}}
	body := `{"id_usuario":1,"id_tipo_comprobante":2,"subtotal":100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/venta", strings.NewReader(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusCreated {
		t.Errorf("esperado 201, obtenido %d", rw.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("error decodificando respuesta: %v", err)
	}
	if resp["code"].(float64) != 201 {
		t.Errorf("esperado code 201, obtenido %v", resp["code"])
	}
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h := &VentaHandler{service: &fakeVentaService{}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/venta", strings.NewReader("{invalid"))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", rw.Code)
	}
}

func TestHandler_Create_ServiceError(t *testing.T) {
	h := &VentaHandler{service: &fakeVentaService{
		CreateFn: func(in input.VentaCreateInput) (*input.VentaOutput, error) {
			return nil, errors.New("error interno")
		},
	}}
	body := `{"id_usuario":1,"id_tipo_comprobante":2,"subtotal":100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/venta", strings.NewReader(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Errorf("esperado 400, obtenido %d", rw.Code)
	}
}

// ---------------------------------------------------------------------------
// Tests List
// ---------------------------------------------------------------------------

func TestHandler_List_OK(t *testing.T) {
	h := &VentaHandler{service: &fakeVentaService{
		ListFn: func(page, size int) ([]input.VentaOutput, int, error) {
			return []input.VentaOutput{{IDVenta: 1}, {IDVenta: 2}}, 2, nil
		},
	}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/venta", nil)
	rw := httptest.NewRecorder()
	h.List(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("esperado 200, obtenido %d", rw.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("error decodificando respuesta: %v", err)
	}
	if _, ok := resp["meta"]; !ok {
		t.Error("se esperaba meta en la respuesta")
	}
}

func TestHandler_List_Error(t *testing.T) {
	h := &VentaHandler{service: &fakeVentaService{
		ListFn: func(page, size int) ([]input.VentaOutput, int, error) {
			return nil, 0, errors.New("db error")
		},
	}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/venta", nil)
	rw := httptest.NewRecorder()
	h.List(rw, req)
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("esperado 500, obtenido %d", rw.Code)
	}
}

// ---------------------------------------------------------------------------
// Tests Delete
// ---------------------------------------------------------------------------

func TestHandler_Delete_OK(t *testing.T) {
	h := &VentaHandler{service: &fakeVentaService{
		DeleteFn: func(id int64) error { return nil },
	}}
	err := h.service.Delete(1)
	if err != nil {
		t.Errorf("no se esperaba error al eliminar: %v", err)
	}
}

func TestHandler_Delete_ServiceError(t *testing.T) {
	h := &VentaHandler{service: &fakeVentaService{
		DeleteFn: func(id int64) error { return errors.New("fail") },
	}}
	err := h.service.Delete(1)
	if err == nil {
		t.Error("debe retornar error")
	}
}
