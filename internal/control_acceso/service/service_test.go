package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"sistema_venta_pasajes/internal/control_acceso/domain"
	"sistema_venta_pasajes/internal/control_acceso/input"
	"sistema_venta_pasajes/internal/control_acceso/service"
)

// mockRepo simula el repositorio
type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetLatest() (*domain.ControlAcceso, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ControlAcceso), args.Error(1)
}

func (m *mockRepo) GetByID(id int64) (*domain.ControlAcceso, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ControlAcceso), args.Error(1)
}

func (m *mockRepo) Create(acceso *domain.ControlAcceso) error {
	args := m.Called(acceso)
	return args.Error(0)
}

func (m *mockRepo) SetEstado(id int64, estado string) error {
	args := m.Called(id, estado)
	return args.Error(0)
}

func (m *mockRepo) Renovar(id int64, nuevaFecha time.Time) error {
	args := m.Called(id, nuevaFecha)
	return args.Error(0)
}

// ---------- GetStatus ----------

func TestGetStatus_Operativo_Normal(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	expiracion := time.Now().AddDate(0, 6, 0)
	repo.On("GetLatest").Return(&domain.ControlAcceso{
		IDAcceso:        1,
		Estado:          "OPERATIVO",
		FechaExpiracion: expiracion,
	}, nil)

	out, err := svc.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, "OPERATIVO", out.EstadoEfectivo)
	assert.False(t, out.EnAlerta)
	assert.False(t, out.EnGracia)
	repo.AssertExpectations(t)
}

func TestGetStatus_Operativo_ConAlerta_FaltanMenos30Dias(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	expiracion := time.Now().AddDate(0, 0, 15)
	repo.On("GetLatest").Return(&domain.ControlAcceso{
		IDAcceso:        1,
		Estado:          "OPERATIVO",
		FechaExpiracion: expiracion,
	}, nil)

	out, err := svc.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, "OPERATIVO", out.EstadoEfectivo)
	assert.True(t, out.EnAlerta, "debe mostrar alerta cuando faltan ≤30 días")
	assert.False(t, out.EnGracia)
	assert.Contains(t, out.Mensaje, "961501468")
	repo.AssertExpectations(t)
}

func TestGetStatus_SoloLectura_EnGracia(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	expiracion := time.Now().AddDate(0, 0, -5)
	repo.On("GetLatest").Return(&domain.ControlAcceso{
		IDAcceso:        1,
		Estado:          "OPERATIVO",
		FechaExpiracion: expiracion,
	}, nil)

	out, err := svc.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, "SOLO_LECTURA", out.EstadoEfectivo)
	assert.True(t, out.EnGracia)
	assert.False(t, out.EnAlerta)
	assert.Contains(t, out.Mensaje, "961501468")
	repo.AssertExpectations(t)
}

func TestGetStatus_Bloqueado_GraciaAgotada(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	expiracion := time.Now().AddDate(0, 0, -40)
	repo.On("GetLatest").Return(&domain.ControlAcceso{
		IDAcceso:        1,
		Estado:          "OPERATIVO",
		FechaExpiracion: expiracion,
	}, nil)

	out, err := svc.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, "BLOQUEADO", out.EstadoEfectivo)
	assert.False(t, out.EnGracia)
	assert.False(t, out.EnAlerta)
	repo.AssertExpectations(t)
}

func TestGetStatus_Bloqueado_ManualOverride(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	expiracion := time.Now().AddDate(1, 0, 0)
	repo.On("GetLatest").Return(&domain.ControlAcceso{
		IDAcceso:        1,
		Estado:          "BLOQUEADO",
		FechaExpiracion: expiracion,
	}, nil)

	out, err := svc.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, "BLOQUEADO", out.EstadoEfectivo)
	assert.False(t, out.EnAlerta)
	assert.False(t, out.EnGracia)
	assert.Contains(t, out.Mensaje, "961501468")
	repo.AssertExpectations(t)
}

func TestGetStatus_SinRegistro(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	repo.On("GetLatest").Return(nil, gorm.ErrRecordNotFound)

	out, err := svc.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, "BLOQUEADO", out.EstadoEfectivo)
	assert.Contains(t, out.Mensaje, "961501468")
	repo.AssertExpectations(t)
}

// ---------- Create ----------

func TestCreate_ConFechasValidas(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	repo.On("Create", mock.AnythingOfType("*domain.ControlAcceso")).Return(nil)

	in := input.ActivarControlAccesoInput{
		FechaActivacion: "2026-01-01",
		FechaExpiracion: "2027-01-01",
	}
	out, err := svc.Create(in)
	assert.NoError(t, err)
	assert.Equal(t, "OPERATIVO", out.EstadoDB)
	assert.Equal(t, "2026-01-01", out.FechaActivacion)
	assert.Equal(t, "2027-01-01", out.FechaExpiracion)
	repo.AssertExpectations(t)
}

func TestCreate_FechasVacias(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	_, err := svc.Create(input.ActivarControlAccesoInput{})
	assert.Error(t, err)
}

func TestCreate_FormatoFechaInvalido(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	_, err := svc.Create(input.ActivarControlAccesoInput{
		FechaActivacion: "01/01/2026",
		FechaExpiracion: "01/01/2027",
	})
	assert.Error(t, err)
}

func TestCreate_FechaExpiracionAnteriorAActivacion(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	_, err := svc.Create(input.ActivarControlAccesoInput{
		FechaActivacion: "2026-06-01",
		FechaExpiracion: "2026-01-01",
	})
	assert.Error(t, err)
}

// ---------- Activar / Bloquear ----------

func TestActivar_OK(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	repo.On("SetEstado", int64(1), "OPERATIVO").Return(nil)

	err := svc.Activar(1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBloquear_OK(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	repo.On("SetEstado", int64(1), "BLOQUEADO").Return(nil)

	err := svc.Bloquear(1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestActivar_NotFound(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	repo.On("SetEstado", int64(99), "OPERATIVO").Return(gorm.ErrRecordNotFound)

	err := svc.Activar(99)
	assert.Error(t, err)
}

// ---------- Renovar ----------

func TestRenovar_OK(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	nuevaFecha := time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC)
	repo.On("Renovar", int64(1), nuevaFecha).Return(nil)
	repo.On("GetByID", int64(1)).Return(&domain.ControlAcceso{
		IDAcceso:        1,
		Estado:          "OPERATIVO",
		FechaExpiracion: nuevaFecha,
	}, nil)

	out, err := svc.Renovar(1, input.RenovarControlAccesoInput{FechaExpiracion: "2028-01-01"})
	assert.NoError(t, err)
	assert.Equal(t, "OPERATIVO", out.EstadoEfectivo)
	assert.Equal(t, "2028-01-01", out.FechaExpiracion)
	repo.AssertExpectations(t)
}

func TestRenovar_FechaVacia(t *testing.T) {
	repo := new(mockRepo)
	svc := service.New(repo)

	_, err := svc.Renovar(1, input.RenovarControlAccesoInput{})
	assert.Error(t, err)
}
