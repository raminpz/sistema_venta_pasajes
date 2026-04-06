package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sistema_venta_pasajes/internal/auth/domain"
	"sistema_venta_pasajes/internal/auth/input"
	"sistema_venta_pasajes/internal/auth/repository"
	"sistema_venta_pasajes/internal/auth/util"
)

const testSecret = "secreto_jwt_de_prueba_256bits_para_tests"

// ─── helpers ────────────────────────────────────────────────────────────────

func hashTok(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func activeUser() *repository.UserAuthData {
	return &repository.UserAuthData{
		IDUsuario: 1,
		Email:     "admin@test.com",
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Estado:    "ACTIVO",
		RolNombre: "ADMIN",
	}
}

// ─── Login ───────────────────────────────────────────────────────────────────

func TestLogin_Exitoso(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo, testSecret)

	mockRepo.On("FindUserForAuth", mock.Anything, "admin@test.com").Return(activeUser(), nil)
	mockRepo.On("RevokeAllUserTokens", mock.Anything, 1).Return(nil)
	mockRepo.On("SaveRefreshToken", mock.Anything, mock.Anything).Return(nil)

	out, err := svc.Login(context.Background(), input.LoginInput{
		Email:    "admin@test.com",
		Password: "password",
	})

	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotEmpty(t, out.AccessToken)
	assert.NotEmpty(t, out.RefreshToken)
	assert.Equal(t, "Bearer", out.TokenType)
	assert.Equal(t, 900, out.ExpiresIn) // 15 min
	mockRepo.AssertExpectations(t)
}

func TestLogin_EmailVacio(t *testing.T) {
	svc := NewAuthService(new(repository.MockAuthRepository), testSecret)
	out, err := svc.Login(context.Background(), input.LoginInput{Email: "", Password: "pass"})
	assert.Nil(t, out)
	assert.EqualError(t, err, util.MSG_EMAIL_REQUERIDO)
}

func TestLogin_PasswordVacio(t *testing.T) {
	svc := NewAuthService(new(repository.MockAuthRepository), testSecret)
	out, err := svc.Login(context.Background(), input.LoginInput{Email: "x@x.com", Password: ""})
	assert.Nil(t, out)
	assert.EqualError(t, err, util.MSG_PASSWORD_REQUERIDO)
}

func TestLogin_UsuarioNoExiste(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo, testSecret)
	mockRepo.On("FindUserForAuth", mock.Anything, "no@existe.com").Return(nil, nil)

	out, err := svc.Login(context.Background(), input.LoginInput{Email: "no@existe.com", Password: "pass"})
	assert.Nil(t, out)
	assert.EqualError(t, err, util.MSG_CREDENCIALES_INVALIDAS)
}

func TestLogin_PasswordIncorrecto(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo, testSecret)
	mockRepo.On("FindUserForAuth", mock.Anything, "admin@test.com").Return(activeUser(), nil)

	out, err := svc.Login(context.Background(), input.LoginInput{Email: "admin@test.com", Password: "wrong"})
	assert.Nil(t, out)
	assert.EqualError(t, err, util.MSG_CREDENCIALES_INVALIDAS)
}

func TestLogin_UsuarioInactivo(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo, testSecret)
	inactivo := activeUser()
	inactivo.Estado = "INACTIVO"
	mockRepo.On("FindUserForAuth", mock.Anything, "admin@test.com").Return(inactivo, nil)

	out, err := svc.Login(context.Background(), input.LoginInput{Email: "admin@test.com", Password: "password"})
	assert.Nil(t, out)
	assert.EqualError(t, err, util.MSG_USUARIO_INACTIVO)
}

// ─── Refresh ──────────────────────────────────────────────────────────────────

func TestRefresh_Exitoso(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo, testSecret)

	rawToken := "mi-refresh-token-valido"
	hash := hashTok(rawToken)

	rt := &domain.RefreshToken{
		ID:        1,
		IDUsuario: 1,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Revocado:  false,
	}

	mockRepo.On("GetRefreshToken", mock.Anything, hash).Return(rt, nil)
	mockRepo.On("RevokeRefreshToken", mock.Anything, hash).Return(nil)
	mockRepo.On("FindUserByID", mock.Anything, 1).Return(activeUser(), nil)
	mockRepo.On("RevokeAllUserTokens", mock.Anything, 1).Return(nil)
	mockRepo.On("SaveRefreshToken", mock.Anything, mock.Anything).Return(nil)

	out, err := svc.Refresh(context.Background(), input.RefreshInput{RefreshToken: rawToken})
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotEmpty(t, out.AccessToken)
}

func TestRefresh_TokenVacio(t *testing.T) {
	svc := NewAuthService(new(repository.MockAuthRepository), testSecret)
	out, err := svc.Refresh(context.Background(), input.RefreshInput{RefreshToken: ""})
	assert.Nil(t, out)
	assert.EqualError(t, err, util.MSG_REFRESH_TOKEN_REQUERIDO)
}

