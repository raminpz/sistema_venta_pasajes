package domain

import "time"

type Pago struct {
	IDPago    int64      `gorm:"column:ID_PAGO;primaryKey;autoIncrement" json:"id_pago"`
	IDVenta   int64      `gorm:"column:ID_VENTA" json:"id_venta"`
	IDMetodo  int64      `gorm:"column:ID_METODO" json:"id_metodo"`
	Monto     float64    `gorm:"column:MONTO" json:"monto"`
	Estado    string     `gorm:"column:ESTADO" json:"estado"`
	CreatedAt *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

func (Pago) TableName() string {
	return "PAGO"
}
