package domain

import "time"

type Conductor struct {
	IDConductor     int64      `json:"id_conductor" gorm:"column:ID_CONDUCTOR;primaryKey;autoIncrement"`
	Nombres         string     `json:"nombres" gorm:"column:NOMBRES"`
	Apellidos       string     `json:"apellidos" gorm:"column:APELLIDOS"`
	DNI             string     `json:"dni" gorm:"column:DNI"`
	NumeroLicencia  string     `json:"numero_licencia" gorm:"column:NUMERO_LICENCIA"`
	Telefono        string     `json:"telefono" gorm:"column:TELEFONO"`
	Direccion       *string    `json:"direccion" gorm:"column:DIRECCION"`
	FechaVencLicencia time.Time   `json:"fecha_venc_licencia" gorm:"column:FECHA_VENC_LICENCIA"`
	CreatedAt       *time.Time `json:"created_at" gorm:"column:CREATED_AT"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"column:UPDATED_AT"`
}

func (Conductor) TableName() string {
	return "CONDUCTOR"
}
