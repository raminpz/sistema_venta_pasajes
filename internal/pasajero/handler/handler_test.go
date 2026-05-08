package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"sistema_venta_pasajes/internal/pasajero/input"
	"sistema_venta_pasajes/pkg"
)

type fakeService struct {
	CreateFn  func(input.CreatePasajeroInput) (input.PasajeroOutput, error)
	UpdateFn  func(int64, input.UpdatePasajeroInput) (input.PasajeroOutput, error)
	DeleteFn  func(int64) error
	GetByIDFn func(int64) (input.PasajeroOutput, error)
	ListFn    func(int, int) ([]input.PasajeroOutput, pkg.PaginationMeta, error)
	SearchFn  func(string) ([]input.PasajeroOutput, error)
}

func (f *fakeService) Create(in input.CreatePasajeroInput) (input.PasajeroOutput, error) {
	return f.CreateFn(in)
}
func (f *fakeService) Update(id int64, in input.UpdatePasajeroInput) (input.PasajeroOutput, error) {
	if f.UpdateFn != nil {
		return f.UpdateFn(id, in)
	}
	return input.PasajeroOutput{IDPasajero: id}, nil
}
func (f *fakeService) Delete(id int64) error {
	if f.DeleteFn != nil {
		return f.DeleteFn(id)
	}
	return nil
}
func (f *fakeService) GetByID(id int64) (input.PasajeroOutput, error) {
	if f.GetByIDFn != nil {
		return f.GetByIDFn(id)
	}
	return input.PasajeroOutput{IDPasajero: id}, nil
}
func (f *fakeService) List(page, size int) ([]input.PasajeroOutput, pkg.PaginationMeta, error) {
	if f.ListFn != nil {
		return f.ListFn(page, size)
	}
	return nil, pkg.PaginationMeta{}, nil
}
func (f *fakeService) Search(query string) ([]input.PasajeroOutput, error) {
	if f.SearchFn != nil {
		return f.SearchFn(query)
	}
	return nil, nil
}

func TestPasajeroHandler_Create(t *testing.T) {
	svc := &fakeService{
		CreateFn: func(in input.CreatePasajeroInput) (input.PasajeroOutput, error) {
			return input.PasajeroOutput{IDPasajero: 1, Nombres: "Juan"}, nil
		},
	}
	h := &PasajeroHandler{service: svc}
	body, _ := json.Marshal(input.CreatePasajeroInput{
		TipoDocumento: "DNI",
		NroDocumento:  "12345678",
		Nombres:       "Juan",
		Apellidos:     "Perez",
		Telefono:      "987654321",
	})
	req := httptest.NewRequest(http.MethodPost, "/pasajero", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusCreated {
		t.Errorf("esperaba status 201, obtuve %d", rw.Code)
	}
}

func TestPasajeroHandler_Create_JSONInvalido(t *testing.T) {
	h := &PasajeroHandler{service: &fakeService{}}
	req := httptest.NewRequest(http.MethodPost, "/pasajero", bytes.NewReader([]byte("{")))
	rw := httptest.NewRecorder()
	h.Create(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("esperaba status 400, obtuve %d", rw.Code)
	}
}

func TestPasajeroHandler_UpdateDeleteGetList(t *testing.T) {
	h := &PasajeroHandler{service: &fakeService{
		UpdateFn: func(id int64, in input.UpdatePasajeroInput) (input.PasajeroOutput, error) {
			return input.PasajeroOutput{IDPasajero: id, Nombres: in.Nombres}, nil
		},
		DeleteFn: func(id int64) error { return nil },
		GetByIDFn: func(id int64) (input.PasajeroOutput, error) {
			if id == 404 {
				return input.PasajeroOutput{}, errors.New("not found")
			}
			return input.PasajeroOutput{IDPasajero: id, Nombres: "Juan"}, nil
		},
		ListFn: func(page, size int) ([]input.PasajeroOutput, pkg.PaginationMeta, error) {
			return []input.PasajeroOutput{{IDPasajero: 1}}, pkg.PaginationMeta{Page: page, Size: size, Total: 1}, nil
		},
	}}

	t.Run("update id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/pasajero/abc", bytes.NewReader([]byte(`{"nombres":"A"}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperaba status 400, obtuve %d", rw.Code)
		}
	})

	t.Run("update ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/pasajero/1", bytes.NewReader([]byte(`{"nombres":"Luis"}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rw := httptest.NewRecorder()
		h.Update(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("esperaba status 200, obtuve %d", rw.Code)
		}
	})

	t.Run("delete missing id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/pasajero", nil)
		rw := httptest.NewRecorder()
		h.Delete(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperaba status 400, obtuve %d", rw.Code)
		}
	})

	t.Run("get not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pasajero/404", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "404"})
		rw := httptest.NewRecorder()
		h.GetByID(rw, req)
		if rw.Code != http.StatusNotFound {
			t.Fatalf("esperaba status 404, obtuve %d", rw.Code)
		}
	})

	t.Run("list paginacion invalida", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pasajero?page=x", nil)
		rw := httptest.NewRecorder()
		h.List(rw, req)
		if rw.Code != http.StatusBadRequest {
			t.Fatalf("esperaba status 400, obtuve %d", rw.Code)
		}
	})

	t.Run("list ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pasajero?page=1&size=2", nil)
		rw := httptest.NewRecorder()
		h.List(rw, req)
		if rw.Code != http.StatusOK {
			t.Fatalf("esperaba status 200, obtuve %d", rw.Code)
		}
	})
}

