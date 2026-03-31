package util

import "errors"

func ValidateEstadoAsiento(estado string) error {
	switch estado {
	case EstadoAsientoActivo, EstadoAsientoOcupado, EstadoAsientoReservado:
		return nil
	default:
		return errors.New("El estado del asiento es inválido. Debe ser ACTIVO, RESERVADO u OCUPADO.")
	}
}
