package service

import (
	"testing"

	vehiculodomain "sistema_venta_pasajes/internal/vehiculo/domain"
	vehiculoinput "sistema_venta_pasajes/internal/vehiculo/input"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVehiculoRepository struct {
	mock.Mock
}

func (m *MockVehiculoRepository) Create(vehiculo *vehiculodomain.Vehiculo) error {
	args := m.Called(vehiculo)
	return args.Error(0)
}
func (m *MockVehiculoRepository) Update(vehiculo *vehiculodomain.Vehiculo) error {
	args := m.Called(vehiculo)
	return args.Error(0)
}
func (m *MockVehiculoRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockVehiculoRepository) GetByID(id int64) (*vehiculodomain.Vehiculo, error) {
	args := m.Called(id)
	return args.Get(0).(*vehiculodomain.Vehiculo), args.Error(1)
}
func (m *MockVehiculoRepository) List(offset, limit int) ([]vehiculodomain.Vehiculo, int, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]vehiculodomain.Vehiculo), args.Int(1), args.Error(2)
}
func (m *MockVehiculoRepository) ExistsByPlaca(placa string) (bool, error) {
	args := m.Called(placa)
	return args.Bool(0), args.Error(1)
}
func (m *MockVehiculoRepository) ExistsByChasis(chasis string) (bool, error) {
	args := m.Called(chasis)
	return args.Bool(0), args.Error(1)
}
func (m *MockVehiculoRepository) ExistsBySoat(soat string) (bool, error) {
	args := m.Called(soat)
	return args.Bool(0), args.Error(1)
}

func TestVehiculoService_Create(t *testing.T) {
	mockRepo := new(MockVehiculoRepository)
	service := NewVehiculoService(mockRepo)
	vehiculoInput := vehiculoinput.CreateVehiculoInput{NroPlaca: "ABC123", NumeroChasis: "CHS123", NroSoat: "SOAT-001"}
	mockRepo.On("ExistsByPlaca", "ABC123").Return(false, nil)
	mockRepo.On("ExistsByChasis", "CHS123").Return(false, nil)
	mockRepo.On("ExistsBySoat", "SOAT-001").Return(false, nil)
	mockRepo.On("Create", mock.Anything).Return(nil)
	output, err := service.Create(vehiculoInput)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestVehiculoService_Create_DuplicatePlaca(t *testing.T) {
	mockRepo := new(MockVehiculoRepository)
	service := NewVehiculoService(mockRepo)
	vehiculoInput := vehiculoinput.CreateVehiculoInput{NroPlaca: "DUP123", NumeroChasis: "CHS999", NroSoat: "SOAT-002"}
	mockRepo.On("ExistsByPlaca", "DUP123").Return(true, nil)
	output, err := service.Create(vehiculoInput)
	assert.Error(t, err)
	assert.Nil(t, output)
}

func TestVehiculoService_Create_DuplicateChasis(t *testing.T) {
	mockRepo := new(MockVehiculoRepository)
	service := NewVehiculoService(mockRepo)
	vehiculoInput := vehiculoinput.CreateVehiculoInput{NroPlaca: "PLACA-OK", NumeroChasis: "CHS-DUP", NroSoat: "SOAT-003"}
	mockRepo.On("ExistsByPlaca", "PLACA-OK").Return(false, nil)
	mockRepo.On("ExistsByChasis", "CHS-DUP").Return(true, nil)
	output, err := service.Create(vehiculoInput)
	assert.Error(t, err)
	assert.Nil(t, output)
}

func TestVehiculoService_Create_DuplicateSoat(t *testing.T) {
	mockRepo := new(MockVehiculoRepository)
	service := NewVehiculoService(mockRepo)
	vehiculoInput := vehiculoinput.CreateVehiculoInput{NroPlaca: "PLACA-OK2", NumeroChasis: "CHS-OK2", NroSoat: "SOAT-DUP"}
	mockRepo.On("ExistsByPlaca", "PLACA-OK2").Return(false, nil)
	mockRepo.On("ExistsByChasis", "CHS-OK2").Return(false, nil)
	mockRepo.On("ExistsBySoat", "SOAT-DUP").Return(true, nil)
	output, err := service.Create(vehiculoInput)
	assert.Error(t, err)
	assert.Nil(t, output)
}
