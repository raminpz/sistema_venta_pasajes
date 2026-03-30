package domain

import "time"

type Empresa struct {
	IDEmpresa       int       `json:"id_empresa" gorm:"column:ID_EMPRESA;primaryKey"`
	RUC             string    `json:"ruc" gorm:"column:RUC"`
	RazonSocial     string    `json:"razon_social" gorm:"column:RAZON_SOCIAL"`
	NombreComercial *string   `json:"nombre_comercial" gorm:"column:NOMBRE_COMERCIAL"`
	Direccion       *string   `json:"direccion" gorm:"column:DIRECCION"`
	Telefono        string    `json:"telefono" gorm:"column:TELEFONO"`
	Email           *string   `json:"email" gorm:"column:EMAIL"`
	Logo            *string   `json:"logo" gorm:"column:LOGO"`
	FechaCreacion   time.Time `json:"fecha_creacion" gorm:"column:FECHA_CREACION"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:CREATED_AT"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:UPDATED_AT"`
}

func (Empresa) TableName() string {
	return "EMPRESA"
}
