package domain

import "time"

type Programacion struct {
	IDProgramacion int64      `gorm:"column:ID_PROGRAMACION;primaryKey;autoIncrement" json:"id_programacion"`
	IDRuta         int64      `gorm:"column:ID_RUTA" json:"id_ruta"`
	IDVehiculo     int64      `gorm:"column:ID_VEHICULO" json:"id_vehiculo"`
	IDConductor    int64      `gorm:"column:ID_CONDUCTOR" json:"id_conductor"`
	FechaSalida    time.Time  `gorm:"column:FECHA_SALIDA" json:"fecha_salida"`
	FechaLlegada   *time.Time `gorm:"column:FECHA_LLEGADA" json:"fecha_llegada"`
	Estado         string     `gorm:"column:ESTADO" json:"estado"`
	CreatedAt      *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

func (Programacion) TableName() string {
	return "PROGRAMACION"
}
