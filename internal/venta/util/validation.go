package util

import (
	"errors"
	"net/http"
	"strconv"
)

func ValidarVentaInput(idUsuario, idTipoComprobante int64, subtotal float64) bool {
	return idUsuario > 0 && idTipoComprobante > 0 && subtotal > 0
}

// SerieFromTipoComprobante determina la serie automáticamente según el tipo de comprobante.
// 1=BOLETA → B001, 2=FACTURA → F001, 3=TICKET → T001
func SerieFromTipoComprobante(idTipo int64) (string, error) {
	switch idTipo {
	case 1:
		return "B001", nil
	case 2:
		return "F001", nil
	case 3:
		return "T001", nil
	default:
		return "", errors.New(MsgVentaErrorTipoComprob)
	}
}

func ParsePaginationParams(r *http.Request) (int, int, error) {
	page := 1
	size := 15

	if p := r.URL.Query().Get("page"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 {
			return 0, 0, errors.New("parámetro 'page' inválido")
		}
		page = v
	}

	if s := r.URL.Query().Get("size"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 {
			return 0, 0, errors.New("parámetro 'size' inválido")
		}
		size = v
	}

	return page, size, nil
}
