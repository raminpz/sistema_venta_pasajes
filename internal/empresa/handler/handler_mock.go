package handler

import (
	"sistema_venta_pasajes/internal/empresa/input"

	"github.com/stretchr/testify/mock"
)

type EmpresaServiceMock struct {
	mock.Mock
}

func (m *EmpresaServiceMock) Create(in input.CreateEmpresaInput) (input.EmpresaOutput, error) {
	args := m.Called(in)
	return args.Get(0).(input.EmpresaOutput), args.Error(1)
}

func (m *EmpresaServiceMock) Update(id int64, in input.UpdateEmpresaInput) (input.EmpresaOutput, error) {
	args := m.Called(id, in)
	return args.Get(0).(input.EmpresaOutput), args.Error(1)
}

func (m *EmpresaServiceMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *EmpresaServiceMock) GetByID(id int64) (input.EmpresaOutput, error) {
	args := m.Called(id)
	return args.Get(0).(input.EmpresaOutput), args.Error(1)
}

func (m *EmpresaServiceMock) List() ([]input.EmpresaOutput, error) {
	args := m.Called()
	return args.Get(0).([]input.EmpresaOutput), args.Error(1)
}
