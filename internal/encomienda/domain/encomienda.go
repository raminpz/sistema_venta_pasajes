package domain

import "time"

type Encomienda struct {
	IDEncomienda       int64      `gorm:"column:ID_ENCOMIENDA;primaryKey;autoIncrement" json:"id_encomienda"`
	IDVenta            int64      `gorm:"column:ID_VENTA" json:"id_venta"`
	IDProgramacion     int64      `gorm:"column:ID_PROGRAMACION" json:"id_programacion"`
	Descripcion        string     `gorm:"column:DESCRIPCION" json:"descripcion"`
	PesoKg             *float64   `gorm:"column:PESO_KG" json:"peso_kg"`
	Costo              float64    `gorm:"column:COSTO" json:"costo"`
	RemitenteNombre    string     `gorm:"column:REMITENTE_NOMBRE" json:"remitente_nombre"`
	RemitenteDoc       string     `gorm:"column:REMITENTE_DOC" json:"remitente_doc"`
	DestinatarioNombre string     `gorm:"column:DESTINATARIO_NOMBRE" json:"destinatario_nombre"`
	DestinatarioDoc    *string    `gorm:"column:DESTINATARIO_DOC" json:"destinatario_doc"`
	DestinatarioTel    string     `gorm:"column:DESTINATARIO_TEL" json:"destinatario_tel"`
	Estado             string     `gorm:"column:ESTADO" json:"estado"`
	CreatedAt          *time.Time `gorm:"column:CREATED_AT" json:"created_at"`
	UpdatedAt          *time.Time `gorm:"column:UPDATED_AT" json:"updated_at"`
}

func (Encomienda) TableName() string {
	return "ENCOMIENDA"
}
