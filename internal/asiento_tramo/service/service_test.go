package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"sistema_venta_pasajes/internal/asiento_tramo/domain"
	"sistema_venta_pasajes/internal/asiento_tramo/input"
	"sistema_venta_pasajes/internal/asiento_tramo/util"
)

// ── Mock ─────────────────────────────────────────────────────────────────

type mockAsientoTramoRepo struct {
	data   map[int64]*domain.AsientoTramo
	nextID int64
}

func newMockRepo() *mockAsientoTramoRepo {
	return &mockAsientoTramoRepo{
		data:   make(map[int64]*domain.AsientoTramo),
		nextID: 1,
	}
}

func (m *mockAsientoTramoRepo) Create(at *domain.AsientoTramo) error {
	at.IDAsientoTramo = m.nextID
	m.nextID++
	m.data[at.IDAsientoTramo] = at
	return nil
}

func (m *mockAsientoTramoRepo) Update(at *domain.AsientoTramo) error {
	if _, ok := m.data[at.IDAsientoTramo]; !ok {
		return gorm.ErrRecordNotFound
	}
	m.data[at.IDAsientoTramo] = at
	return nil
}

func (m *mockAsientoTramoRepo) Delete(id int64) error {
	if _, ok := m.data[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.data, id)
	return nil
}

