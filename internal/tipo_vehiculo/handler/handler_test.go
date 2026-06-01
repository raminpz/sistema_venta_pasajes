package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/tipo_vehiculo/input"
	"sistema_venta_pasajes/pkg"
	"testing"

	"github.com/gorilla/mux"
)

type fakeTipoVehiculoService struct {
	createFn  func(input.CreateTipoVehiculoInput) (*input.TipoVehiculoOutput, error)
	updateFn  func(int64, input.UpdateTipoVehiculoInput) (*input.TipoVehiculoOutput, error)
	deleteFn  func(int64) error
	getByIDFn func(int64) (*input.TipoVehiculoOutput, error)
	listFn    func() ([]input.TipoVehiculoOutput, error)
}

func (f *fakeTipoVehiculoService) Create(in input.CreateTipoVehiculoInput) (*input.TipoVehiculoOutput, error) {
	return f.createFn(in)
}
func (f *fakeTipoVehiculoService) Update(id int64, in input.UpdateTipoVehiculoInput) (*input.TipoVehiculoOutput, error) {
	return f.updateFn(id, in)
}
func (f *fakeTipoVehiculoService) Delete(id int64) error { return f.deleteFn(id) }
func (f *fakeTipoVehiculoService) GetByID(id int64) (*input.TipoVehiculoOutput, error) {
	return f.getByIDFn(id)
}
func (f *fakeTipoVehiculoService) List() ([]input.TipoVehiculoOutput, error) { return f.listFn() }

