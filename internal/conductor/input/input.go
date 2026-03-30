package input

import (
	"errors"
	"sistema_venta_pasajes/internal/conductor/util"
	"sistema_venta_pasajes/pkg"
	"time"
)

type CreateConductorInput struct {
	Nombres           string  `json:"nombres" binding:"required"`
	Apellidos         string  `json:"apellidos" binding:"required"`
	DNI               string  `json:"dni" binding:"required"`
	NumeroLicencia    string  `json:"numero_licencia" binding:"required"`
	Telefono          string  `json:"telefono" binding:"required"`
	Direccion         *string `json:"direccion"`
	FechaVencLicencia string  `json:"fecha_venc_licencia" binding:"required"`
}

type UpdateConductorInput struct {
	Nombres        *string `json:"nombres"`
	Apellidos      *string `json:"apellidos"`
	DNI            *string `json:"dni"`
	NumeroLicencia *string `json:"numero_licencia"`
	Telefono       *string `json:"telefono"`
	Direccion      *string `json:"direccion"`
}

type ConductorOutput struct {
	IDConductor       int64      `json:"id_conductor"`
	Nombres           string     `json:"nombres"`
	Apellidos         string     `json:"apellidos"`
	DNI               string     `json:"dni"`
	NumeroLicencia    string     `json:"numero_licencia"`
	Telefono          string     `json:"telefono"`
	Direccion         *string    `json:"direccion"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	FechaVencLicencia time.Time  `json:"fecha_venc_licencia"`
}

func (in *CreateConductorInput) Validate() error {
	if in.Nombres == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": nombres")
	}
	if in.Apellidos == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": apellidos")
	}
	if in.DNI == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": dni")
	}
	if len(in.DNI) != 8 {
		return errors.New("El DNI debe tener exactamente 8 dígitos")
	}
	if in.NumeroLicencia == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": numero_licencia")
	}
	if !util.IsValidNumeroLicencia(in.NumeroLicencia) {
		return errors.New(util.ERR_NUMERO_LICENCIA_FORMAT)
	}
	if in.Telefono == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": telefono")
	}
	if !pkg.ValidatePhone(in.Telefono) {
		return errors.New("El teléfono debe tener exactamente 9 dígitos numéricos")
	}
	if in.FechaVencLicencia == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": fecha_venc_licencia")
	}
	return nil
}

func (in *UpdateConductorInput) Validate() error {
	if in.DNI != nil && *in.DNI != "" && len(*in.DNI) != 8 {
		return errors.New("El DNI debe tener exactamente 8 dígitos")
	}
	if in.NumeroLicencia != nil && *in.NumeroLicencia != "" && !util.IsValidNumeroLicencia(*in.NumeroLicencia) {
		return errors.New(util.ERR_NUMERO_LICENCIA_FORMAT)
	}
	if in.Telefono != nil && *in.Telefono != "" && !pkg.ValidatePhone(*in.Telefono) {
		return errors.New("El teléfono debe tener exactamente 9 dígitos numéricos")
	}
	return nil
}
