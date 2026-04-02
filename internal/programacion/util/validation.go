package util

import (
	"errors"
	"net/http"
	"sistema_venta_pasajes/pkg"
	"time"
)

var dateTimeLayouts = []string{
	"2006-01-02 15:04:05",
	time.RFC3339,
}

func ParseDateTime(value string) (time.Time, error) {
	for _, layout := range dateTimeLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New(MSG_PROGRAMACION_INVALID_DATETIME)
}

func IsValidEstado(estado string) bool {
	switch estado {
	case STATUS_PROGRAMADO, STATUS_EN_CURSO, STATUS_COMPLETADO, STATUS_CANCELADO:
		return true
	default:
		return false
	}
}

func ParsePaginationParams(r *http.Request) (int, int, error) {
	return pkg.ParsePaginationParams(r, 1, 15)
}
