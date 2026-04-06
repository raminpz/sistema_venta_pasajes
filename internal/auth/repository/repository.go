package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"sistema_venta_pasajes/internal/auth/domain"
)

// UserAuthData agrupa los datos del usuario necesarios para la autenticación.
type UserAuthData struct {
	IDUsuario int
	Email     string
	Password  string
	Estado    string
	RolNombre string
}

// AuthRepository define las operaciones de persistencia del módulo auth.
type AuthRepository interface {
	FindUserForAuth(ctx context.Context, email string) (*UserAuthData, error)
	FindUserByID(ctx context.Context, id int) (*UserAuthData, error)
	SaveRefreshToken(ctx context.Context, rt *domain.RefreshToken) error
	GetRefreshToken(ctx context.Context, hash string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, hash string) error
	RevokeAllUserTokens(ctx context.Context, userID int) error
}

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository crea una instancia real del repositorio con MySQL.
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserForAuth(ctx context.Context, email string) (*UserAuthData, error) {
	var data UserAuthData
	err := r.db.WithContext(ctx).
		Table("USUARIO u").
		Select("u.ID_USUARIO, u.EMAIL, u.PASSWORD, u.ESTADO, r.NOMBRE as rol_nombre").
		Joins("JOIN ROL r ON u.ID_ROL = r.ID_ROL").
		Where("u.EMAIL = ?", email).
		Limit(1).
		Scan(&data).Error
	if err != nil {
		return nil, err
	}
	if data.IDUsuario == 0 {
		return nil, nil
	}
	return &data, nil
}

func (r *authRepository) FindUserByID(ctx context.Context, id int) (*UserAuthData, error) {
	var data UserAuthData
	err := r.db.WithContext(ctx).
		Table("USUARIO u").
		Select("u.ID_USUARIO, u.EMAIL, u.PASSWORD, u.ESTADO, r.NOMBRE as rol_nombre").
		Joins("JOIN ROL r ON u.ID_ROL = r.ID_ROL").
		Where("u.ID_USUARIO = ?", id).
		Limit(1).
		Scan(&data).Error
	if err != nil {
		return nil, err
	}
	if data.IDUsuario == 0 {
		return nil, nil
	}
	return &data, nil
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, rt *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(rt).Error
}

func (r *authRepository) GetRefreshToken(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var rt domain.RefreshToken
	err := r.db.WithContext(ctx).Where("TOKEN_HASH = ? AND REVOCADO = false", hash).First(&rt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rt, err
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, hash string) error {
	return r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("TOKEN_HASH = ?", hash).
		Update("REVOCADO", true).Error
}

func (r *authRepository) RevokeAllUserTokens(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("ID_USUARIO = ? AND REVOCADO = false", userID).
		Update("REVOCADO", true).Error
}
