package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"sistema_venta_pasajes/internal/auth/domain"
	"sistema_venta_pasajes/internal/auth/input"
	"sistema_venta_pasajes/internal/auth/repository"
	"sistema_venta_pasajes/internal/auth/util"
	"sistema_venta_pasajes/pkg"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

// AuthClaims representa el payload del JWT.
type AuthClaims struct {
	IDUsuario int    `json:"id_usuario"`
	Email     string `json:"email"`
	Rol       string `json:"rol"`
	jwt.RegisteredClaims
}

// AuthService define los casos de uso de autenticación.
type AuthService interface {
	Login(ctx context.Context, in input.LoginInput) (*input.TokenPairOutput, error)
	Refresh(ctx context.Context, in input.RefreshInput) (*input.TokenPairOutput, error)
	Logout(ctx context.Context, in input.RefreshInput) error
	ValidarToken(tokenStr string) (*AuthClaims, error)
}

type authService struct {
	repo      repository.AuthRepository
	jwtSecret string
}

// NewAuthService crea el servicio de autenticación.
func NewAuthService(repo repository.AuthRepository, jwtSecret string) AuthService {
	return &authService{repo: repo, jwtSecret: jwtSecret}
}

// Login valida credenciales y devuelve el par de tokens JWT.
func (s *authService) Login(ctx context.Context, in input.LoginInput) (*input.TokenPairOutput, error) {
	msg, ok := util.ValidarLoginInput(in.Email, in.Password)
	if !ok {
		return nil, pkg.BadRequest(util.ERR_CODE_BODY_INVALIDO, msg)
	}

	user, err := s.repo.FindUserForAuth(ctx, in.Email)
	if err != nil {
		return nil, pkg.Internal("Error al buscar el usuario.")
	}
	if user == nil {
		return nil, pkg.Unauthorized(util.ERR_CODE_CREDENCIALES, util.MSG_CREDENCIALES_INVALIDAS)
	}
	if user.Estado != "ACTIVO" {
		return nil, pkg.Forbidden(util.ERR_CODE_USUARIO_INACT, util.MSG_USUARIO_INACTIVO)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)); err != nil {
		return nil, pkg.Unauthorized(util.ERR_CODE_CREDENCIALES, util.MSG_CREDENCIALES_INVALIDAS)
	}

	// Revocar tokens anteriores del usuario (sesión única)
	_ = s.repo.RevokeAllUserTokens(ctx, user.IDUsuario)

	return s.generateTokenPair(ctx, user)
}

// Refresh valida el refresh token y emite un nuevo par de tokens.
func (s *authService) Refresh(ctx context.Context, in input.RefreshInput) (*input.TokenPairOutput, error) {
	if in.RefreshToken == "" {
		return nil, pkg.BadRequest(util.ERR_CODE_BODY_INVALIDO, util.MSG_REFRESH_TOKEN_REQUERIDO)
	}

	hash := hashToken(in.RefreshToken)
	rt, err := s.repo.GetRefreshToken(ctx, hash)
	if err != nil {
		return nil, pkg.Internal("Error al verificar el token de refresco.")
	}
	if rt == nil {
		return nil, pkg.Unauthorized(util.ERR_CODE_TOKEN_INVALIDO, util.MSG_REFRESH_TOKEN_INVALIDO)
	}
	if rt.Revocado || time.Now().After(rt.ExpiresAt) {
		return nil, pkg.Unauthorized(util.ERR_CODE_TOKEN_INVALIDO, util.MSG_REFRESH_TOKEN_INVALIDO)
	}

	// Revocar el token usado (rotación de tokens)
	_ = s.repo.RevokeRefreshToken(ctx, hash)

	// Obtener datos actualizados del usuario por ID
	user, err := s.repo.FindUserByID(ctx, rt.IDUsuario)
	if err != nil || user == nil {
		return nil, pkg.Unauthorized(util.ERR_CODE_TOKEN_INVALIDO, util.MSG_REFRESH_TOKEN_INVALIDO)
	}
	if user.Estado != "ACTIVO" {
		return nil, pkg.Forbidden(util.ERR_CODE_USUARIO_INACT, util.MSG_USUARIO_INACTIVO)
	}

	return s.generateTokenPair(ctx, user)
}

// Logout revoca el refresh token del usuario.
func (s *authService) Logout(ctx context.Context, in input.RefreshInput) error {
	if in.RefreshToken == "" {
		return pkg.BadRequest(util.ERR_CODE_BODY_INVALIDO, util.MSG_REFRESH_TOKEN_REQUERIDO)
	}
	hash := hashToken(in.RefreshToken)
	return s.repo.RevokeRefreshToken(ctx, hash)
}

// ValidarToken valida la firma y expiración del access token y devuelve los claims.
func (s *authService) ValidarToken(tokenStr string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pkg.Unauthorized(util.ERR_CODE_TOKEN_INVALIDO, util.MSG_TOKEN_INVALIDO)
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, pkg.Unauthorized(util.ERR_CODE_TOKEN_INVALIDO, util.MSG_TOKEN_INVALIDO)
	}
	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, pkg.Unauthorized(util.ERR_CODE_TOKEN_INVALIDO, util.MSG_TOKEN_INVALIDO)
	}
	return claims, nil
}

// generateTokenPair genera access token + refresh token y persiste el refresh en BD.
func (s *authService) generateTokenPair(ctx context.Context, user *repository.UserAuthData) (*input.TokenPairOutput, error) {
	now := time.Now()

	claims := &AuthClaims{
		IDUsuario: user.IDUsuario,
		Email:     user.Email,
		Rol:       user.RolNombre,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "auth",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, pkg.Internal("Error al generar el token de acceso.")
	}

	// Refresh token: UUID aleatorio, se guarda su hash en BD
	rawRefresh := uuid.New().String()
	hash := hashToken(rawRefresh)

	rt := &domain.RefreshToken{
		IDUsuario: user.IDUsuario,
		TokenHash: hash,
		ExpiresAt: now.Add(refreshTokenDuration),
		Revocado:  false,
	}
	if err := s.repo.SaveRefreshToken(ctx, rt); err != nil {
		return nil, pkg.Internal("Error al guardar el token de refresco.")
	}

	return &input.TokenPairOutput{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessTokenDuration.Seconds()),
	}, nil
}

// hashToken calcula el SHA-256 del token para almacenarlo de forma segura.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
