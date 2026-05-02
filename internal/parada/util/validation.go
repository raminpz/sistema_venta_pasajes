package util

import (
	"sistema_venta_pasajes/internal/parada/input"
	"sistema_venta_pasajes/pkg"
)

func ValidarCamposCreate(in input.CreateParadaInput) error {
	if in.IDRuta <= 0 {
		return pkg.BadRequest("required_ruta", ERR_REQUIRED_RUTA)
	}
	if in.IDTerminal <= 0 {
		return pkg.BadRequest("required_terminal", ERR_REQUIRED_TERMINAL)
	}
	if in.Orden <= 0 {
		return pkg.BadRequest("required_orden", ERR_REQUIRED_ORDEN)
	}
	return nil
}

func ValidarCamposUpdate(in input.UpdateParadaInput) error {
	if in.IDTerminal == nil && in.Orden == nil {
		return pkg.BadRequest("empty_update", ERR_EMPTY_UPDATE)
	}
	if in.IDTerminal != nil && *in.IDTerminal <= 0 {
		return pkg.BadRequest("required_terminal", ERR_REQUIRED_TERMINAL)
	}
	if in.Orden != nil && *in.Orden <= 0 {
		return pkg.BadRequest("required_orden", ERR_REQUIRED_ORDEN)
	}
	return nil
}

