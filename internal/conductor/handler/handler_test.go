package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"sistema_venta_pasajes/internal/conductor/domain"
	"sistema_venta_pasajes/internal/conductor/input"
)

type mockService struct{}

func (m *mockService) List(ctx context.Context) ([]domain.Conductor, error) {
	return []domain.Conductor{{IDConductor: 1, Nombres: "Juan", Apellidos: "Perez", NumeroLicencia: "ABC123456", Telefono: "987654321"}}, nil
}
func (m *mockService) GetByID(ctx context.Context, id int64) (*domain.Conductor, error) {
	return &domain.Conductor{IDConductor: id, Nombres: "Juan", Apellidos: "Perez", NumeroLicencia: "ABC123456", Telefono: "987654321"}, nil
}
func (m *mockService) Create(ctx context.Context, in input.CreateConductorInput) (*domain.Conductor, error) {
	return &domain.Conductor{IDConductor: 2, Nombres: in.Nombres, Apellidos: in.Apellidos, NumeroLicencia: in.NumeroLicencia, Telefono: in.Telefono}, nil
}
func (m *mockService) Update(ctx context.Context, id int64, in input.UpdateConductorInput) (*domain.Conductor, error) {
	return &domain.Conductor{IDConductor: id, Nombres: "Mod", Apellidos: "Mod", NumeroLicencia: "MOD123456", Telefono: "987654321"}, nil
}
func (m *mockService) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestHandler_Create(t *testing.T) {
	h := New(&mockService{})
	r := mux.NewRouter()
	r.HandleFunc("/conductor", h.Create).Methods("POST")
	body := input.CreateConductorInput{
		Nombres:      "Juan",
		Apellidos:    "Perez",
		NumeroLicencia: "ABC123456",
		Telefono:     "987654321",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/conductor", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}
