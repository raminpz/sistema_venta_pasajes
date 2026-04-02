package util

import (
	"errors"
	"net/http"
	"sistema_venta_pasajes/pkg"
)

func ValidarVentaInput(idUsuario, idTipoComprobante int64, subtotal float64) bool {
	return idUsuario > 0 && idTipoComprobante > 0 && subtotal > 0
}

// SerieFromTipoComprobante determina la serie automáticamente según el tipo de comprobante.
// 1=BOLETA -> B001, 2=FACTURA -> F001, 3=TICKET -> T001
func SerieFromTipoComprobante(idTipo int64) (string, error) {
	switch idTipo {
	case 1:
		return "B001", nil
	case 2:
		return "F001", nil
	case 3:
		return "T001", nil
	default:
		return "", errors.New(MSG_VENTA_TIPO_COMPROBANTE_ERROR)
	}
}

func ParsePaginationParams(r *http.Request) (int, int, error) {
	return pkg.ParsePaginationParams(r, 1, 15)
}
