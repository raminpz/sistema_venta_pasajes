package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/vehiculo/util"
	"testing"

	vehiculoinput "sistema_venta_pasajes/internal/vehiculo/input"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVehiculoService struct {
	mock.Mock
}

func (m *MockVehiculoService) Create(input vehiculoinput.CreateVehiculoInput) (*vehiculoinput.VehiculoOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*vehiculoinput.VehiculoOutput), args.Error(1)
}
func (m *MockVehiculoService) Update(input vehiculoinput.UpdateVehiculoInput) (*vehiculoinput.VehiculoOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*vehiculoinput.VehiculoOutput), args.Error(1)
}
func (m *MockVehiculoService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockVehiculoService) GetByID(id int64) (*vehiculoinput.VehiculoOutput, error) {
	args := m.Called(id)
	return args.Get(0).(*vehiculoinput.VehiculoOutput), args.Error(1)
}
func (m *MockVehiculoService) List(page, size int) ([]vehiculoinput.VehiculoOutput, int, error) {
	args := m.Called(page, size)
	return args.Get(0).([]vehiculoinput.VehiculoOutput), args.Int(1), args.Error(2)
}

func TestVehiculoHandler_Create(t *testing.T) {
	mockService := new(MockVehiculoService)
	handler := NewVehiculoHandler(mockService)
	vehiculoInput := vehiculoinput.CreateVehiculoInput{NroPlaca: "ABC123"}
	output := &vehiculoinput.VehiculoOutput{IDVehiculo: 1, NroPlaca: "ABC123"}
	mockService.On("Create", vehiculoInput).Return(output, nil)
	body, _ := json.Marshal(vehiculoInput)
	req := httptest.NewRequest(http.MethodPost, "/vehiculo", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.Create(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestVehiculoHandler_GetByID(t *testing.T) {
	mockService := new(MockVehiculoService)
	handler := NewVehiculoHandler(mockService)
	output := &vehiculoinput.VehiculoOutput{IDVehiculo: 5, NroPlaca: "XYZ123"}
	mockService.On("GetByID", int64(5)).Return(output, nil)

	req := httptest.NewRequest(http.MethodGet, "/vehiculo/5", nil)
	r := mux.NewRouter()
	r.HandleFunc("/vehiculo/{id}", handler.GetByID).Methods("GET")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, util.MSG_GET, resp["message"])
	assert.NotNil(t, resp["data"])
}

func TestVehiculoHandler_Delete(t *testing.T) {
	mockService := new(MockVehiculoService)
	handler := NewVehiculoHandler(mockService)
	mockService.On("Delete", int64(19)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/vehiculo/19", nil)
	r := mux.NewRouter()
	r.HandleFunc("/vehiculo/{id}", handler.Delete).Methods("DELETE")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, util.MSG_DELETED, resp["message"])
}

func TestVehiculoHandler_Delete_InvalidID(t *testing.T) {
	mockService := new(MockVehiculoService)
	handler := NewVehiculoHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/vehiculo/abc", nil)
	r := mux.NewRouter()
	r.HandleFunc("/vehiculo/{id}", handler.Delete).Methods("DELETE")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVehiculoHandler_Update(t *testing.T) {
	mockService := new(MockVehiculoService)
	handler := NewVehiculoHandler(mockService)

	in := vehiculoinput.UpdateVehiculoInput{IDVehiculo: 5, NroPlaca: "ABC-123"}
	output := &vehiculoinput.VehiculoOutput{IDVehiculo: 5, NroPlaca: "ABC-123"}
	mockService.On("Update", in).Return(output, nil)

	body, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPut, "/vehiculo/5", bytes.NewReader(body))
	r := mux.NewRouter()
	r.HandleFunc("/vehiculo/{id}", handler.Update).Methods("PUT")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, util.MSG_UPDATED, resp["message"])
}
