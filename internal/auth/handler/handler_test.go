package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sistema_venta_pasajes/internal/auth/input"
	"sistema_venta_pasajes/internal/auth/service"
	"sistema_venta_pasajes/internal/auth/util"
	"sistema_venta_pasajes/pkg"
)

// ─── Mock del servicio ────────────────────────────────────────────────────────

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Login(ctx context.Context, in input.LoginInput) (*input.TokenPairOutput, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*input.TokenPairOutput), args.Error(1)
}

func (m *mockAuthService) Refresh(ctx context.Context, in input.RefreshInput) (*input.TokenPairOutput, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*input.TokenPairOutput), args.Error(1)
}

func (m *mockAuthService) Logout(ctx context.Context, in input.RefreshInput) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *mockAuthService) ValidarToken(tokenStr string) (*service.AuthClaims, error) {
	args := m.Called(tokenStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.AuthClaims), args.Error(1)
}

// tokenPair de respuesta de ejemplo
func tokenPairOK() *input.TokenPairOutput {
	return &input.TokenPairOutput{
		AccessToken:  "access.token.ejemplo",
		RefreshToken: "refresh-token-ejemplo",
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}
}

// ─── Login ────────────────────────────────────────────────────────────────────

func TestLoginHandler_Exitoso(t *testing.T) {
	svc := new(mockAuthService)
	h := NewAuthHandler(svc)

	body, _ := json.Marshal(input.LoginInput{Email: "admin@test.com", Password: "pass123"})
	svc.On("Login", mock.Anything, input.LoginInput{Email: "admin@test.com", Password: "pass123"}).
		Return(tokenPairOK(), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "access.token.ejemplo")
	assert.Contains(t, w.Body.String(), util.MSG_LOGIN_OK)
}

func TestLoginHandler_BodyVacio(t *testing.T) {
	h := NewAuthHandler(new(mockAuthService))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte("")))
	w := httptest.NewRecorder()
	h.Login(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginHandler_CredencialesInvalidas(t *testing.T) {
	svc := new(mockAuthService)
	h := NewAuthHandler(svc)

	in := input.LoginInput{Email: "mal@test.com", Password: "wrong"}
	body, _ := json.Marshal(in)
	svc.On("Login", mock.Anything, in).
		Return(nil, pkg.Unauthorized(util.ERR_CODE_CREDENCIALES, util.MSG_CREDENCIALES_INVALIDAS))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── Refresh ─────────────────────────────────────────────────────────────────

func TestRefreshHandler_Exitoso(t *testing.T) {
	svc := new(mockAuthService)
	h := NewAuthHandler(svc)

	in := input.RefreshInput{RefreshToken: "valid-refresh-token"}
	body, _ := json.Marshal(in)
	svc.On("Refresh", mock.Anything, in).Return(tokenPairOK(), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Refresh(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), util.MSG_REFRESH_OK)
}

func TestRefreshHandler_TokenInvalido(t *testing.T) {
	svc := new(mockAuthService)
	h := NewAuthHandler(svc)

	in := input.RefreshInput{RefreshToken: "token-falso"}
	body, _ := json.Marshal(in)
	svc.On("Refresh", mock.Anything, in).
		Return(nil, pkg.Unauthorized(util.ERR_CODE_TOKEN_INVALIDO, util.MSG_REFRESH_TOKEN_INVALIDO))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Refresh(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── Logout ───────────────────────────────────────────────────────────────────

func TestLogoutHandler_Exitoso(t *testing.T) {
	svc := new(mockAuthService)
	h := NewAuthHandler(svc)

	in := input.RefreshInput{RefreshToken: "valid-refresh-token"}
	body, _ := json.Marshal(in)
	svc.On("Logout", mock.Anything, in).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Logout(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), util.MSG_LOGOUT_OK)
}

func TestLogoutHandler_BodyVacio(t *testing.T) {
	h := NewAuthHandler(new(mockAuthService))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewReader([]byte("")))
	w := httptest.NewRecorder()
	h.Logout(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
