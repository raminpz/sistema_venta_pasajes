package handler

import (
  "bytes"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"

  "sistema_venta_pasajes/internal/pasajero/input"
  "sistema_venta_pasajes/pkg"
)

type fakeService struct {
  CreateFn func(input.CreatePasajeroInput) (input.PasajeroOutput, error)
  SearchFn func(string) ([]input.PasajeroOutput, error)
}

func (f *fakeService) Create(in input.CreatePasajeroInput) (input.PasajeroOutput, error) {
  return f.CreateFn(in)
}
func (f *fakeService) Update(id int64, in input.UpdatePasajeroInput) (input.PasajeroOutput, error) {
  return input.PasajeroOutput{}, nil
}
func (f *fakeService) Delete(id int64) error { return nil }
func (f *fakeService) GetByID(id int64) (input.PasajeroOutput, error) {
  return input.PasajeroOutput{}, nil
}
func (f *fakeService) List(page, size int) ([]input.PasajeroOutput, pkg.PaginationMeta, error) {
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
    Code    int                  `json:"code"`
    Message string               `json:"message"`
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

func ptrStr(s string) *string { return &s }
