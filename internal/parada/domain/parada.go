package domain

import "time"

type Parada struct {
	IDParada     int64      `gorm:"column:ID_PARADA;primaryKey;autoIncrement" json:"id_parada"`
	IDRuta       int64      `gorm:"column:ID_RUTA" json:"id_ruta"`
	NombreParada string     `gorm:"column:NOMBRE_PARADA" json:"nombre_parada"`
	Orden        int        `gorm:"column:ORDEN" json:"orden"`
	CreatedAt    *time.Time `gorm:"column:CREATED_AT" json:"-"`
	UpdatedAt    *time.Time `gorm:"column:UPDATED_AT" json:"-"`
}

func (Parada) TableName() string {
	return "PARADA"
}
