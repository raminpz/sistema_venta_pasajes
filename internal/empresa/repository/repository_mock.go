package repository

import (
	"sistema_venta_pasajes/internal/empresa/domain"

	"github.com/stretchr/testify/mock"
)

type EmpresaRepositoryMock struct {
	mock.Mock
}

func (m *EmpresaRepositoryMock) Create(empresa *domain.Empresa) error {
	args := m.Called(empresa)
	return args.Error(0)
}

func (m *EmpresaRepositoryMock) GetByID(id int64) (*domain.Empresa, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Empresa), args.Error(1)
}

func (m *EmpresaRepositoryMock) Update(empresa *domain.Empresa) error {
	args := m.Called(empresa)
	return args.Error(0)
}

func (m *EmpresaRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *EmpresaRepositoryMock) List() ([]domain.Empresa, error) {
	args := m.Called()
	return args.Get(0).([]domain.Empresa), args.Error(1)
}
