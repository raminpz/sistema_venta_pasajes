package repository

import (
	"sistema_venta_pasajes/internal/pasajero/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePasajero_Ok(t *testing.T) {
	mockRepo := &PasajeroRepositoryMock{}
	pasajero := &domain.Pasajero{Nombres: "Juan"}
	mockRepo.On("Create", pasajero).Return(nil)
	err := mockRepo.Create(pasajero)
	assert.NoError(t, err)
}

func TestGetByID_NotFound(t *testing.T) {
	mockRepo := &PasajeroRepositoryMock{}
	mockRepo.On("GetByID", 999).Return(nil, assert.AnError)
	_, err := mockRepo.GetByID(999)
	assert.Error(t, err)
}

func TestListPasajeros_Ok(t *testing.T) {
	mockRepo := &PasajeroRepositoryMock{}
	pasajeros := []domain.Pasajero{{Nombres: "Juan"}, {Nombres: "Ana"}}
	mockRepo.On("List", 1, 10).Return(pasajeros, 2, nil)
	res, total, err := mockRepo.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, 2, total)
}
