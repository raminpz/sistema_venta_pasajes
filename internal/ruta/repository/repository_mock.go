package repository

import (
	"sistema_venta_pasajes/internal/ruta/domain"

	"github.com/stretchr/testify/mock"
)

type RutaRepositoryMock struct {
	mock.Mock
}

func (m *RutaRepositoryMock) Create(ruta *domain.Ruta) error {
	args := m.Called(ruta)
	return args.Error(0)
}

func (m *RutaRepositoryMock) GetByID(id int) (*domain.Ruta, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ruta), args.Error(1)
}

func (m *RutaRepositoryMock) Update(ruta *domain.Ruta) error {
	args := m.Called(ruta)
	return args.Error(0)
}

func (m *RutaRepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RutaRepositoryMock) List() ([]domain.Ruta, error) {
	args := m.Called()
	return args.Get(0).([]domain.Ruta), args.Error(1)
}
