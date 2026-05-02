package domain

import "time"

type Tramo struct {
	IDTramo         int64      `gorm:"column:ID_TRAMO;primaryKey;autoIncrement" json:"id_tramo"`
	IDRuta          int64      `gorm:"column:ID_RUTA" json:"id_ruta"`
	IDParadaOrigen  int64      `gorm:"column:ID_PARADA_ORIGEN" json:"id_parada_origen"`
	IDParadaDestino int64      `gorm:"column:ID_PARADA_DESTINO" json:"id_parada_destino"`
	CreatedAt       *time.Time `gorm:"column:CREATED_AT" json:"-"`
	UpdatedAt       *time.Time `gorm:"column:UPDATED_AT" json:"-"`
}

func (Tramo) TableName() string {
	return "TRAMO"
}
