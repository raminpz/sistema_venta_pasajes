package service

import (
	"errors"
	"testing"

	"sistema_venta_pasajes/internal/terminal/domain"
	"sistema_venta_pasajes/internal/terminal/input"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

