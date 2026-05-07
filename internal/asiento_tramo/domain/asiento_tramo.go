package domain

import "time"

type AsientoTramo struct {
	IDAsientoTramo int64      `gorm:"column:ID_ASIENTO_TRAMO;primaryKey;autoIncrement" json:"id_asiento_tramo"`
	IDVenta        *int64     `gorm:"column:ID_VENTA" json:"id_venta,omitempty"`
	IDAsiento      int64      `gorm:"column:ID_ASIENTO" json:"id_asiento"`
	IDTramo        int64      `gorm:"column:ID_TRAMO" json:"id_tramo"`
	Estado         string     `gorm:"column:ESTADO" json:"estado"`
	CreatedAt      *time.Time `gorm:"column:CREATED_AT" json:"created_at,omitempty"`
	UpdatedAt      *time.Time `gorm:"column:UPDATED_AT" json:"updated_at,omitempty"`
}

func (AsientoTramo) TableName() string {
	return "ASIENTO_TRAMO"
}
