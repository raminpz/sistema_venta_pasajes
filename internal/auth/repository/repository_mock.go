package repository

import (
	"context"

	"github.com/stretchr/testify/mock"

	"sistema_venta_pasajes/internal/auth/domain"
)

// MockAuthRepository es el mock del repositorio para pruebas unitarias.
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) FindUserForAuth(ctx context.Context, email string) (*UserAuthData, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserAuthData), args.Error(1)
}

func (m *MockAuthRepository) FindUserByID(ctx context.Context, id int) (*UserAuthData, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserAuthData), args.Error(1)
}

func (m *MockAuthRepository) SaveRefreshToken(ctx context.Context, rt *domain.RefreshToken) error {
	args := m.Called(ctx, rt)
	return args.Error(0)
}

func (m *MockAuthRepository) GetRefreshToken(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockAuthRepository) RevokeRefreshToken(ctx context.Context, hash string) error {
	args := m.Called(ctx, hash)
	return args.Error(0)
}

func (m *MockAuthRepository) RevokeAllUserTokens(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
