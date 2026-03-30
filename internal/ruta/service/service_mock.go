package service

import (
	"context"
	"sistema_venta_pasajes/internal/ruta/domain"
	"sistema_venta_pasajes/internal/ruta/input"

	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) List(ctx context.Context) ([]domain.Ruta, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Ruta), args.Error(1)
}

func (m *ServiceMock) GetByID(ctx context.Context, id int) (*domain.Ruta, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ruta), args.Error(1)
}

func (m *ServiceMock) Create(ctx context.Context, in input.CreateRutaInput) (*domain.Ruta, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ruta), args.Error(1)
}

func (m *ServiceMock) Update(ctx context.Context, id int, in input.UpdateRutaInput) (*domain.Ruta, error) {
	args := m.Called(ctx, id, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ruta), args.Error(1)
}

func (m *ServiceMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
