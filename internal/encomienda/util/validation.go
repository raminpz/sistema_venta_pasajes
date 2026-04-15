package util

import (
	"net/http"

	"sistema_venta_pasajes/pkg"
)

func IsValidEstado(estado string) bool {
	switch estado {
	case STATUS_PENDIENTE, STATUS_EN_TRANSITO, STATUS_ENTREGADO, STATUS_DEVUELTO:
		return true
	default:
		return false
	}
}

func ParsePaginationParams(r *http.Request) (int, int, error) {
	return pkg.ParsePaginationParams(r, 1, 15)
}

func ParseDBError(err error, errCode, genericMsg string) error {
	fkMessages := map[string]string{
		"FK_ENCOMIENDA_VENTA":        MSG_ENCOMIENDA_VENTA_NOT_FOUND,
		"FK_ENCOMIENDA_PROGRAMACION": MSG_ENCOMIENDA_PROG_NOT_FOUND,
		"*":                          MSG_ENCOMIENDA_FOREIGN_KEY_ERROR,
	}
	return pkg.ParseDBError(err, errCode, genericMsg, fkMessages, nil)
}
