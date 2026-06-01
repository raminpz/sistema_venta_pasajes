package repository

import (
	"sistema_venta_pasajes/internal/tipo_vehiculo/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// mockTipoVehiculoDB representa una persistencia en memoria para validar contrato del repositorio.
type mockTipoVehiculoDB struct {
	items  map[int64]*domain.TipoVehiculo
	nextID int64
}

func newMockTipoVehiculoDB() *mockTipoVehiculoDB {
	return &mockTipoVehiculoDB{items: make(map[int64]*domain.TipoVehiculo), nextID: 1}
}

func (m *mockTipoVehiculoDB) Create(tv *domain.TipoVehiculo) error {
	tv.IDTipoVehiculo = m.nextID
	cp := *tv
	m.items[m.nextID] = &cp
	m.nextID++
	return nil
}
func (m *mockTipoVehiculoDB) Update(tv *domain.TipoVehiculo) error {
	if _, ok := m.items[tv.IDTipoVehiculo]; !ok {
		return gorm.ErrRecordNotFound
	}
	cp := *tv
	m.items[tv.IDTipoVehiculo] = &cp
	return nil
}
func (m *mockTipoVehiculoDB) Delete(id int64) error {
	if _, ok := m.items[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.items, id)
	return nil
}
func (m *mockTipoVehiculoDB) GetByID(id int64) (*domain.TipoVehiculo, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *item
	return &cp, nil
}
func (m *mockTipoVehiculoDB) List() ([]domain.TipoVehiculo, error) {
	list := make([]domain.TipoVehiculo, 0, len(m.items))
	for _, item := range m.items {
		list = append(list, *item)
	}
	return list, nil
}

func TestTipoVehiculoRepositoryContract(t *testing.T) {
	db := newMockTipoVehiculoDB()

	tv := &domain.TipoVehiculo{Nombre: "Auto", Descripcion: "Ejecutivo"}
	err := db.Create(tv)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), tv.IDTipoVehiculo)

	found, err := db.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Auto", found.Nombre)

	found.Descripcion = "Interprovincial"
	err = db.Update(found)
	assert.NoError(t, err)

	list, err := db.List()
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	err = db.Delete(1)
	assert.NoError(t, err)

	_, err = db.GetByID(1)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

