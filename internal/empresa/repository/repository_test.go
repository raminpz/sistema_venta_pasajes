package repository

import (
	"sistema_venta_pasajes/internal/empresa/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpresaRepository_CRUD(t *testing.T) {
	mockRepo := &EmpresaRepositoryMock{}
	// Crear
	emp := &domain.Empresa{
		IDEmpresa:   1,
		RUC:         "12345678901",
		RazonSocial: "Empresa Test",
		Telefono:    "987654321",
	}
	mockRepo.On("Create", emp).Return(nil).Once()
	mockRepo.On("GetByID", int64(emp.IDEmpresa)).Return(emp, nil).Once()
	mockRepo.On("Update", emp).Return(nil).Once()
	mockRepo.On("List").Return([]domain.Empresa{*emp}, nil).Once()
	mockRepo.On("Delete", int64(emp.IDEmpresa)).Return(nil).Once()
	mockRepo.On("GetByID", int64(emp.IDEmpresa)).Return(nil, assert.AnError).Once()

	err := mockRepo.Create(emp)
	assert.NoError(t, err)
	got, err := mockRepo.GetByID(int64(emp.IDEmpresa))
	assert.NoError(t, err)
	assert.Equal(t, emp.RUC, got.RUC)

	emp.RazonSocial = "Empresa Actualizada"
	err = mockRepo.Update(emp)
	assert.NoError(t, err)

	list, err := mockRepo.List()
	assert.NoError(t, err)
	assert.True(t, len(list) > 0)

	err = mockRepo.Delete(int64(emp.IDEmpresa))
	assert.NoError(t, err)
	_, err = mockRepo.GetByID(int64(emp.IDEmpresa))
	assert.Error(t, err)
}
