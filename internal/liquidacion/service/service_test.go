package service
import (
"errors"
"testing"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/mock"
"sistema_venta_pasajes/internal/liquidacion/domain"
"sistema_venta_pasajes/internal/liquidacion/input"
)
// ── Mock del repositorio ──────────────────────────────────────────────────────
type MockLiquidacionRepository struct {
mock.Mock
}
func (m *MockLiquidacionRepository) Create(liq *domain.LiquidacionViaje) error {
args := m.Called(liq)
return args.Error(0)
}
func (m *MockLiquidacionRepository) Update(liq *domain.LiquidacionViaje) error {
args := m.Called(liq)
return args.Error(0)
}
func (m *MockLiquidacionRepository) Delete(id int64) error {
args := m.Called(id)
return args.Error(0)
}
func (m *MockLiquidacionRepository) GetByID(id int64) (*domain.LiquidacionViaje, error) {
args := m.Called(id)
v := args.Get(0)
if v == nil {
return nil, args.Error(1)
}
return v.(*domain.LiquidacionViaje), args.Error(1)
}
func (m *MockLiquidacionRepository) GetByProgramacion(id int64) (*domain.LiquidacionViaje, error) {
args := m.Called(id)
v := args.Get(0)
if v == nil {
return nil, args.Error(1)
}
return v.(*domain.LiquidacionViaje), args.Error(1)
}
func (m *MockLiquidacionRepository) List(offset, limit int) ([]domain.LiquidacionViaje, int, error) {
args := m.Called(offset, limit)
return args.Get(0).([]domain.LiquidacionViaje), args.Int(1), args.Error(2)
}
func (m *MockLiquidacionRepository) ExistsByProgramacion(id int64) (bool, error) {
args := m.Called(id)
return args.Bool(0), args.Error(1)
}
func (m *MockLiquidacionRepository) GetConductorByProgramacion(id int64) (int64, error) {
args := m.Called(id)
return args.Get(0).(int64), args.Error(1)
}
func (m *MockLiquidacionRepository) SumarVentas(id int64) (float64, int, error) {
args := m.Called(id)
return args.Get(0).(float64), args.Int(1), args.Error(2)
}
func (m *MockLiquidacionRepository) SumarEncomiendas(id int64) (float64, int, error) {
args := m.Called(id)
return args.Get(0).(float64), args.Int(1), args.Error(2)
}
// ── Tests ─────────────────────────────────────────────────────────────────────
func TestGenerar_OK(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
mockRepo.On("ExistsByProgramacion", int64(1)).Return(false, nil)
mockRepo.On("GetConductorByProgramacion", int64(1)).Return(int64(5), nil)
mockRepo.On("SumarVentas", int64(1)).Return(float64(450.00), 9, nil)
mockRepo.On("SumarEncomiendas", int64(1)).Return(float64(80.00), 3, nil)
mockRepo.On("Create", mock.Anything).Return(nil)
out, err := svc.Generar(input.GenerarLiquidacionInput{IDProgramacion: 1, Observaciones: "Viaje Lima-Huancayo"})
assert.NoError(t, err)
assert.NotNil(t, out)
assert.Equal(t, float64(530.00), out.TotalCaja)
assert.Equal(t, float64(450.00), out.TotalPasajes)
assert.Equal(t, float64(80.00), out.TotalEncomiendas)
assert.Equal(t, "PENDIENTE", out.Estado)
assert.Equal(t, 9, out.CantidadPasajes)
assert.Equal(t, 3, out.CantidadEncomiendas)
}
func TestGenerar_IDProgramacionInvalido(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
out, err := svc.Generar(input.GenerarLiquidacionInput{IDProgramacion: 0})
assert.Error(t, err)
assert.Nil(t, out)
}
func TestGenerar_LiquidacionYaExiste(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
mockRepo.On("ExistsByProgramacion", int64(2)).Return(true, nil)
out, err := svc.Generar(input.GenerarLiquidacionInput{IDProgramacion: 2})
assert.Error(t, err)
assert.Nil(t, out)
}
func TestGenerar_ProgramacionNoExiste(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
mockRepo.On("ExistsByProgramacion", int64(99)).Return(false, nil)
mockRepo.On("GetConductorByProgramacion", int64(99)).Return(int64(0), errors.New("not found"))
out, err := svc.Generar(input.GenerarLiquidacionInput{IDProgramacion: 99})
assert.Error(t, err)
assert.Nil(t, out)
}
func TestActualizarEstado_AEntregado(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
liq := &domain.LiquidacionViaje{IDLiquidacion: 1, Estado: "PENDIENTE", TotalCaja: 530.00}
mockRepo.On("GetByID", int64(1)).Return(liq, nil)
mockRepo.On("Update", mock.Anything).Return(nil)
out, err := svc.ActualizarEstado(1, input.ActualizarEstadoInput{Estado: "ENTREGADO"})
assert.NoError(t, err)
assert.Equal(t, "ENTREGADO", out.Estado)
assert.NotNil(t, out.FechaLiquidacion)
}
func TestActualizarEstado_EstadoInvalido(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
out, err := svc.ActualizarEstado(1, input.ActualizarEstadoInput{Estado: "INVALIDO"})
assert.Error(t, err)
assert.Nil(t, out)
}
func TestActualizarEstado_NotFound(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
mockRepo.On("GetByID", int64(999)).Return(nil, errors.New("not found"))
out, err := svc.ActualizarEstado(999, input.ActualizarEstadoInput{Estado: "ENTREGADO"})
assert.Error(t, err)
assert.Nil(t, out)
}
func TestGetByID_OK(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
liq := &domain.LiquidacionViaje{IDLiquidacion: 3, TotalCaja: 200.00, Estado: "PENDIENTE"}
mockRepo.On("GetByID", int64(3)).Return(liq, nil)
out, err := svc.GetByID(3)
assert.NoError(t, err)
assert.Equal(t, int64(3), out.IDLiquidacion)
}
func TestGetByID_NotFound(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
mockRepo.On("GetByID", int64(999)).Return(nil, errors.New("record not found"))
out, err := svc.GetByID(999)
assert.Error(t, err)
assert.Nil(t, out)
}
func TestDelete_OK(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
liq := &domain.LiquidacionViaje{IDLiquidacion: 1}
mockRepo.On("GetByID", int64(1)).Return(liq, nil)
mockRepo.On("Delete", int64(1)).Return(nil)
err := svc.Delete(1)
assert.NoError(t, err)
}
func TestDelete_NotFound(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
mockRepo.On("GetByID", int64(99)).Return(nil, errors.New("not found"))
err := svc.Delete(99)
assert.Error(t, err)
}
func TestList_OK(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
liqs := []domain.LiquidacionViaje{
{IDLiquidacion: 1, TotalCaja: 100},
{IDLiquidacion: 2, TotalCaja: 200},
}
mockRepo.On("List", 0, 15).Return(liqs, 2, nil)
out, total, err := svc.List(1, 15)
assert.NoError(t, err)
assert.Equal(t, 2, total)
assert.Len(t, out, 2)
}
func TestList_Vacio(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
mockRepo.On("List", 0, 15).Return([]domain.LiquidacionViaje{}, 0, nil)
out, total, err := svc.List(1, 15)
assert.NoError(t, err)
assert.Equal(t, 0, total)
assert.Empty(t, out)
}
func TestObtenerResumenCaja_OK(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
mockRepo.On("GetConductorByProgramacion", int64(5)).Return(int64(2), nil)
mockRepo.On("SumarVentas", int64(5)).Return(float64(300.00), 6, nil)
mockRepo.On("SumarEncomiendas", int64(5)).Return(float64(50.00), 2, nil)
out, err := svc.ObtenerResumenCaja(5)
assert.NoError(t, err)
assert.Equal(t, float64(350.00), out.TotalCaja)
assert.Equal(t, 6, out.CantidadPasajes)
assert.Equal(t, 2, out.CantidadEncomiendas)
}
func TestObtenerResumenCaja_ProgramacionNoExiste(t *testing.T) {
mockRepo := new(MockLiquidacionRepository)
svc := NewLiquidacionService(mockRepo)
mockRepo.On("GetConductorByProgramacion", int64(99)).Return(int64(0), errors.New("not found"))
out, err := svc.ObtenerResumenCaja(99)
assert.Error(t, err)
assert.Nil(t, out)
}
