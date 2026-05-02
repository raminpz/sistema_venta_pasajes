package repository

import (
	"sistema_venta_pasajes/internal/tramo/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTramoRepository struct {
	CreateFunc              func(*domain.Tramo) error
	ExistsByRutaParadasFunc func(int64, int64, int64) (bool, error)
}

func (m *mockTramoRepository) Create(t *domain.Tramo) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(t)
	}
	return nil
}
func (m *mockTramoRepository) Update(t *domain.Tramo) error            { return nil }
func (m *mockTramoRepository) Delete(id int64) error                   { return nil }
func (m *mockTramoRepository) GetByID(id int64) (*domain.Tramo, error) { return nil, nil }
func (m *mockTramoRepository) List(offset, limit int) ([]domain.Tramo, int, error) {
	return []domain.Tramo{}, 0, nil
}
func (m *mockTramoRepository) ListByRuta(idRuta int64) ([]domain.Tramo, error) {
	return []domain.Tramo{}, nil
}
func (m *mockTramoRepository) ExistsByRutaParadas(idRuta, idOrigen, idDestino int64) (bool, error) {
	if m.ExistsByRutaParadasFunc != nil {
		return m.ExistsByRutaParadasFunc(idRuta, idOrigen, idDestino)
	}
	return false, nil
}

func TestTramoRepository_Create_OK(t *testing.T) {
	tramo := &domain.Tramo{IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 2}
	mock := &mockTramoRepository{
		CreateFunc: func(_ *domain.Tramo) error { return nil },
	}
	err := mock.Create(tramo)
	assert.NoError(t, err)
}

func TestTramoRepository_Create_DatosMinimos(t *testing.T) {
	tramo := &domain.Tramo{IDRuta: 7, IDParadaOrigen: 10, IDParadaDestino: 11}
	mock := &mockTramoRepository{
		CreateFunc: func(tr *domain.Tramo) error {
			assert.Equal(t, int64(7), tr.IDRuta)
			assert.Equal(t, int64(10), tr.IDParadaOrigen)
			assert.Equal(t, int64(11), tr.IDParadaDestino)
			return nil
		},
	}
	err := mock.Create(tramo)
	assert.NoError(t, err)
}

func TestTramoRepository_ExistsRutaParadas_True(t *testing.T) {
	mock := &mockTramoRepository{
		ExistsByRutaParadasFunc: func(idRuta, idOrigen, idDestino int64) (bool, error) {
			return true, nil
		},
	}
	exists, err := mock.ExistsByRutaParadas(1, 1, 2)
	assert.NoError(t, err)
	assert.True(t, exists)
}
