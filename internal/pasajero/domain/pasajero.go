package domain

import "time"

type Pasajero struct {
	IDPasajero      int        `json:"id_pasajero" gorm:"column=ID_PASAJERO;primaryKey"`
	TipoDocumento   string     `json:"tipo_documento" gorm:"column:TIPO_DOCUMENTO"`
	NroDocumento    string     `json:"nro_documento"`
	Nombres         string     `json:"nombres"`
	Apellidos       string     `json:"apellidos" gorm:"column:APELLIDOS"`
	Telefono        string     `json:"telefono"`
	Email           *string    `json:"email,omitempty"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (Pasajero) TableName() string {
	return "PASAJERO"
}
