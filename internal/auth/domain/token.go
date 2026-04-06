package domain

import "time"

// RefreshToken representa el token de refresco persistido en BD.
type RefreshToken struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:ID"`
	IDUsuario int       `gorm:"column:ID_USUARIO;not null"`
	TokenHash string    `gorm:"column:TOKEN_HASH;not null;unique"`
	ExpiresAt time.Time `gorm:"column:EXPIRES_AT;not null"`
	Revocado  bool      `gorm:"column:REVOCADO;not null;default:false"`
	CreatedAt time.Time `gorm:"column:CREATED_AT;autoCreateTime"`
}

func (RefreshToken) TableName() string {
	return "REFRESH_TOKEN"
}
