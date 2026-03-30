package domain

import (
	"time"
)

type LicenciaSistema struct {
	IDLicencia      int64      `json:"id_licencia" gorm:"column:ID_LICENCIA;primaryKey;autoIncrement"`
	ClaveLicencia   string     `json:"clave_licencia" gorm:"column:CLAVE_LICENCIA"`
	FechaActivacion time.Time  `json:"fecha_activacion" gorm:"column:FECHA_ACTIVACION"`
	FechaExpiracion time.Time  `json:"fecha_expiracion" gorm:"column:FECHA_EXPIRACION"`
	Estado          string     `json:"estado" gorm:"column:ESTADO;default:OPERATIVO"`
	FechaRegistro   time.Time  `json:"fecha_registro" gorm:"column:FECHA_REGISTRO"`
	CreatedAt       *time.Time `json:"created_at,omitempty" gorm:"column:CREATED_AT"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty" gorm:"column:UPDATED_AT"`
}

func (LicenciaSistema) TableName() string {
	return "LICENCIA_SISTEMA"
}
