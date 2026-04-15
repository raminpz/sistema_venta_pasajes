package repository
import (
"errors"
"sistema_venta_pasajes/internal/liquidacion/domain"
"testing"
"github.com/stretchr/testify/assert"
)
type mockLiqRepo struct {
CreateFunc func(liq *domain.LiquidacionViaje) error
GetByIDFunc func(id int64) (*domain.LiquidacionViaje, error)
}
func (m *mockLiqRepo) Create(liq *domain.LiquidacionViaje) error {
if m.CreateFunc != nil {
return m.CreateFunc(liq)
}
return nil
}
func (m *mockLiqRepo) GetByID(id int64) (*domain.LiquidacionViaje, error) {
if m.GetByIDFunc != nil {
return m.GetByIDFunc(id)
}
return nil, errors.New("not found")
}
func TestLiquidacionRepository_Create_OK(t *testing.T) {
liq := &domain.LiquidacionViaje{
IDProgramacion:   1,
IDConductor:      2,
TotalPasajes:     450.00,
TotalEncomiendas: 80.00,
TotalCaja:        530.00,
Estado:           "PENDIENTE",
}
mockRepo := &mockLiqRepo{
CreateFunc: func(l *domain.LiquidacionViaje) error {
if l.IDProgramacion == 1 {
l.IDLiquidacion = 1
return nil
}
return errors.New("programacion invalida")
},
}
err := mockRepo.Create(liq)
assert.NoError(t, err)
assert.Equal(t, int64(1), liq.IDLiquidacion)
}
func TestLiquidacionRepository_GetByID_NotFound(t *testing.T) {
mockRepo := &mockLiqRepo{
GetByIDFunc: func(id int64) (*domain.LiquidacionViaje, error) {
return nil, errors.New("record not found")
},
}
result, err := mockRepo.GetByID(999)
assert.Error(t, err)
assert.Nil(t, result)
}
func TestLiquidacionRepository_Create_Error(t *testing.T) {
liq := &domain.LiquidacionViaje{IDProgramacion: 999}
mockRepo := &mockLiqRepo{
CreateFunc: func(l *domain.LiquidacionViaje) error {
return errors.New("constraint violation")
},
}
err := mockRepo.Create(liq)
assert.Error(t, err)
}
