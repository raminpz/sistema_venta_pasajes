package service

import (
	"errors"
	"strings"
	"testing"

	"sistema_venta_pasajes/internal/terminal/domain"
	"sistema_venta_pasajes/internal/terminal/input"
	"sistema_venta_pasajes/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Create(t *domain.Terminal) error {
	args := m.Called(t)
	return args.Error(0)
}
func (m *mockRepo) GetByID(id int64) (*domain.Terminal, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Terminal), args.Error(1)
}
func (m *mockRepo) Update(t *domain.Terminal) error {
	args := m.Called(t)
	return args.Error(0)
}
func (m *mockRepo) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockRepo) List() ([]domain.Terminal, error) {
	args := m.Called()
	return args.Get(0).([]domain.Terminal), args.Error(1)
}

func TestCreateTerminal(t *testing.T) {
	repo := new(mockRepo)
	svc := NewTerminalService(repo)
	input := input.CreateTerminalInput{
		Nombre:       "Terminal 1",
		Ciudad:       "Ciudad 1",
		Departamento: "Depto 1",
		Direccion:    "Dir 1",
		Estado:       "Activo",
	}
	repo.On("Create", mock.Anything).Return(nil)
	terminal, err := svc.Create(input)
	assert.NoError(t, err)
	assert.Equal(t, "Terminal 1", terminal.NOMBRE)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := new(mockRepo)
	svc := NewTerminalService(repo)
	repo.On("GetByID", int64(1)).Return(&domain.Terminal{}, errors.New("not found"))
	_, err := svc.GetByID(1)
	assert.Error(t, err)
}

func TestDeleteTerminal_OK(t *testing.T) {
	repo := new(mockRepo)
	svc := NewTerminalService(repo)
	terminal := &domain.Terminal{IDTerminal: 1}
	repo.On("GetByID", int64(1)).Return(terminal, nil)
	repo.On("Delete", int64(1)).Return(nil)
	err := svc.Delete(1)
	assert.NoError(t, err)
}

func TestDeleteTerminal_DeleteError(t *testing.T) {
	repo := new(mockRepo)
	svc := NewTerminalService(repo)
	terminal := &domain.Terminal{IDTerminal: 2}
	repo.On("GetByID", int64(2)).Return(terminal, nil)
	repo.On("Delete", int64(2)).Return(errors.New("fk constraint"))
	err := svc.Delete(2)
	assert.Error(t, err)
}

func TestCreateTerminal_ValidationAndDuplicate(t *testing.T) {
	repo := new(mockRepo)
	svc := NewTerminalService(repo)

	_, err := svc.Create(input.CreateTerminalInput{})
	assert.Error(t, err)

	dup := input.CreateTerminalInput{Nombre: "A", Ciudad: "B", Departamento: "C", Direccion: "D", Estado: "ACTIVO"}
	repo.On("Create", mock.Anything).Return(pkg.Conflict("duplicate_resource", "dup")).Once()
	_, err = svc.Create(dup)
	assert.Error(t, err)
}

func TestTerminalService_UpdateAndList(t *testing.T) {
	repo := new(mockRepo)
	svc := NewTerminalService(repo)

	base := &domain.Terminal{IDTerminal: 1, NOMBRE: "Old", CIUDAD: "Ay", DEPARTAMENTO: "Ay", DIRECCION: "Dir", ESTADO: "ACTIVO"}
	repo.On("GetByID", int64(1)).Return(base, nil).Once()
	repo.On("Update", mock.Anything).Return(nil).Once()
	out, err := svc.Update(1, input.UpdateTerminalInput{Nombre: "Nuevo", Ciudad: "Ayacucho", Departamento: "Ayacucho", Direccion: "Dir", Estado: "ACTIVO"})
	assert.NoError(t, err)
	assert.Equal(t, "Nuevo", out.NOMBRE)

	repo.On("GetByID", int64(2)).Return(&domain.Terminal{}, errors.New("db")).Once()
	_, err = svc.Update(2, input.UpdateTerminalInput{Nombre: "X"})
	assert.Error(t, err)

	repo.On("GetByID", int64(3)).Return(&domain.Terminal{}, nil).Once()
	_, err = svc.Update(3, input.UpdateTerminalInput{Nombre: strings.Repeat("a", 101)})
	assert.Error(t, err)

	repo.On("List").Return([]domain.Terminal{{IDTerminal: 1}}, nil).Once()
	list, err := svc.List()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestDeleteTerminal_NotFound(t *testing.T) {
	repo := new(mockRepo)
	svc := NewTerminalService(repo)
	repo.On("GetByID", int64(9)).Return(&domain.Terminal{}, gorm.ErrRecordNotFound).Once()
	err := svc.Delete(9)
	assert.Error(t, err)
}
