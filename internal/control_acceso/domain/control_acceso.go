package domain

import "time"

type ControlAcceso struct {
	IDAcceso        int64      `gorm:"column:ID_ACCESO;primaryKey;autoIncrement"`
	FechaActivacion time.Time  `gorm:"column:FECHA_ACTIVACION"`
	FechaExpiracion time.Time  `gorm:"column:FECHA_EXPIRACION"`
	Estado          string     `gorm:"column:ESTADO;default:OPERATIVO"`
	CreatedAt       *time.Time `gorm:"column:CREATED_AT"`
	UpdatedAt       *time.Time `gorm:"column:UPDATED_AT"`
}

func (ControlAcceso) TableName() string {
	return "CONTROL_ACCESO"
}
