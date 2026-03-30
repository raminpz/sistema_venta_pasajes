package input

import (
	"errors"
	"sistema_venta_pasajes/internal/empresa/util"
	"sistema_venta_pasajes/pkg"
)

type CreateEmpresaInput struct {
	RUC             string  `json:"ruc"`
	RazonSocial     string  `json:"razon_social"`
	NombreComercial *string `json:"nombre_comercial"`
	Direccion       *string `json:"direccion"`
	Telefono        string  `json:"telefono"`
	Email           *string `json:"email"`
	Logo            *string `json:"logo"`
	FechaCreacion   string  `json:"fecha_creacion"`
}

type UpdateEmpresaInput struct {
	RUC             *string `json:"ruc"`
	RazonSocial     string  `json:"razon_social"`
	NombreComercial *string `json:"nombre_comercial"`
	Direccion       *string `json:"direccion"`
	Telefono        string  `json:"telefono"`
	Email           *string `json:"email"`
	Logo            *string `json:"logo"`
	FechaCreacion   *string `json:"fecha_creacion"`
}

type EmpresaOutput struct {
	IDEmpresa       int64   `json:"id_empresa"`
	RUC             string  `json:"ruc"`
	RazonSocial     string  `json:"razon_social"`
	NombreComercial *string `json:"nombre_comercial"`
	Direccion       *string `json:"direccion"`
	Telefono        string  `json:"telefono"`
	Email           *string `json:"email"`
	Logo            *string `json:"logo"`
	FechaCreacion   string  `json:"fecha_creacion"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func (in *CreateEmpresaInput) Validate() error {
	if in.RUC == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": ruc")
	}
	if !pkg.ValidateRUC(in.RUC) {
		return errors.New(util.ERR_RUC_FORMAT)
	}
	if in.RazonSocial == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": razon_social")
	}
	if in.Telefono == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": telefono")
	}
	if !pkg.ValidatePhone(in.Telefono) {
		return errors.New(util.ERR_PHONE_FORMAT)
	}
	if in.Email != nil && *in.Email != "" && !pkg.ValidateEmail(*in.Email) {
		return errors.New(util.ERR_EMAIL_FORMAT)
	}
	if in.FechaCreacion == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": fecha_creacion")
	}
	if !pkg.ValidateDate(in.FechaCreacion) {
		return errors.New(util.ERR_DATE_FORMAT)
	}
	return nil
}

func (in *UpdateEmpresaInput) Validate() error {
	if in.RUC != nil && *in.RUC != "" && !pkg.ValidateRUC(*in.RUC) {
		return errors.New(util.ERR_RUC_FORMAT)
	}
	if in.RazonSocial == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": razon_social")
	}
	if in.Telefono == "" {
		return errors.New(util.ERR_REQUIRED_FIELD + ": telefono")
	}
	if !pkg.ValidatePhone(in.Telefono) {
		return errors.New(util.ERR_PHONE_FORMAT)
	}
	if in.Email != nil && *in.Email != "" && !pkg.ValidateEmail(*in.Email) {
		return errors.New(util.ERR_EMAIL_FORMAT)
	}
	if in.FechaCreacion != nil && *in.FechaCreacion != "" && !pkg.ValidateDate(*in.FechaCreacion) {
		return errors.New(util.ERR_DATE_FORMAT)
	}
	return nil
}
