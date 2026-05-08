package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/pago/input"
	"testing"

	"github.com/gorilla/mux"
)

type fakePagoService struct {
	createFn  func(input.CreatePagoInput) (*input.PagoOutput, error)
	updateFn  func(int64, input.UpdatePagoInput) (*input.PagoOutput, error)
	deleteFn  func(int64) error
	getByIDFn func(int64) (*input.PagoOutput, error)
	listFn    func(int, int, *int64) ([]input.PagoOutput, int, error)
}

func (f *fakePagoService) Create(in input.CreatePagoInput) (*input.PagoOutput, error) {
	return f.createFn(in)
}
func (f *fakePagoService) Update(id int64, in input.UpdatePagoInput) (*input.PagoOutput, error) {
	return f.updateFn(id, in)
}
func (f *fakePagoService) Delete(id int64) error { return f.deleteFn(id) }
func (f *fakePagoService) GetByID(id int64) (*input.PagoOutput, error) {
	return f.getByIDFn(id)
}
func (f *fakePagoService) List(page, size int, idVenta *int64) ([]input.PagoOutput, int, error) {
	return f.listFn(page, size, idVenta)
}

func TestHandlerCreateEstadoRegistrada(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{createFn: func(in input.CreatePagoInput) (*input.PagoOutput, error) {
		return &input.PagoOutput{IDPago: 1, IDVenta: in.IDVenta, Estado: "REGISTRADA"}, nil
	}}}
	body := `{"id_venta":1,"id_metodo":1,"monto":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pago", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatalf("esperado 201, obtuvo %d", rw.Code)
	}
}

func TestHandlerCreateEstadoParcial(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{createFn: func(in input.CreatePagoInput) (*input.PagoOutput, error) {
		return &input.PagoOutput{IDPago: 2, IDVenta: in.IDVenta, Estado: "PARCIAL"}, nil
	}}}
	body := `{"id_venta":1,"id_metodo":2,"monto":50,"estado":"PARCIAL"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pago", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatalf("esperado 201, obtuvo %d", rw.Code)
	}
}

func TestHandlerCreateEstadoPagada(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{createFn: func(in input.CreatePagoInput) (*input.PagoOutput, error) {
		return &input.PagoOutput{IDPago: 3, IDVenta: in.IDVenta, Estado: "PAGADA"}, nil
	}}}
	body := `{"id_venta":1,"id_metodo":1,"monto":100,"estado":"PAGADA"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pago", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatalf("esperado 201, obtuvo %d", rw.Code)
	}
}

func TestHandlerCreateEstadoAnulada(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{createFn: func(in input.CreatePagoInput) (*input.PagoOutput, error) {
		return &input.PagoOutput{IDPago: 4, IDVenta: in.IDVenta, Estado: "ANULADA"}, nil
	}}}
	body := `{"id_venta":2,"id_metodo":1,"monto":0,"estado":"ANULADA"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pago", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatalf("esperado 201, obtuvo %d", rw.Code)
	}
}

func TestHandlerCreateServiceError(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{createFn: func(in input.CreatePagoInput) (*input.PagoOutput, error) {
		return nil, errors.New("estado invalido")
	}}}
	body := `{"id_venta":1,"id_metodo":1,"monto":0,"estado":"PENDIENTE"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pago", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code == http.StatusCreated {
		t.Fatal("se esperaba error por estado invalido")
	}
}

func TestHandlerListOK(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{listFn: func(page, size int, idVenta *int64) ([]input.PagoOutput, int, error) {
		return []input.PagoOutput{{IDPago: 1, Estado: "REGISTRADA"}, {IDPago: 2, Estado: "PARCIAL"}}, 2, nil
	}}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/pago?page=1&size=10", nil)
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

func TestHandlerGetByIDOK(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{getByIDFn: func(id int64) (*input.PagoOutput, error) {
		return &input.PagoOutput{IDPago: id, IDVenta: 1, Estado: "PAGADA"}, nil
	}}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/pago/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rw := httptest.NewRecorder()
	h.GetByID(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtuvo %d", rw.Code)
	}
}

func TestHandlerUpdateEstadoValido(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{updateFn: func(id int64, in input.UpdatePagoInput) (*input.PagoOutput, error) {
		return &input.PagoOutput{IDPago: id, Estado: "PAGADA"}, nil
	}}}
	body := `{"estado":"PAGADA"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/pago/1", bytes.NewBufferString(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rw := httptest.NewRecorder()
	h.Update(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtuvo %d", rw.Code)
	}
}

func TestHandlerDeleteOK(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{deleteFn: func(int64) error { return nil }}}
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/pago/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rw := httptest.NewRecorder()
	h.Delete(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtuvo %d", rw.Code)
	}
}

func TestHandlerDeleteServiceError(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{deleteFn: func(int64) error { return errors.New("error") }}}
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/pago/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rw := httptest.NewRecorder()
	h.Delete(rw, req)
	if rw.Code == http.StatusOK {
		t.Fatal("se esperaba error")
	}
}

func TestPagoHandler_ErrorBranches(t *testing.T) {
	h := &PagoHandler{service: &fakePagoService{
		updateFn:  func(int64, input.UpdatePagoInput) (*input.PagoOutput, error) { return nil, errors.New("boom") },
		getByIDFn: func(int64) (*input.PagoOutput, error) { return nil, errors.New("boom") },
		listFn: func(page, size int, idVenta *int64) ([]input.PagoOutput, int, error) {
			if page == 9 {
				return nil, 0, errors.New("boom")
			}
			return nil, 0, nil
		},
	}}

	t.Run("create invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/pago", bytes.NewBufferString("{"))
		rw := httptest.NewRecorder()
		h.Create(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("update invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/pago/abc", bytes.NewBufferString(`{"estado":"PAGADA"}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("update invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/pago/1", bytes.NewBufferString("{"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("update service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/pago/1", bytes.NewBufferString(`{"estado":"PAGADA"}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusInternalServerError {
			t.Fatalf("esperado 500, obtuvo %d", rw.Code)
		}
	})

	t.Run("get invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pago/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("get service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pago/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusInternalServerError {
			t.Fatalf("esperado 500, obtuvo %d", rw.Code)
		}
	})

	t.Run("list invalid pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pago?page=abc", nil)
		rw := httptest.NewRecorder()
		h.List(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("list invalid id_venta", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pago?page=1&size=10&id_venta=abc", nil)
		rw := httptest.NewRecorder()
		h.List(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("list service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pago?page=9&size=10", nil)
		rw := httptest.NewRecorder()
		h.List(rw, req)
		if rw.Code != http.StatusInternalServerError {
			t.Fatalf("esperado 500, obtuvo %d", rw.Code)
		}
	})

	t.Run("list nil slice", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/pago?page=1&size=10", nil)
		rw := httptest.NewRecorder()
		h.List(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("esperado 200, obtuvo %d", rw.Code)
		}
	})
}
