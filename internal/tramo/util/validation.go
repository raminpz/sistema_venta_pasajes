package util

import (
	"sistema_venta_pasajes/internal/tramo/input"
	"sistema_venta_pasajes/pkg"
)

func ValidarCamposCreate(in input.CreateTramoInput) error {
	if in.IDRuta <= 0 {
		return pkg.BadRequest("required_ruta", ERR_REQUIRED_RUTA)
	}
	if in.IDParadaOrigen <= 0 {
		return pkg.BadRequest("required_parada_origen", ERR_REQUIRED_PARADA_ORIGEN)
	}
	if in.IDParadaDestino <= 0 {
		return pkg.BadRequest("required_parada_destino", ERR_REQUIRED_PARADA_DESTINO)
	}
	if in.IDParadaOrigen == in.IDParadaDestino {
		return pkg.BadRequest("paradas_iguales", ERR_PARADAS_IGUALES)
	}
	return nil
}

func ValidarCamposUpdate(in input.UpdateTramoInput) error {
	if in.IDRuta == nil && in.IDParadaOrigen == nil && in.IDParadaDestino == nil {
		return pkg.BadRequest("empty_update", ERR_EMPTY_UPDATE)
	}
	if in.IDRuta != nil && *in.IDRuta <= 0 {
		return pkg.BadRequest("required_ruta", ERR_REQUIRED_RUTA)
	}
	if in.IDParadaOrigen != nil && *in.IDParadaOrigen <= 0 {
		return pkg.BadRequest("required_parada_origen", ERR_REQUIRED_PARADA_ORIGEN)
	}
	if in.IDParadaDestino != nil && *in.IDParadaDestino <= 0 {
		return pkg.BadRequest("required_parada_destino", ERR_REQUIRED_PARADA_DESTINO)
	}
	if in.IDParadaOrigen != nil && in.IDParadaDestino != nil && *in.IDParadaOrigen == *in.IDParadaDestino {
		return pkg.BadRequest("paradas_iguales", ERR_PARADAS_IGUALES)
	}
	return nil
}
