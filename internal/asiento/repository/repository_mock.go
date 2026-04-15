package repository

import (
	"sistema_venta_pasajes/internal/asiento/domain"

	"github.com/stretchr/testify/mock"
)

type AsientoRepositoryMock struct {
	mock.Mock
}

func (m *AsientoRepositoryMock) Create(asiento *domain.Asiento) error {
	args := m.Called(asiento)
	return args.Error(0)
}

func (m *AsientoRepositoryMock) GetByID(id int64) (*domain.Asiento, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Asiento), args.Error(1)
}

func (m *AsientoRepositoryMock) ListByVehiculo(idVehiculo int64) ([]*domain.Asiento, error) {
	args := m.Called(idVehiculo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Asiento), args.Error(1)
}

func (m *AsientoRepositoryMock) Update(asiento *domain.Asiento) error {
	args := m.Called(asiento)
	return args.Error(0)
}

func (m *AsientoRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *AsientoRepositoryMock) CambiarEstado(id int64, estado string) error {
	args := m.Called(id, estado)
	return args.Error(0)
}

