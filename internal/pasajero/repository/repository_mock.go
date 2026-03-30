package repository

import (
	"sistema_venta_pasajes/internal/pasajero/domain"

	"github.com/stretchr/testify/mock"
)

type PasajeroRepositoryMock struct {
	mock.Mock
}

func (m *PasajeroRepositoryMock) Create(pasajero *domain.Pasajero) error {
	args := m.Called(pasajero)
	return args.Error(0)
}

func (m *PasajeroRepositoryMock) GetByID(id int) (*domain.Pasajero, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pasajero), args.Error(1)
}

func (m *PasajeroRepositoryMock) Update(pasajero *domain.Pasajero) error {
	args := m.Called(pasajero)
	return args.Error(0)
}

func (m *PasajeroRepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *PasajeroRepositoryMock) List(page, pageSize int) ([]domain.Pasajero, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]domain.Pasajero), args.Int(1), args.Error(2)
}
