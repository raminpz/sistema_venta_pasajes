package util

import (
	"sistema_venta_pasajes/internal/tipo_vehiculo/input"
	"sistema_venta_pasajes/pkg"
)

func ValidateCreate(in input.CreateTipoVehiculoInput) error {
	if in.Nombre == "" {
		return pkg.BadRequest("required_nombre", ERR_REQUIRED_NAME)
	}
	if in.Descripcion == "" {
		return pkg.BadRequest("required_descripcion", ERR_REQUIRED_DESC)
	}
	return nil
}

func ValidateUpdate(in input.UpdateTipoVehiculoInput) error {
	if in.Nombre == nil && in.Descripcion == nil {
		return pkg.BadRequest("empty_update", ERR_EMPTY_UPDATE)
	}
	if in.Nombre != nil && *in.Nombre == "" {
		return pkg.BadRequest("required_nombre", ERR_REQUIRED_NAME)
	}
	if in.Descripcion != nil && *in.Descripcion == "" {
		return pkg.BadRequest("required_descripcion", ERR_REQUIRED_DESC)
	}
	return nil
}

