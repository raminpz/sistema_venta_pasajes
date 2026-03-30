package repository

import (
	"sistema_venta_pasajes/internal/terminal/domain"

	"github.com/stretchr/testify/mock"
)

type TerminalRepositoryMock struct {
	mock.Mock
}

func (m *TerminalRepositoryMock) Create(terminal *domain.Terminal) error {
	args := m.Called(terminal)
	return args.Error(0)
}

func (m *TerminalRepositoryMock) GetByID(id int) (*domain.Terminal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Terminal), args.Error(1)
}

func (m *TerminalRepositoryMock) Update(terminal *domain.Terminal) error {
	args := m.Called(terminal)
	return args.Error(0)
}

func (m *TerminalRepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *TerminalRepositoryMock) List() ([]domain.Terminal, error) {
	args := m.Called()
	return args.Get(0).([]domain.Terminal), args.Error(1)
}
