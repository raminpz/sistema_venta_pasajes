package repository

import (
	"testing"

	"sistema_venta_pasajes/internal/asiento_tramo/domain"
	"sistema_venta_pasajes/internal/asiento_tramo/util"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// mockAsientoTramoDB simula persistencia en memoria para pruebas de contrato del repositorio.
type mockAsientoTramoDB struct {
	data   map[int64]*domain.AsientoTramo
	nextID int64
}

func newMockAsientoTramoDB() *mockAsientoTramoDB {
	return &mockAsientoTramoDB{data: make(map[int64]*domain.AsientoTramo), nextID: 1}
}

func (m *mockAsientoTramoDB) Create(at *domain.AsientoTramo) error {
	at.IDAsientoTramo = m.nextID
	cp := *at
	m.data[m.nextID] = &cp
	m.nextID++
	return nil
}

func (m *mockAsientoTramoDB) Update(at *domain.AsientoTramo) error {
	if _, ok := m.data[at.IDAsientoTramo]; !ok {
		return gorm.ErrRecordNotFound
	}
	cp := *at
	m.data[at.IDAsientoTramo] = &cp
	return nil
}

func (m *mockAsientoTramoDB) Delete(id int64) error {
	if _, ok := m.data[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.data, id)
	return nil
}

func (m *mockAsientoTramoDB) GetByID(id int64) (*domain.AsientoTramo, error) {
	at, ok := m.data[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *at
	return &cp, nil
}

func (m *mockAsientoTramoDB) GetByAsientoTramo(idAsiento, idTramo int64) (*domain.AsientoTramo, error) {
	for _, at := range m.data {
		if at.IDAsiento == idAsiento && at.IDTramo == idTramo {
			cp := *at
			return &cp, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAsientoTramoDB) GetDisponiblesEnTramo(idTramo int64) ([]domain.AsientoTramo, error) {
	res := make([]domain.AsientoTramo, 0)
	for _, at := range m.data {
		if at.IDTramo == idTramo && at.Estado == util.ESTADO_DISPONIBLE {
			res = append(res, *at)
		}
	}
	return res, nil
}

func (m *mockAsientoTramoDB) MarkAsOccupied(idAsiento, idTramo int64, idVenta *int64) error {
	for _, at := range m.data {
		if at.IDAsiento == idAsiento && at.IDTramo == idTramo {
			at.Estado = util.ESTADO_OCUPADO
			at.IDVenta = idVenta
		}
	}
	return nil
}

func (m *mockAsientoTramoDB) MarkAsAvailable(idAsiento, idTramo int64) error {
	for _, at := range m.data {
		if at.IDAsiento == idAsiento && at.IDTramo == idTramo {
			at.Estado = util.ESTADO_DISPONIBLE
			at.IDVenta = nil
		}
	}
	return nil
}

func (m *mockAsientoTramoDB) DeleteByVenta(idVenta int64) error {
	for id, at := range m.data {
		if at.IDVenta != nil && *at.IDVenta == idVenta {
			delete(m.data, id)
		}
	}
	return nil
}

func TestNewAsientoTramoRepository_NotNil(t *testing.T) {
	repo := NewAsientoTramoRepository(nil)
	assert.NotNil(t, repo)
}

func TestAsientoTramoRepository_CreateAndGetByID(t *testing.T) {
	db := newMockAsientoTramoDB()
	at := &domain.AsientoTramo{IDAsiento: 10, IDTramo: 2, Estado: util.ESTADO_DISPONIBLE}
	err := db.Create(at)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), at.IDAsientoTramo)

	got, err := db.GetByID(at.IDAsientoTramo)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), got.IDAsiento)
	assert.Equal(t, int64(2), got.IDTramo)
}

