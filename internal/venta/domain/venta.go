package domain

import (
	"time"
)

type Venta struct {
	IDVenta           int64     `gorm:"column:ID_VENTA;primaryKey;autoIncrement" json:"id_venta"`
	IDUsuario         int64     `gorm:"column:ID_USUARIO" json:"id_usuario"`
	IDTipoComprobante int64     `gorm:"column:ID_TIPO_COMPROBANTE" json:"id_tipo_comprobante"`
	Serie             string    `gorm:"column:SERIE" json:"serie"`
	Correlativo       uint      `gorm:"column:CORRELATIVO" json:"correlativo"`
	Nota              string    `gorm:"column:NOTA" json:"nota"`
	Observaciones     string    `gorm:"column:OBSERVACIONES" json:"observaciones"`
	QRCode            string    `gorm:"column:QR_CODE" json:"qr_code"`
	Subtotal          float64   `gorm:"column:SUBTOTAL" json:"subtotal"`
	IGV               float64   `gorm:"column:IGV" json:"igv"`
	Total             float64   `gorm:"column:TOTAL" json:"total"`
	Estado            string    `gorm:"column:ESTADO" json:"estado"`
	CreatedAt         time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

func (Venta) TableName() string {
	return "VENTA"
}
