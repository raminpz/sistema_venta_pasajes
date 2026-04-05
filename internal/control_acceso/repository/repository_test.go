package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sistema_venta_pasajes/internal/control_acceso/domain"
	"sistema_venta_pasajes/internal/control_acceso/repository"
)

// mockControlAccesoRepository implementa ControlAccesoRepository para tests
type mockControlAccesoRepository struct {
	mock.Mock
}

func (m *mockControlAccesoRepository) GetLatest() (*domain.ControlAcceso, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ControlAcceso), args.Error(1)
}

func (m *mockControlAccesoRepository) GetByID(id int64) (*domain.ControlAcceso, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ControlAcceso), args.Error(1)
}

func (m *mockControlAccesoRepository) Create(acceso *domain.ControlAcceso) error {
	args := m.Called(acceso)
	return args.Error(0)
}

func (m *mockControlAccesoRepository) SetEstado(id int64, estado string) error {
	args := m.Called(id, estado)
	return args.Error(0)
}

func (m *mockControlAccesoRepository) Renovar(id int64, nuevaFecha time.Time) error {
	args := m.Called(id, nuevaFecha)
	return args.Error(0)
}

// Asegura que mockControlAccesoRepository implementa la interfaz
var _ repository.ControlAccesoRepository = (*mockControlAccesoRepository)(nil)

func TestMockRepository_GetLatest(t *testing.T) {
	repo := new(mockControlAccesoRepository)
	expected := &domain.ControlAcceso{
		IDAcceso: 1,
		Estado:   "OPERATIVO",
	}
	repo.On("GetLatest").Return(expected, nil)

	result, err := repo.GetLatest()
	assert.NoError(t, err)
	assert.Equal(t, expected.IDAcceso, result.IDAcceso)
	assert.Equal(t, "OPERATIVO", result.Estado)
	repo.AssertExpectations(t)
}

func TestMockRepository_GetByID(t *testing.T) {
	repo := new(mockControlAccesoRepository)
	expected := &domain.ControlAcceso{IDAcceso: 5, Estado: "BLOQUEADO"}
	repo.On("GetByID", int64(5)).Return(expected, nil)

	result, err := repo.GetByID(5)
	assert.NoError(t, err)
	assert.Equal(t, "BLOQUEADO", result.Estado)
	repo.AssertExpectations(t)
}

func TestMockRepository_Create(t *testing.T) {
	repo := new(mockControlAccesoRepository)
	a := &domain.ControlAcceso{Estado: "OPERATIVO"}
	repo.On("Create", a).Return(nil)

	err := repo.Create(a)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMockRepository_SetEstado(t *testing.T) {
	repo := new(mockControlAccesoRepository)
	repo.On("SetEstado", int64(1), "BLOQUEADO").Return(nil)

	err := repo.SetEstado(1, "BLOQUEADO")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMockRepository_Renovar(t *testing.T) {
	repo := new(mockControlAccesoRepository)
	fecha := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	repo.On("Renovar", int64(1), fecha).Return(nil)

	err := repo.Renovar(1, fecha)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
