package domain

import (
	"time"
)

type ProveedorSistema struct {
	IDProveedor     int64      `json:"id_proveedor" gorm:"column:ID_PROVEEDOR;primaryKey;autoIncrement"`
	RUC             string     `json:"ruc" gorm:"column:RUC"`
	RazonSocial     string     `json:"razon_social" gorm:"column:RAZON_SOCIAL"`
	NombreComercial string     `json:"nombre_comercial,omitempty" gorm:"column:NOMBRE_COMERCIAL"`
	Direccion       string     `json:"direccion,omitempty" gorm:"column:DIRECCION"`
	Telefono        string     `json:"telefono,omitempty" gorm:"column:TELEFONO"`
	Email           string     `json:"email,omitempty" gorm:"column:EMAIL"`
	Web             string     `json:"web,omitempty" gorm:"column:WEB"`
	CreatedAt       *time.Time `json:"-" gorm:"column:CREATED_AT"`
	UpdatedAt       *time.Time `json:"-" gorm:"column:UPDATED_AT"`
}

func (ProveedorSistema) TableName() string {
	return "PROVEEDOR_SISTEMA"
}