func (m *mockAsientoTramoRepo) GetByID(id int64) (*domain.AsientoTramo, error) {
	if at, ok := m.data[id]; ok {
		return at, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAsientoTramoRepo) GetByAsientoTramo(idAsiento, idTramo int64) (*domain.AsientoTramo, error) {
	for _, at := range m.data {
		if at.IDAsiento == idAsiento && at.IDTramo == idTramo {
			return at, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAsientoTramoRepo) GetDisponiblesEnTramo(idTramo int64) ([]domain.AsientoTramo, error) {
	var result []domain.AsientoTramo
	for _, at := range m.data {
		if at.IDTramo == idTramo && at.Estado == util.ESTADO_DISPONIBLE {
			result = append(result, *at)
		}
	}
	return result, nil
}

func (m *mockAsientoTramoRepo) MarkAsOccupied(idAsiento, idTramo int64, idVenta *int64) error {
	for _, at := range m.data {
		if at.IDAsiento == idAsiento && at.IDTramo == idTramo {
			at.Estado = util.ESTADO_OCUPADO
			at.IDVenta = idVenta
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockAsientoTramoRepo) MarkAsAvailable(idAsiento, idTramo int64) error {
	for _, at := range m.data {
		if at.IDAsiento == idAsiento && at.IDTramo == idTramo {
			at.Estado = util.ESTADO_DISPONIBLE
			at.IDVenta = nil
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockAsientoTramoRepo) DeleteByVenta(idVenta int64) error {
	for id, at := range m.data {
		if at.IDVenta != nil && *at.IDVenta == idVenta {
			delete(m.data, id)
		}
	}
	return nil
}

// ── Tests ────────────────────────────────────────────────────────────────

func TestAsientoTramo_Create_OK(t *testing.T) {
	svc := NewAsientoTramoService(newMockRepo())
	out, err := svc.Create(input.CreateAsientoTramoInput{
		IDAsiento: 1,
		IDTramo:   1,
		Estado:    util.ESTADO_DISPONIBLE,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), out.IDAsientoTramo)
	assert.Equal(t, util.ESTADO_DISPONIBLE, out.Estado)
}

func TestAsientoTramo_Create_InvalidAsiento(t *testing.T) {
	svc := NewAsientoTramoService(newMockRepo())
	_, err := svc.Create(input.CreateAsientoTramoInput{
		IDAsiento: 0,
		IDTramo:   1,
		Estado:    util.ESTADO_DISPONIBLE,
	})
	assert.Error(t, err)
}

func TestAsientoTramo_Create_InvalidTramo(t *testing.T) {
	svc := NewAsientoTramoService(newMockRepo())
	_, err := svc.Create(input.CreateAsientoTramoInput{
		IDAsiento: 1,
		IDTramo:   0,
		Estado:    util.ESTADO_DISPONIBLE,
	})
	assert.Error(t, err)
}

func TestAsientoTramo_GetByID_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewAsientoTramoService(repo)
	_, _ = svc.Create(input.CreateAsientoTramoInput{
		IDAsiento: 5,
		IDTramo:   2,
		Estado:    util.ESTADO_DISPONIBLE,
	})
	out, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), out.IDAsiento)
	assert.Equal(t, int64(2), out.IDTramo)
}

func TestAsientoTramo_GetByID_NotFound(t *testing.T) {
	svc := NewAsientoTramoService(newMockRepo())
	_, err := svc.GetByID(999)
	assert.Error(t, err)
}

func TestAsientoTramo_MarkAsOccupied_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewAsientoTramoService(repo)
	_, _ = svc.Create(input.CreateAsientoTramoInput{
		IDAsiento: 5,
		IDTramo:   2,
		Estado:    util.ESTADO_DISPONIBLE,
	})

	idVenta := int64(100)
	err := svc.MarkAsOccupied(5, 2, &idVenta)
	assert.NoError(t, err)

	out, _ := svc.GetByID(1)
	assert.Equal(t, util.ESTADO_OCUPADO, out.Estado)
	assert.Equal(t, idVenta, *out.IDVenta)
}

func TestAsientoTramo_MarkAsAvailable_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewAsientoTramoService(repo)
	_, _ = svc.Create(input.CreateAsientoTramoInput{
		IDAsiento: 5,
		IDTramo:   2,
		Estado:    util.ESTADO_OCUPADO,
	})

	err := svc.MarkAsAvailable(5, 2)
	assert.NoError(t, err)

	out, _ := svc.GetByID(1)
	assert.Equal(t, util.ESTADO_DISPONIBLE, out.Estado)
	assert.Nil(t, out.IDVenta)
}

func TestAsientoTramo_GetDisponiblesEnTramo_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewAsientoTramoService(repo)
	_, _ = svc.Create(input.CreateAsientoTramoInput{IDAsiento: 1, IDTramo: 2, Estado: util.ESTADO_DISPONIBLE})
	_, _ = svc.Create(input.CreateAsientoTramoInput{IDAsiento: 2, IDTramo: 2, Estado: util.ESTADO_DISPONIBLE})
	_, _ = svc.Create(input.CreateAsientoTramoInput{IDAsiento: 3, IDTramo: 2, Estado: util.ESTADO_OCUPADO})

	list, err := svc.GetDisponiblesEnTramo(2)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestAsientoTramo_IsAsientoDisponible_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewAsientoTramoService(repo)
	_, _ = svc.Create(input.CreateAsientoTramoInput{IDAsiento: 5, IDTramo: 2, Estado: util.ESTADO_DISPONIBLE})

	disponible, err := svc.IsAsientoDisponible(5, 2)
	assert.NoError(t, err)
	assert.True(t, disponible)
}

func TestAsientoTramo_IsAsientoDisponible_Occupied(t *testing.T) {
	repo := newMockRepo()
	svc := NewAsientoTramoService(repo)
	_, _ = svc.Create(input.CreateAsientoTramoInput{IDAsiento: 5, IDTramo: 2, Estado: util.ESTADO_OCUPADO})

	disponible, err := svc.IsAsientoDisponible(5, 2)
	assert.NoError(t, err)
	assert.False(t, disponible)
}

func TestAsientoTramo_DeleteByVenta_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewAsientoTramoService(repo)
	idVenta := int64(100)
	_, _ = svc.Create(input.CreateAsientoTramoInput{IDAsiento: 5, IDTramo: 2, Estado: util.ESTADO_OCUPADO})
	_ = svc.MarkAsOccupied(5, 2, &idVenta)

	err := svc.DeleteByVenta(idVenta)
	assert.NoError(t, err)
}
