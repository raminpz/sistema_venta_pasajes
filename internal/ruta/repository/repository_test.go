package repository

import (
	"testing"
	"sistema_venta_pasajes/internal/ruta/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetByID(t *testing.T) {
	mockRepo := &RutaRepositoryMock{}
	ruta := &domain.Ruta{IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5}
	mockRepo.On("Create", ruta).Return(nil)
	mockRepo.On("GetByID", ruta.IDRuta).Return(ruta, nil)

	err := mockRepo.Create(ruta)
	assert.NoError(t, err)
	got, err := mockRepo.GetByID(ruta.IDRuta)
	assert.NoError(t, err)
	assert.Equal(t, ruta, got)
}

func TestUpdate(t *testing.T) {
	mockRepo := &RutaRepositoryMock{}
	ruta := &domain.Ruta{IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5}
	mockRepo.On("Update", ruta).Return(nil)
	err := mockRepo.Update(ruta)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	mockRepo := &RutaRepositoryMock{}
	rutaID := 1
	mockRepo.On("Delete", rutaID).Return(nil)
	err := mockRepo.Delete(rutaID)
	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	mockRepo := &RutaRepositoryMock{}
	rutas := []domain.Ruta{
		{IDRuta: 1, IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5},
		{IDRuta: 2, IDOrigenTerminal: 2, IDDestinoTerminal: 3, DuracionHoras: 8.0},
	}
	mockRepo.On("List").Return(rutas, nil)
	result, err := mockRepo.List()
	assert.NoError(t, err)
	assert.Equal(t, rutas, result)
}
