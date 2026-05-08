package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/programacion/input"
	"testing"

	"github.com/gorilla/mux"
)

type fakeProgramacionService struct {
	createFn  func(input.CreateProgramacionInput) (*input.ProgramacionOutput, error)
	updateFn  func(int64, input.UpdateProgramacionInput) (*input.ProgramacionOutput, error)
	deleteFn  func(int64) error
	getByIDFn func(int64) (*input.ProgramacionOutput, error)
	listFn    func(int, int) ([]input.ProgramacionOutput, int, error)
}

func (f *fakeProgramacionService) Create(in input.CreateProgramacionInput) (*input.ProgramacionOutput, error) {
	return f.createFn(in)
}

func (f *fakeProgramacionService) Update(id int64, in input.UpdateProgramacionInput) (*input.ProgramacionOutput, error) {
	return f.updateFn(id, in)
}

func (f *fakeProgramacionService) Delete(id int64) error { return f.deleteFn(id) }

func (f *fakeProgramacionService) GetByID(id int64) (*input.ProgramacionOutput, error) {
	return f.getByIDFn(id)
}

func (f *fakeProgramacionService) List(page, size int) ([]input.ProgramacionOutput, int, error) {
	return f.listFn(page, size)
}

func TestHandlerCreateOK(t *testing.T) {
	h := &ProgramacionHandler{service: &fakeProgramacionService{createFn: func(in input.CreateProgramacionInput) (*input.ProgramacionOutput, error) {
		return &input.ProgramacionOutput{IDProgramacion: 1, IDRuta: in.IDRuta}, nil
	}}}

	body := `{"id_ruta":1,"id_vehiculo":2,"id_conductor":3,"fecha_salida":"2026-04-05 02:00:00"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/programacion", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)

	if rw.Code != http.StatusCreated {
		t.Fatalf("esperado 201, obtuvo %d", rw.Code)
	}
}

func TestHandlerCreateInvalidBody(t *testing.T) {
	h := &ProgramacionHandler{service: &fakeProgramacionService{}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/programacion", bytes.NewBufferString("{invalid"))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtuvo %d", rw.Code)
	}
}

func TestHandlerListOK(t *testing.T) {
	h := &ProgramacionHandler{service: &fakeProgramacionService{listFn: func(page, size int) ([]input.ProgramacionOutput, int, error) {
		if page != 2 || size != 5 {
			t.Fatalf("paginacion esperada page=2,size=5; obtuvo page=%d,size=%d", page, size)
		}
		return []input.ProgramacionOutput{{IDProgramacion: 1}}, 7, nil
	}}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/programacion?page=2&size=5", nil)
	rw := httptest.NewRecorder()
	h.List(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtuvo %d", rw.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("json invalido: %v", err)
	}
	metaRaw, ok := resp["meta"]
	if !ok {
		t.Fatal("meta no presente")
	}
	meta, ok := metaRaw.(map[string]interface{})
	if !ok {
		t.Fatal("meta invalida")
	}
	if meta["page"].(float64) != 2 {
		t.Fatalf("page esperado 2, obtuvo %v", meta["page"])
	}
	if meta["size"].(float64) != 5 {
		t.Fatalf("size esperado 5, obtuvo %v", meta["size"])
	}
	if meta["total"].(float64) != 7 {
		t.Fatalf("total esperado 7, obtuvo %v", meta["total"])
	}
	if meta["total_pages"].(float64) != 2 {
		t.Fatalf("total_pages esperado 2, obtuvo %v", meta["total_pages"])
	}
}

func TestHandlerListEmptyDataReturnsArray(t *testing.T) {
	h := &ProgramacionHandler{service: &fakeProgramacionService{listFn: func(page, size int) ([]input.ProgramacionOutput, int, error) {
		return nil, 0, nil
	}}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/programacion?page=1&size=10", nil)
	rw := httptest.NewRecorder()
	h.List(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtuvo %d", rw.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("json invalido: %v", err)
	}
	data, ok := resp["data"].([]interface{})
	if !ok {
		t.Fatal("data debe ser array")
	}
	if len(data) != 0 {
		t.Fatalf("data esperado vacio, obtuvo %d", len(data))
	}
}

func TestHandlerGetByIDInvalidID(t *testing.T) {
	h := &ProgramacionHandler{service: &fakeProgramacionService{}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/programacion/abc", nil)
	rw := httptest.NewRecorder()
	h.GetByID(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtuvo %d", rw.Code)
	}
}

func TestHandlerDeleteServiceError(t *testing.T) {
	h := &ProgramacionHandler{service: &fakeProgramacionService{deleteFn: func(int64) error {
		return errors.New("error")
	}}}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/programacion/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rw := httptest.NewRecorder()
	h.Delete(rw, req)
	if rw.Code == http.StatusOK {
		t.Fatal("se esperaba error")
	}
}

func TestHandlerUpdateDeleteGetAndListBranches(t *testing.T) {
	h := &ProgramacionHandler{service: &fakeProgramacionService{
		updateFn: func(id int64, in input.UpdateProgramacionInput) (*input.ProgramacionOutput, error) {
			return &input.ProgramacionOutput{IDProgramacion: id}, nil
		},
		deleteFn: func(id int64) error { return nil },
		getByIDFn: func(id int64) (*input.ProgramacionOutput, error) {
			if id == 404 {
				return nil, errors.New("not found")
			}
			return &input.ProgramacionOutput{IDProgramacion: id}, nil
		},
		listFn: func(page, size int) ([]input.ProgramacionOutput, int, error) {
			if page == 9 {
				return nil, 0, errors.New("db")
			}
			return nil, 0, nil
		},
	}}

	t.Run("update id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/programacion/abc", bytes.NewBufferString(`{"estado":"PROGRAMADO"}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("update body invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/programacion/1", bytes.NewBufferString("{"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("update ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/programacion/1", bytes.NewBufferString(`{"estado":"PROGRAMADO"}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("esperado 200, obtuvo %d", rw.Code)
		}
	})

	t.Run("delete id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/programacion/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.Delete(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("delete ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/programacion/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Delete(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("esperado 200, obtuvo %d", rw.Code)
		}
	})

	t.Run("get service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/programacion/404", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "404"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusInternalServerError {
			t.Fatalf("esperado 500, obtuvo %d", rw.Code)
		}
	})

	t.Run("get ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/programacion/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("esperado 200, obtuvo %d", rw.Code)
		}
	})

	t.Run("list invalid page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/programacion?page=abc", nil)
		rw := httptest.NewRecorder()
		h.List(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperado 400, obtuvo %d", rw.Code)
		}
	})

	t.Run("list service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/programacion?page=9&size=10", nil)
		rw := httptest.NewRecorder()
		h.List(rw, req)
		if rw.Code != http.StatusInternalServerError {
			t.Fatalf("esperado 500, obtuvo %d", rw.Code)
		}
	})
}
