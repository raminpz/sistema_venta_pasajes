package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sistema_venta_pasajes/internal/auth/domain"
)

func TestMockAuthRepository_FindUserForAuth(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	ctx := context.Background()
	expected := &UserAuthData{IDUsuario: 5, Email: "rami@test.com", RolNombre: "ADMIN", Estado: "ACTIVO"}
	mockRepo.On("FindUserForAuth", ctx, "rami@test.com").Return(expected, nil)

	result, err := mockRepo.FindUserForAuth(ctx, "rami@test.com")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_FindUserForAuth_NoExiste(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	ctx := context.Background()
	mockRepo.On("FindUserForAuth", ctx, "noexiste@test.com").Return(nil, nil)

	result, err := mockRepo.FindUserForAuth(ctx, "noexiste@test.com")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestMockAuthRepository_FindUserByID(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	ctx := context.Background()
	expected := &UserAuthData{IDUsuario: 3, Email: "user@test.com", RolNombre: "VENDEDOR", Estado: "ACTIVO"}
	mockRepo.On("FindUserByID", ctx, 3).Return(expected, nil)

	result, err := mockRepo.FindUserByID(ctx, 3)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestMockAuthRepository_SaveRefreshToken(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	ctx := context.Background()
	rt := &domain.RefreshToken{IDUsuario: 1, TokenHash: "abc123", ExpiresAt: time.Now().Add(24 * time.Hour)}
	mockRepo.On("SaveRefreshToken", ctx, rt).Return(nil)

	err := mockRepo.SaveRefreshToken(ctx, rt)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockAuthRepository_GetRefreshToken_Valido(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	ctx := context.Background()
	rt := &domain.RefreshToken{ID: 1, IDUsuario: 1, TokenHash: "hash-valido", ExpiresAt: time.Now().Add(24 * time.Hour), Revocado: false}
	mockRepo.On("GetRefreshToken", ctx, "hash-valido").Return(rt, nil)

	result, err := mockRepo.GetRefreshToken(ctx, "hash-valido")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Revocado)
}

func TestMockAuthRepository_GetRefreshToken_NoExiste(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	ctx := context.Background()
	mockRepo.On("GetRefreshToken", ctx, "hash-inexistente").Return(nil, nil)

	result, err := mockRepo.GetRefreshToken(ctx, "hash-inexistente")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestMockAuthRepository_RevokeRefreshToken(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	ctx := context.Background()
	mockRepo.On("RevokeRefreshToken", ctx, "hash-a-revocar").Return(nil)

	err := mockRepo.RevokeRefreshToken(ctx, "hash-a-revocar")
	assert.NoError(t, err)
}

func TestMockAuthRepository_RevokeAllUserTokens(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	ctx := context.Background()
	mockRepo.On("RevokeAllUserTokens", ctx, mock.AnythingOfType("int")).Return(nil)

	err := mockRepo.RevokeAllUserTokens(ctx, 7)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
