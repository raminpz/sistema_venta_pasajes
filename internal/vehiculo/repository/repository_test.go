package repository

import (
	"errors"
	"sistema_venta_pasajes/internal/vehiculo/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock del repositorio
type mockVehiculoRepository struct {
	CreateFunc func(vehiculo *domain.Vehiculo) error
}

func (m *mockVehiculoRepository) Create(vehiculo *domain.Vehiculo) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(vehiculo)
	}
	return nil
}

func TestVehiculoRepository_Create_Mock(t *testing.T) {
	vehiculo := &domain.Vehiculo{
		IDTipoVehiculo: 1,
		NroPlaca:       "MOCK-123",
		Marca:          "MarcaTest",
	}
	mockRepo := &mockVehiculoRepository{
		CreateFunc: func(v *domain.Vehiculo) error {
			if v.NroPlaca == "MOCK-123" {
				return nil
			}
			return errors.New("placa incorrecta")
		},
	}
	err := mockRepo.Create(vehiculo)
	assert.NoError(t, err)
}
