package util

import "errors"

func ValidateEstadoAsiento(estado string) error {
	switch estado {
	case STATUS_SEAT_ACTIVE, STATUS_SEAT_OCCUPIED, STATUS_SEAT_RESERVED:
		return nil
	default:
		return errors.New("El estado del asiento es inválido. Debe ser ACTIVO, RESERVADO u OCUPADO.")
	}
}