func TestPasajeroHandler_Search_TipoDocumentoVacio(t *testing.T) {
	svc := &fakeService{
		SearchFn: func(query string) ([]input.PasajeroOutput, error) {
			return []input.PasajeroOutput{
				{IDPasajero: 1, TipoDocumento: "", NroDocumento: "73000001", Nombres: "Lucia", Apellidos: "Herrera Paz", Telefono: "987500001", Email: ptrStr("lucia.herrera@gmail.com"), FechaNacimiento: ptrStr("1995-05-10"), CreatedAt: "2026-03-26T14:55:38-05:00", UpdatedAt: "2026-03-26T14:55:38-05:00"},
				{IDPasajero: 7, TipoDocumento: "", NroDocumento: "12345600", Nombres: "Emily Latiana", Apellidos: "Benz", Telefono: "987654111", Email: ptrStr("emi@email.com"), FechaNacimiento: ptrStr("1990-05-14"), CreatedAt: "2026-03-27T12:16:04-05:00", UpdatedAt: "2026-03-27T12:16:04-05:00"},
			}, nil
		},
	}
	h := &PasajeroHandler{service: svc}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/pasajeros/search?q=Juan&size=2", nil)
	rw := httptest.NewRecorder()
	h.Search(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("esperaba status 200, obtuve %d", rw.Code)
	}
	var resp struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    []input.PasajeroOutput `json:"data"`
	}
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error decodificando respuesta: %v", err)
	}
	for _, p := range resp.Data {
		if p.TipoDocumento != "" {
			t.Errorf("esperaba tipo_documento vacio, obtuve '%s'", p.TipoDocumento)
		}
	}
}

func TestPasajeroHandler_Search_Errores(t *testing.T) {
	h := &PasajeroHandler{service: &fakeService{SearchFn: func(string) ([]input.PasajeroOutput, error) {
		return nil, errors.New("boom")
	}}}

	reqMissing := httptest.NewRequest(http.MethodGet, "/api/v1/pasajeros/search", nil)
	rwMissing := httptest.NewRecorder()
	h.Search(rwMissing, reqMissing)
	if rwMissing.Code != http.StatusBadRequest {
		t.Fatalf("esperaba status 400, obtuve %d", rwMissing.Code)
	}

	reqErr := httptest.NewRequest(http.MethodGet, "/api/v1/pasajeros/search?q=juan", nil)
	rwErr := httptest.NewRecorder()
	h.Search(rwErr, reqErr)
	if rwErr.Code != http.StatusInternalServerError {
		t.Fatalf("esperaba status 500, obtuve %d", rwErr.Code)
	}
}

func ptrStr(s string) *string { return &s }
