package domain

import "time"

type Asiento struct {
	IDAsiento     int        `gorm:"column:ID_ASIENTO;primaryKey;autoIncrement" json:"id_asiento"`
	IDVehiculo    int        `gorm:"column:ID_VEHICULO" json:"id_vehiculo"`
	NumeroAsiento string     `gorm:"column:NUMERO_ASIENTO" json:"numero_asiento"`
	Estado        string     `gorm:"column:ESTADO" json:"estado"`
	CreatedAt     *time.Time `gorm:"column:CREATED_AT" json:"created_at,omitempty"`
	UpdatedAt     *time.Time `gorm:"column:UPDATED_AT" json:"updated_at,omitempty"`
}

func (Asiento) TableName() string {
	return "ASIENTO"
}
