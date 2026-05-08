package service

import (
	"errors"
	"testing"
	"time"

	vehiculodomain "sistema_venta_pasajes/internal/vehiculo/domain"
	vehiculoinput "sistema_venta_pasajes/internal/vehiculo/input"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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

func dateOnlyPtr(t time.Time) *vehiculoinput.DateOnly {
	return &vehiculoinput.DateOnly{Time: t}
}

func TestVehiculoService_Create(t *testing.T) {
	mockRepo := new(MockVehiculoRepository)
	service := NewVehiculoService(mockRepo)
	vehiculoInput := vehiculoinput.CreateVehiculoInput{
		IDTipoVehiculo:       1,
		NroPlaca:             "ABC-123",
		Marca:                "Toyota",
		Modelo:               "Coaster",
		AnioFabricacion:      2022,
		NumeroChasis:         "CHS123456",
		Capacidad:            30,
		NroSoat:              "SOAT-001",
		FechaVencSoat:        dateOnlyPtr(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)),
		NroRevisionTecnica:   "REV-001",
		FechaVencRevisionTec: dateOnlyPtr(time.Date(2027, 11, 30, 0, 0, 0, 0, time.UTC)),
		Estado:               "ACTIVO",
	}
	mockRepo.On("ExistsByPlaca", "ABC-123").Return(false, nil)
	mockRepo.On("Create", mock.Anything).Return(nil)
	output, err := service.Create(vehiculoInput)
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

func TestVehiculoService_Create_MissingFields(t *testing.T) {
	mockRepo := new(MockVehiculoRepository)
	service := NewVehiculoService(mockRepo)

	// Sin IDTipoVehiculo
	_, err := service.Create(vehiculoinput.CreateVehiculoInput{NroPlaca: "ABC-123"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "obligatorio")

	// Sin placa
	_, err = service.Create(vehiculoinput.CreateVehiculoInput{IDTipoVehiculo: 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "obligatorio")

	// Fecha SOAT faltante
	_, err = service.Create(vehiculoinput.CreateVehiculoInput{
		IDTipoVehiculo: 1, NroPlaca: "ABC-123", Marca: "Toyota", Modelo: "Coaster",
		AnioFabricacion: 2022, NumeroChasis: "CHS1", Capacidad: 10,
		NroSoat: "SOAT-1", NroRevisionTecnica: "REV-1",
		FechaVencRevisionTec: dateOnlyPtr(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)),
		Estado:               "ACTIVO",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SOAT")

	// Estado inválido
	_, err = service.Create(vehiculoinput.CreateVehiculoInput{
		IDTipoVehiculo: 1, NroPlaca: "ABC-123", Marca: "Toyota", Modelo: "Coaster",
		AnioFabricacion: 2022, NumeroChasis: "CHS1", Capacidad: 10,
		NroSoat: "SOAT-1", FechaVencSoat: dateOnlyPtr(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)),
		NroRevisionTecnica: "REV-1", FechaVencRevisionTec: dateOnlyPtr(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)),
		Estado: "INVALIDO",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ACTIVO")
}

func TestVehiculoService_Create_DuplicatePlaca(t *testing.T) {
	mockRepo := new(MockVehiculoRepository)
	service := NewVehiculoService(mockRepo)
	vehiculoInput := vehiculoinput.CreateVehiculoInput{
		IDTipoVehiculo:       1,
		NroPlaca:             "DUP123",
		Marca:                "Toyota",
		Modelo:               "Coaster",
		AnioFabricacion:      2022,
		NumeroChasis:         "CHS999",
		Capacidad:            20,
		NroSoat:              "SOAT-002",
		FechaVencSoat:        dateOnlyPtr(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)),
		NroRevisionTecnica:   "REV-002",
		FechaVencRevisionTec: dateOnlyPtr(time.Date(2027, 10, 31, 0, 0, 0, 0, time.UTC)),
		Estado:               "ACTIVO",
	}
	mockRepo.On("ExistsByPlaca", "DUP123").Return(true, nil)
	output, err := service.Create(vehiculoInput)
	assert.Error(t, err)
	assert.Nil(t, output)
}

func TestVehiculoService_Update_Get_Delete_List(t *testing.T) {
	mockRepo := new(MockVehiculoRepository)
	service := NewVehiculoService(mockRepo)

	nombreMarca := "toyota"
	nombreModelo := "coaster"
	nuevaPlaca := "abC-123"
	estado := "activo"
	vehiculo := &vehiculodomain.Vehiculo{IDVehiculo: 1, NroPlaca: "ABC-111", Marca: "Toyota", Modelo: "Old", Estado: "ACTIVO"}

	mockRepo.On("GetByID", int64(1)).Return(vehiculo, nil).Once()
	mockRepo.On("Update", mock.Anything).Return(nil).Once()
	out, err := service.Update(vehiculoinput.UpdateVehiculoInput{IDVehiculo: 1, NroPlaca: &nuevaPlaca, Marca: &nombreMarca, Modelo: &nombreModelo, Estado: &estado})
	assert.NoError(t, err)
	assert.Equal(t, "ABC-123", out.NroPlaca)

	mockRepo.On("GetByID", int64(404)).Return((*vehiculodomain.Vehiculo)(nil), errors.New("not found")).Once()
	_, err = service.GetByID(404)
	assert.Error(t, err)

	mockRepo.On("GetByID", int64(2)).Return(&vehiculodomain.Vehiculo{IDVehiculo: 2}, nil).Once()
	mockRepo.On("Delete", int64(2)).Return(nil).Once()
	err = service.Delete(2)
	assert.NoError(t, err)

	mockRepo.On("GetByID", int64(3)).Return((*vehiculodomain.Vehiculo)(nil), gorm.ErrRecordNotFound).Once()
	err = service.Delete(3)
	assert.Error(t, err)

	mockRepo.On("List", 0, 15).Return([]vehiculodomain.Vehiculo{{IDVehiculo: 1, NroPlaca: "ABC-123", Modelo: "Coaster", Estado: "ACTIVO"}}, 1, nil).Once()
	list, total, err := service.List(1, 15)
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)
}