func TestAsientoTramoRepository_GetByID_NotFound(t *testing.T) {
	db := newMockAsientoTramoDB()
	_, err := db.GetByID(999)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestAsientoTramoRepository_Update_OK(t *testing.T) {
	db := newMockAsientoTramoDB()
	_ = db.Create(&domain.AsientoTramo{IDAsiento: 1, IDTramo: 1, Estado: util.ESTADO_DISPONIBLE})

	err := db.Update(&domain.AsientoTramo{IDAsientoTramo: 1, IDAsiento: 1, IDTramo: 1, Estado: util.ESTADO_OCUPADO})
	assert.NoError(t, err)

	got, _ := db.GetByID(1)
	assert.Equal(t, util.ESTADO_OCUPADO, got.Estado)
}

func TestAsientoTramoRepository_Update_NotFound(t *testing.T) {
	db := newMockAsientoTramoDB()
	err := db.Update(&domain.AsientoTramo{IDAsientoTramo: 500})
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestAsientoTramoRepository_Delete_OK(t *testing.T) {
	db := newMockAsientoTramoDB()
	_ = db.Create(&domain.AsientoTramo{IDAsiento: 2, IDTramo: 3, Estado: util.ESTADO_DISPONIBLE})

	err := db.Delete(1)
	assert.NoError(t, err)
	_, err = db.GetByID(1)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestAsientoTramoRepository_Delete_NotFound(t *testing.T) {
	db := newMockAsientoTramoDB()
	err := db.Delete(404)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestAsientoTramoRepository_GetByAsientoTramo(t *testing.T) {
	db := newMockAsientoTramoDB()
	_ = db.Create(&domain.AsientoTramo{IDAsiento: 7, IDTramo: 9, Estado: util.ESTADO_DISPONIBLE})

	got, err := db.GetByAsientoTramo(7, 9)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), got.IDAsiento)
	assert.Equal(t, int64(9), got.IDTramo)

	_, err = db.GetByAsientoTramo(7, 10)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestAsientoTramoRepository_GetDisponiblesEnTramo(t *testing.T) {
	db := newMockAsientoTramoDB()
	_ = db.Create(&domain.AsientoTramo{IDAsiento: 1, IDTramo: 5, Estado: util.ESTADO_DISPONIBLE})
	_ = db.Create(&domain.AsientoTramo{IDAsiento: 2, IDTramo: 5, Estado: util.ESTADO_OCUPADO})
	_ = db.Create(&domain.AsientoTramo{IDAsiento: 3, IDTramo: 5, Estado: util.ESTADO_DISPONIBLE})

	list, err := db.GetDisponiblesEnTramo(5)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestAsientoTramoRepository_MarkAsOccupiedAndAvailable(t *testing.T) {
	db := newMockAsientoTramoDB()
	_ = db.Create(&domain.AsientoTramo{IDAsiento: 11, IDTramo: 8, Estado: util.ESTADO_DISPONIBLE})
	idVenta := int64(1001)

	err := db.MarkAsOccupied(11, 8, &idVenta)
	assert.NoError(t, err)
	got, _ := db.GetByID(1)
	assert.Equal(t, util.ESTADO_OCUPADO, got.Estado)
	assert.NotNil(t, got.IDVenta)
	assert.Equal(t, idVenta, *got.IDVenta)

	err = db.MarkAsAvailable(11, 8)
	assert.NoError(t, err)
	got, _ = db.GetByID(1)
	assert.Equal(t, util.ESTADO_DISPONIBLE, got.Estado)
	assert.Nil(t, got.IDVenta)
}

func TestAsientoTramoRepository_DeleteByVenta(t *testing.T) {
	db := newMockAsientoTramoDB()
	idVenta := int64(200)
	_ = db.Create(&domain.AsientoTramo{IDAsiento: 1, IDTramo: 1, Estado: util.ESTADO_OCUPADO, IDVenta: &idVenta})
	_ = db.Create(&domain.AsientoTramo{IDAsiento: 2, IDTramo: 1, Estado: util.ESTADO_DISPONIBLE})

	err := db.DeleteByVenta(idVenta)
	assert.NoError(t, err)
	assert.Len(t, db.data, 1)
	_, exists := db.data[2]
	assert.True(t, exists)
}