func TestRefresh_TokenNoExiste(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo, testSecret)
	mockRepo.On("GetRefreshToken", mock.Anything, mock.Anything).Return(nil, nil)

	out, err := svc.Refresh(context.Background(), input.RefreshInput{RefreshToken: "token-falso"})
	assert.Nil(t, out)
	assert.EqualError(t, err, util.MSG_REFRESH_TOKEN_INVALIDO)
}

func TestRefresh_TokenRevocado(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo, testSecret)

	rawToken := "token-revocado"
	hash := hashTok(rawToken)
	rt := &domain.RefreshToken{
		ID:        2,
		IDUsuario: 1,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Revocado:  true,
	}
	mockRepo.On("GetRefreshToken", mock.Anything, hash).Return(rt, nil)

	out, err := svc.Refresh(context.Background(), input.RefreshInput{RefreshToken: rawToken})
	assert.Nil(t, out)
	assert.EqualError(t, err, util.MSG_REFRESH_TOKEN_INVALIDO)
}

func TestRefresh_TokenExpirado(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo, testSecret)

	rawToken := "token-expirado"
	hash := hashTok(rawToken)
	rt := &domain.RefreshToken{
		ID:        3,
		IDUsuario: 1,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // ya expiró
		Revocado:  false,
	}
	mockRepo.On("GetRefreshToken", mock.Anything, hash).Return(rt, nil)

	out, err := svc.Refresh(context.Background(), input.RefreshInput{RefreshToken: rawToken})
	assert.Nil(t, out)
	assert.EqualError(t, err, util.MSG_REFRESH_TOKEN_INVALIDO)
}

// ─── Logout ───────────────────────────────────────────────────────────────────

func TestLogout_Exitoso(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	svc := NewAuthService(mockRepo, testSecret)

	rawToken := "mi-refresh-token-para-logout"
	hash := hashTok(rawToken)
	mockRepo.On("RevokeRefreshToken", mock.Anything, hash).Return(nil)

	err := svc.Logout(context.Background(), input.RefreshInput{RefreshToken: rawToken})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLogout_TokenVacio(t *testing.T) {
	svc := NewAuthService(new(repository.MockAuthRepository), testSecret)
	err := svc.Logout(context.Background(), input.RefreshInput{RefreshToken: ""})
	assert.EqualError(t, err, util.MSG_REFRESH_TOKEN_REQUERIDO)
}

// ─── ValidarToken ─────────────────────────────────────────────────────────────

func TestValidarToken_Valido(t *testing.T) {
	mockRepo := new(repository.MockAuthRepository)
	mockRepo.On("RevokeAllUserTokens", mock.Anything, 1).Return(nil)
	mockRepo.On("SaveRefreshToken", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("FindUserForAuth", mock.Anything, "admin@test.com").Return(activeUser(), nil)

	svc := NewAuthService(mockRepo, testSecret)

	// Generar un token real haciendo login
	out, err := svc.Login(context.Background(), input.LoginInput{
		Email:    "admin@test.com",
		Password: "password",
	})
	assert.NoError(t, err)

	claims, err := svc.ValidarToken(out.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, 1, claims.IDUsuario)
	assert.Equal(t, "admin@test.com", claims.Email)
	assert.Equal(t, "ADMIN", claims.Rol)
}

func TestValidarToken_TokenInvalido(t *testing.T) {
	svc := NewAuthService(new(repository.MockAuthRepository), testSecret)
	claims, err := svc.ValidarToken("token.invalido.firma")
	assert.Nil(t, claims)
	assert.EqualError(t, err, util.MSG_TOKEN_INVALIDO)
}

func TestValidarToken_SecretoDistinto(t *testing.T) {
	// Token firmado con secreto A, validado con secreto B → debe fallar
	svcA := NewAuthService(new(repository.MockAuthRepository), "secreto-A")
	svcB := NewAuthService(new(repository.MockAuthRepository), "secreto-B")

	mockRepo := new(repository.MockAuthRepository)
	mockRepo.On("FindUserForAuth", mock.Anything, "admin@test.com").Return(activeUser(), nil)
	mockRepo.On("RevokeAllUserTokens", mock.Anything, 1).Return(nil)
	mockRepo.On("SaveRefreshToken", mock.Anything, mock.Anything).Return(nil)
	svcA2 := NewAuthService(mockRepo, "secreto-A")

	out, _ := svcA2.Login(context.Background(), input.LoginInput{
		Email: "admin@test.com", Password: "password",
	})
	_ = svcA

	claims, err := svcB.ValidarToken(out.AccessToken)
	assert.Nil(t, claims)
	assert.Error(t, err)
}
