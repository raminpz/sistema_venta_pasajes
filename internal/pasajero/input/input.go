package input

import (
	"errors"
	"sistema_venta_pasajes/internal/pasajero/util"
	"sistema_venta_pasajes/pkg"
)

// CreatePasajeroInput Input para crear pasajero
type CreatePasajeroInput struct {
	TipoDocumento   string  `json:"tipo_documento" binding:"required"`
	NroDocumento    string  `json:"nro_documento" binding:"required"`
	Nombres         string  `json:"nombres" binding:"required"`
	Apellidos       string  `json:"apellidos" binding:"required"`
	Telefono        string  `json:"telefono" binding:"required"`
	Email           *string `json:"email"`
	FechaNacimiento *string `json:"fecha_nacimiento"`
}

// Input para actualizar pasajero
// Todos los campos opcionales excepto el identificador en la ruta
// (puedes ajustar según tu lógica de negocio)
type UpdatePasajeroInput struct {
	TipoDocumento   string  `json:"tipo_documento"`
	NroDocumento    string  `json:"nro_documento"`
	Nombres         string  `json:"nombres"`
	Apellidos       string  `json:"apellidos"`
	Telefono        string  `json:"telefono"`
	Email           *string `json:"email"`
	FechaNacimiento *string `json:"fecha_nacimiento"`
}

// Output para respuesta de pasajero
// Puedes ajustar los tipos según tu dominio
// IDPasajero es int64 para consistencia con terminal
// Fechas como string (formato ISO) para respuesta JSON

type PasajeroOutput struct {
	IDPasajero      int64   `json:"id_pasajero"`
	TipoDocumento   string  `json:"tipo_documento"`
	NroDocumento    string  `json:"nro_documento"`
	Nombres         string  `json:"nombres"`
	Apellidos       string  `json:"apellidos"`
	Telefono        string  `json:"telefono"`
	Email           *string `json:"email"`
	FechaNacimiento *string `json:"fecha_nacimiento"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func (in *CreatePasajeroInput) Validate() error {
	if in.TipoDocumento == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": tipo_documento")
	}
	if in.NroDocumento == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": nro_documento")
	}
	if in.Nombres == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": nombres")
	}
	if in.Apellidos == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": apellidos")
	}
	if !pkg.ValidatePhone(in.Telefono) {
		return errors.New(util.ERR_PHONE_FORMAT)
	}
	if in.Email != nil && *in.Email != "" && !pkg.ValidateEmail(*in.Email) {
		return errors.New(util.ERR_EMAIL_FORMAT)
	}
	if in.FechaNacimiento != nil && *in.FechaNacimiento != "" && !pkg.ValidateDate(*in.FechaNacimiento) {
		return errors.New(util.ERR_DATE_FORMAT)
	}
	return nil
}

func (in *UpdatePasajeroInput) Validate() error {
	if in.Telefono != "" && !pkg.ValidatePhone(in.Telefono) {
		return errors.New(util.ERR_PHONE_FORMAT)
	}
	if in.Email != nil && *in.Email != "" && !pkg.ValidateEmail(*in.Email) {
		return errors.New(util.ERR_EMAIL_FORMAT)
	}
	if in.FechaNacimiento != nil && *in.FechaNacimiento != "" && !pkg.ValidateDate(*in.FechaNacimiento) {
		return errors.New(util.ERR_DATE_FORMAT)
	}
	return nil
}
