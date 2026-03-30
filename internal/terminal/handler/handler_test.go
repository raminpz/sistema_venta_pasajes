package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sistema_venta_pasajes/internal/terminal/domain"
	"sistema_venta_pasajes/internal/terminal/input"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) Create(in input.CreateTerminalInput) (*domain.Terminal, error) {
	args := m.Called(in)
	return args.Get(0).(*domain.Terminal), args.Error(1)
}
func (m *mockService) GetByID(id int64) (*domain.Terminal, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Terminal), args.Error(1)
}
func (m *mockService) Update(id int64, in input.UpdateTerminalInput) (*domain.Terminal, error) {
	args := m.Called(id, in)
	return args.Get(0).(*domain.Terminal), args.Error(1)
}
func (m *mockService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockService) List() ([]domain.Terminal, error) {
	args := m.Called()
	return args.Get(0).([]domain.Terminal), args.Error(1)
}

func TestCreateTerminalHandler(t *testing.T) {
	ms := new(mockService)
	h := NewTerminalHandler(ms)
	r := mux.NewRouter()
	registerRoutesWithHandler(r, h)
	terminalInput := input.CreateTerminalInput{
		Nombre:       "Terminal Test",
		Ciudad:       "Ciudad Test",
		Departamento: "Depto Test",
		Direccion:    "Dir Test",
		Estado:       "Activo",
	}
	ms.On("Create", terminalInput).Return(&domain.Terminal{NOMBRE: "Terminal Test"}, nil)
	body, _ := json.Marshal(terminalInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/terminal", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
