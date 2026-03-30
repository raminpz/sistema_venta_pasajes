package repository

import (
	"testing"
	"sistema_venta_pasajes/internal/terminal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetByID(t *testing.T) {
	mockRepo := &TerminalRepositoryMock{}
	terminal := &domain.Terminal{
		IDTerminal:   1,
		NOMBRE:       "Terminal Test",
		CIUDAD:       "Ciudad Test",
		DEPARTAMENTO: "Depto Test",
		DIRECCION:    "Dir Test",
		ESTADO:       "Activo",
	}
	mockRepo.On("Create", terminal).Return(nil)
	mockRepo.On("GetByID", int(terminal.IDTerminal)).Return(terminal, nil)

	err := mockRepo.Create(terminal)
	assert.NoError(t, err)
	fetched, err := mockRepo.GetByID(int(terminal.IDTerminal))
	assert.NoError(t, err)
	assert.Equal(t, "Terminal Test", fetched.NOMBRE)
}
