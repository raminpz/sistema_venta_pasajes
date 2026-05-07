package util

import (
	"sistema_venta_pasajes/internal/parada/input"
	"sistema_venta_pasajes/pkg"
	"strings"
)

func ValidarCamposCreate(in input.CreateParadaInput) error {
	if in.IDRuta <= 0 {
		return pkg.BadRequest("required_ruta", ERR_REQUIRED_RUTA)
	}
	if strings.TrimSpace(in.NombreParada) == "" {
		return pkg.BadRequest("required_nombre_parada", ERR_REQUIRED_NOMBRE)
	}
	if in.Orden <= 0 {
		return pkg.BadRequest("required_orden", ERR_REQUIRED_ORDEN)
	}
	return nil
}

func ValidarCamposUpdate(in input.UpdateParadaInput) error {
	if in.NombreParada == nil && in.Orden == nil {
		return pkg.BadRequest("empty_update", ERR_EMPTY_UPDATE)
	}
	if in.NombreParada != nil && strings.TrimSpace(*in.NombreParada) == "" {
		return pkg.BadRequest("required_nombre_parada", ERR_REQUIRED_NOMBRE)
	}
	if in.Orden != nil && *in.Orden <= 0 {
		return pkg.BadRequest("required_orden", ERR_REQUIRED_ORDEN)
	}
	return nil
}