func TestHandlerCreateOK(t *testing.T) {
	h := &TipoVehiculoHandler{service: &fakeTipoVehiculoService{createFn: func(in input.CreateTipoVehiculoInput) (*input.TipoVehiculoOutput, error) {
		return &input.TipoVehiculoOutput{IDTipoVehiculo: 1, Nombre: in.Nombre, Descripcion: in.Descripcion}, nil
	}}}
	body := `{"nombre":"AUTO","descripcion":"Auto ejecutivo"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tipo-vehiculo", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatalf("esperado 201, obtuvo %d", rw.Code)
	}
}

func TestHandlerCreateErrors(t *testing.T) {
	h := &TipoVehiculoHandler{service: &fakeTipoVehiculoService{createFn: func(input.CreateTipoVehiculoInput) (*input.TipoVehiculoOutput, error) {
		return nil, errors.New("duplicado")
	}}}

	reqInvalidJSON := httptest.NewRequest(http.MethodPost, "/api/v1/tipo-vehiculo", bytes.NewBufferString("{"))
	rwInvalidJSON := httptest.NewRecorder()
	h.Create(rwInvalidJSON, reqInvalidJSON)
	if rwInvalidJSON.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtuvo %d", rwInvalidJSON.Code)
	}

	body := `{"nombre":"AUTO","descripcion":"x"}`
	reqSvcErr := httptest.NewRequest(http.MethodPost, "/api/v1/tipo-vehiculo", bytes.NewBufferString(body))
	rwSvcErr := httptest.NewRecorder()
	h.Create(rwSvcErr, reqSvcErr)
	if rwSvcErr.Code != http.StatusInternalServerError && rwSvcErr.Code != http.StatusBadRequest && rwSvcErr.Code != http.StatusConflict {
		t.Fatalf("esperado error controlado, obtuvo %d", rwSvcErr.Code)
	}
}

func TestHandlerCreateExactDBErrorMessage(t *testing.T) {
	h := &TipoVehiculoHandler{service: &fakeTipoVehiculoService{createFn: func(input.CreateTipoVehiculoInput) (*input.TipoVehiculoOutput, error) {
		return nil, pkg.BadRequest("tipo_vehiculo_create_error", "Error 1265 (01000): Data truncated for column 'NOMBRE' at row 1").WithDetails("Data truncated for column 'NOMBRE' at row 1")
	}}}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tipo-vehiculo", bytes.NewBufferString(`{"nombre":"Bus","descripcion":"Bus de alto tonelaje"}`))
	rw := httptest.NewRecorder()
	h.Create(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtuvo %d", rw.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON invalida: %v", err)
	}
	if resp["message"] != "Error 1265 (01000): Data truncated for column 'NOMBRE' at row 1" {
		t.Fatalf("mensaje inesperado: %#v", resp["message"])
	}
}

func TestHandlerUpdateOKAndErrors(t *testing.T) {
	nombre := "Combi"
	h := &TipoVehiculoHandler{service: &fakeTipoVehiculoService{updateFn: func(id int64, in input.UpdateTipoVehiculoInput) (*input.TipoVehiculoOutput, error) {
		if id == 99 {
			return nil, errors.New("no existe")
		}
		return &input.TipoVehiculoOutput{IDTipoVehiculo: id, Nombre: *in.Nombre, Descripcion: "desc"}, nil
	}}}

	reqInvalidID := httptest.NewRequest(http.MethodPut, "/api/v1/tipo-vehiculo/abc", bytes.NewBufferString(`{"nombre":"A"}`))
	reqInvalidID = mux.SetURLVars(reqInvalidID, map[string]string{"id": "abc"})
	rwInvalidID := httptest.NewRecorder()
	h.Update(rwInvalidID, reqInvalidID)
	if rwInvalidID.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtuvo %d", rwInvalidID.Code)
	}

	reqInvalidJSON := httptest.NewRequest(http.MethodPut, "/api/v1/tipo-vehiculo/1", bytes.NewBufferString("{"))
	reqInvalidJSON = mux.SetURLVars(reqInvalidJSON, map[string]string{"id": "1"})
	rwInvalidJSON := httptest.NewRecorder()
	h.Update(rwInvalidJSON, reqInvalidJSON)
	if rwInvalidJSON.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtuvo %d", rwInvalidJSON.Code)
	}

	payload, _ := json.Marshal(input.UpdateTipoVehiculoInput{Nombre: &nombre})
	reqOK := httptest.NewRequest(http.MethodPut, "/api/v1/tipo-vehiculo/1", bytes.NewReader(payload))
	reqOK = mux.SetURLVars(reqOK, map[string]string{"id": "1"})
	rwOK := httptest.NewRecorder()
	h.Update(rwOK, reqOK)
	if rwOK.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtuvo %d", rwOK.Code)
	}

	reqErr := httptest.NewRequest(http.MethodPut, "/api/v1/tipo-vehiculo/99", bytes.NewReader(payload))
	reqErr = mux.SetURLVars(reqErr, map[string]string{"id": "99"})
	rwErr := httptest.NewRecorder()
	h.Update(rwErr, reqErr)
	if rwErr.Code == http.StatusOK {
		t.Fatal("se esperaba error")
	}
}

func TestHandlerDeleteGetByIDAndList(t *testing.T) {
	h := &TipoVehiculoHandler{service: &fakeTipoVehiculoService{
		deleteFn: func(id int64) error {
			if id == 2 {
				return errors.New("fk")
			}
			return nil
		},
		getByIDFn: func(id int64) (*input.TipoVehiculoOutput, error) {
			if id == 9 {
				return nil, errors.New("not found")
			}
			return &input.TipoVehiculoOutput{IDTipoVehiculo: id, Nombre: "Auto", Descripcion: "desc"}, nil
		},
		listFn: func() ([]input.TipoVehiculoOutput, error) {
			return []input.TipoVehiculoOutput{{IDTipoVehiculo: 1, Nombre: "Auto", Descripcion: "desc"}}, nil
		},
	}}

	t.Run("delete invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tipo-vehiculo/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.Delete(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("delete error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tipo-vehiculo/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rw := httptest.NewRecorder()
		h.Delete(rw, req)
		if rw.Code == http.StatusOK {
			t.Fatal("se esperaba error")
		}
	})

	t.Run("delete ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tipo-vehiculo/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Delete(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("esperado 200, obtuvo %d", rw.Code)
		}
	})

	t.Run("get invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tipo-vehiculo/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("get service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tipo-vehiculo/9", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "9"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code == http.StatusOK {
			t.Fatal("se esperaba error")
		}
	})

	t.Run("get ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tipo-vehiculo/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("esperado 200, obtuvo %d", rw.Code)
		}
	})

	t.Run("list ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tipo-vehiculo", nil)
		rw := httptest.NewRecorder()
		h.List(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("esperado 200, obtuvo %d", rw.Code)
		}
	})
}
