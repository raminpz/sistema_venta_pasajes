package domain

import "time"

type Vehiculo struct {
	IDVehiculo           int64      `gorm:"column:ID_VEHICULO;primaryKey;autoIncrement" json:"id_vehiculo"`
	IDTipoVehiculo       int64      `gorm:"column:ID_TIPO_VEHICULO;not null" json:"id_tipo_vehiculo"`
	NroPlaca             string     `gorm:"column:NRO_PLACA;not null" json:"nro_placa"`
	Marca                string     `gorm:"column:MARCA;not null" json:"marca"`
	Modelo               string     `gorm:"column:MODELO;not null" json:"modelo"`
	AnioFabricacion      int        `gorm:"column:ANIO_FABRICACION;not null" json:"anio_fabricacion"`
	NumeroChasis         string     `gorm:"column:NUMERO_CHASIS;not null" json:"numero_chasis"`
	Capacidad            int        `gorm:"column:CAPACIDAD;not null" json:"capacidad"`
	NroSoat              string     `gorm:"column:NRO_SOAT;not null" json:"nro_soat"`
	FechaVencSoat        *time.Time `gorm:"column:FECHA_VENC_SOAT;not null" json:"fecha_venc_soat"`
	NroRevisionTecnica   string     `gorm:"column:NRO_REVISION_TECNICA;not null" json:"nro_revision_tecnica"`
	FechaVencRevisionTec *time.Time `gorm:"column:FECHA_VENC_REV_TECNICA;not null" json:"fecha_venc_revision_tecnica"`
	Estado               string     `gorm:"column:ESTADO;not null" json:"estado"`
}

func (Vehiculo) TableName() string {
	return "VEHICULO"
}
