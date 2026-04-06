package domain

import "time"

type LiquidacionViaje struct {
	IDLiquidacion    int64      `gorm:"column:ID_LIQUIDACION;primaryKey;autoIncrement"`
	IDProgramacion   int64      `gorm:"column:ID_PROGRAMACION;not null;uniqueIndex"`
	IDConductor      int64      `gorm:"column:ID_CONDUCTOR;not null"`
	TotalPasajes     float64    `gorm:"column:TOTAL_PASAJES;not null;default:0"`
	TotalEncomiendas float64    `gorm:"column:TOTAL_ENCOMIENDAS;not null;default:0"`
	TotalCaja        float64    `gorm:"column:TOTAL_CAJA;not null;default:0"`
	Estado           string     `gorm:"column:ESTADO;not null;default:PENDIENTE"`
	FechaLiquidacion *time.Time `gorm:"column:FECHA_LIQUIDACION"`
	Observaciones    string     `gorm:"column:OBSERVACIONES"`
	CreatedAt        *time.Time `gorm:"column:CREATED_AT;autoCreateTime"`
	UpdatedAt        *time.Time `gorm:"column:UPDATED_AT;autoUpdateTime"`
}

func (LiquidacionViaje) TableName() string {
	return "LIQUIDACION_VIAJE"
}
