package domain

import "time"

type Vehiculo struct {
	IDVehiculo           int64      `gorm:"column:ID_VEHICULO;primaryKey;autoIncrement" json:"id_vehiculo"`
	IDTipoVehiculo       int64      `gorm:"column:ID_TIPO_VEHICULO" json:"id_tipo_vehiculo"`
	NroPlaca             string     `gorm:"column:NRO_PLACA" json:"nro_placa"`
	Marca                string     `gorm:"column:MARCA" json:"marca"`
	Modelo               string     `gorm:"column:MODELO" json:"modelo"`
	AnioFabricacion      int        `gorm:"column:ANIO_FABRICACION" json:"anio_fabricacion"`
	NumeroChasis         string     `gorm:"column:NUMERO_CHASIS" json:"numero_chasis"`
	Capacidad            int        `gorm:"column:CAPACIDAD" json:"capacidad"`
	NroSoat              string     `gorm:"column:NRO_SOAT" json:"nro_soat"`
	FechaVencSoat        *time.Time `gorm:"column:FECHA_VENC_SOAT" json:"fecha_venc_soat"`
	NroRevisionTecnica   string     `gorm:"column:NRO_REVISION_TECNICA" json:"nro_revision_tecnica"`
	FechaVencRevisionTec *time.Time `gorm:"column:FECHA_VENC_REV_TECNICA" json:"fecha_venc_revision_tecnica"`
	Estado               string     `gorm:"column:ESTADO" json:"estado"`
}

func (Vehiculo) TableName() string {
	return "VEHICULO"
}
