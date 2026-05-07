package handler
import (
"bytes"
"encoding/json"
"errors"
"net/http"
"net/http/httptest"
"testing"
"github.com/gorilla/mux"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/mock"
"sistema_venta_pasajes/internal/liquidacion/input"
"sistema_venta_pasajes/internal/liquidacion/util"
"sistema_venta_pasajes/pkg"
)
// ── Mock del servicio ─────────────────────────────────────────────────────────
type MockLiquidacionService struct {
mock.Mock
}
func (m *MockLiquidacionService) Generar(in input.GenerarLiquidacionInput) (*input.LiquidacionOutput, error) {
args := m.Called(in)
v := args.Get(0)
if v == nil {
return nil, args.Error(1)
}
return v.(*input.LiquidacionOutput), args.Error(1)
}
func (m *MockLiquidacionService) ActualizarEstado(id int64, in input.ActualizarEstadoInput) (*input.LiquidacionOutput, error) {
args := m.Called(id, in)
v := args.Get(0)
if v == nil {
return nil, args.Error(1)
}
return v.(*input.LiquidacionOutput), args.Error(1)
}
func (m *MockLiquidacionService) Delete(id int64) error {
args := m.Called(id)
return args.Error(0)
}
func (m *MockLiquidacionService) GetByID(id int64) (*input.LiquidacionOutput, error) {
args := m.Called(id)
v := args.Get(0)
if v == nil {
return nil, args.Error(1)
}
return v.(*input.LiquidacionOutput), args.Error(1)
}
func (m *MockLiquidacionService) List(page, size int) ([]input.LiquidacionOutput, int, error) {
args := m.Called(page, size)
return args.Get(0).([]input.LiquidacionOutput), args.Int(1), args.Error(2)
}
func (m *MockLiquidacionService) ObtenerResumenCaja(idProgramacion int64) (*input.ResumenCajaOutput, error) {
args := m.Called(idProgramacion)
v := args.Get(0)
if v == nil {
return nil, args.Error(1)
}
return v.(*input.ResumenCajaOutput), args.Error(1)
}
// ── Tests ─────────────────────────────────────────────────────────────────────
func TestHandler_Generar_OK(t *testing.T) {
mockSvc := new(MockLiquidacionService)
h := NewLiquidacionHandler(mockSvc)
in := input.GenerarLiquidacionInput{IDProgramacion: 1, Observaciones: "Viaje Lima-Ayacucho"}
out := &input.LiquidacionOutput{IDLiquidacion: 1, TotalCaja: 530.00, Estado: "PENDIENTE"}
mockSvc.On("Generar", in).Return(out, nil)
body, _ := json.Marshal(in)
req := httptest.NewRequest(http.MethodPost, "/liquidacion", bytes.NewReader(body))
w := httptest.NewRecorder()
h.Generar(w, req)
assert.Equal(t, http.StatusCreated, w.Code)
var resp map[string]interface{}
json.Unmarshal(w.Body.Bytes(), &resp)
assert.Equal(t, util.MSG_CREATED, resp["message"])
}
func TestHandler_Generar_ServiceError(t *testing.T) {
mockSvc := new(MockLiquidacionService)
h := NewLiquidacionHandler(mockSvc)
in := input.GenerarLiquidacionInput{IDProgramacion: 2}
mockSvc.On("Generar", in).Return(nil, errors.New("ya existe"))
body, _ := json.Marshal(in)
req := httptest.NewRequest(http.MethodPost, "/liquidacion", bytes.NewReader(body))
w := httptest.NewRecorder()
h.Generar(w, req)
assert.NotEqual(t, http.StatusCreated, w.Code)
}
func TestHandler_GetByID_OK(t *testing.T) {
mockSvc := new(MockLiquidacionService)
h := NewLiquidacionHandler(mockSvc)
out := &input.LiquidacionOutput{IDLiquidacion: 5, TotalCaja: 200.00, Estado: "PENDIENTE"}
mockSvc.On("GetByID", int64(5)).Return(out, nil)
req := httptest.NewRequest(http.MethodGet, "/liquidacion/5", nil)
r := mux.NewRouter()
r.HandleFunc("/liquidacion/{id}", h.GetByID).Methods("GET")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
var resp map[string]interface{}
json.Unmarshal(w.Body.Bytes(), &resp)
assert.Equal(t, util.MSG_GET, resp["message"])
}
func TestHandler_GetByID_NotFound(t *testing.T) {
mockSvc := new(MockLiquidacionService)
h := NewLiquidacionHandler(mockSvc)
mockSvc.On("GetByID", int64(99)).Return(nil, pkg.NotFound("liquidacion_not_found", util.ERR_NOT_FOUND))
req := httptest.NewRequest(http.MethodGet, "/liquidacion/99", nil)
r := mux.NewRouter()
r.HandleFunc("/liquidacion/{id}", h.GetByID).Methods("GET")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestHandler_ActualizarEstado_OK(t *testing.T) {
mockSvc := new(MockLiquidacionService)
h := NewLiquidacionHandler(mockSvc)
in := input.ActualizarEstadoInput{Estado: "ENTREGADO"}
out := &input.LiquidacionOutput{IDLiquidacion: 1, Estado: "ENTREGADO"}
mockSvc.On("ActualizarEstado", int64(1), in).Return(out, nil)
body, _ := json.Marshal(in)
req := httptest.NewRequest(http.MethodPut, "/liquidacion/1", bytes.NewReader(body))
r := mux.NewRouter()
r.HandleFunc("/liquidacion/{id}", h.ActualizarEstado).Methods("PUT")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
var resp map[string]interface{}
json.Unmarshal(w.Body.Bytes(), &resp)
assert.Equal(t, util.MSG_UPDATED, resp["message"])
}
func TestHandler_Delete_OK(t *testing.T) {
mockSvc := new(MockLiquidacionService)
h := NewLiquidacionHandler(mockSvc)
mockSvc.On("Delete", int64(3)).Return(nil)
req := httptest.NewRequest(http.MethodDelete, "/liquidacion/3", nil)
r := mux.NewRouter()
r.HandleFunc("/liquidacion/{id}", h.Delete).Methods("DELETE")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
}
func TestHandler_List_OK(t *testing.T) {
mockSvc := new(MockLiquidacionService)
h := NewLiquidacionHandler(mockSvc)
liqs := []input.LiquidacionOutput{{IDLiquidacion: 1}, {IDLiquidacion: 2}}
mockSvc.On("List", 1, 15).Return(liqs, 2, nil)
req := httptest.NewRequest(http.MethodGet, "/liquidaciones", nil)
w := httptest.NewRecorder()
h.List(w, req)
assert.Equal(t, http.StatusOK, w.Code)
}
func TestHandler_ObtenerResumenCaja_OK(t *testing.T) {
mockSvc := new(MockLiquidacionService)
h := NewLiquidacionHandler(mockSvc)
out := &input.ResumenCajaOutput{IDProgramacion: 7, TotalCaja: 350.00, CantidadPasajes: 6}
mockSvc.On("ObtenerResumenCaja", int64(7)).Return(out, nil)
req := httptest.NewRequest(http.MethodGet, "/programacion/7/caja", nil)
r := mux.NewRouter()
r.HandleFunc("/programacion/{id_programacion}/caja", h.ObtenerResumenCaja).Methods("GET")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
var resp map[string]interface{}
json.Unmarshal(w.Body.Bytes(), &resp)
assert.Equal(t, util.MSG_RESUMEN_CAJA, resp["message"])
}
