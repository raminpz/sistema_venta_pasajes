package util

import (
	"errors"
	"net/http"

	"sistema_venta_pasajes/pkg"
)

// ValidarCreateInput valida los campos obligatorios del input de creación.
func ValidarCreateInput(idUsuario, idTipoComprobante, idProgramacion, idPasajero, idAsiento int64, precio float64, descuento *float64) error {
	details := map[string]string{}
	if idUsuario <= 0 {
		details["id_usuario"] = MSG_VENTA_USUARIO_REQUIRED
	}
	if idTipoComprobante <= 0 {
		details["id_tipo_comprobante"] = MSG_VENTA_COMPROBANTE_REQUIRED
	}
	if idProgramacion <= 0 {
		details["id_programacion"] = MSG_VENTA_PROGRAMACION_REQUIRED
	}
	if idPasajero <= 0 {
		details["id_pasajero"] = MSG_VENTA_PASAJERO_REQUIRED
	}
	if idAsiento <= 0 {
		details["id_asiento"] = MSG_VENTA_ASIENTO_REQUIRED
	}
	if precio < 0 {
		details["precio"] = MSG_VENTA_PRECIO_INVALID
	}
	if descuento != nil {
		if *descuento < 0 || *descuento > precio {
			details["descuento"] = MSG_VENTA_DESCUENTO_INVALID
		}
	}
	if len(details) > 0 {
		return pkg.Validation(MSG_VENTA_VALIDATION_ERROR, details)
	}
	return nil
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

// ParseDBError interpreta errores MySQL de VENTA con mensajes claros.
func ParseDBError(err error, errCode, genericMsg string) error {
	fkMessages := map[string]string{
		"FK_VENTA_PROGRAMACION":     MSG_VENTA_FK_PROGRAMACION,
		"FK_VENTA_PASAJERO":         MSG_VENTA_FK_PASAJERO,
		"FK_VENTA_ASIENTO":          MSG_VENTA_FK_ASIENTO,
		"FK_VENTA_USUARIO":          MSG_VENTA_FK_USUARIO,
		"FK_VENTA_TIPO_COMPROBANTE": MSG_VENTA_FK_COMPROBANTE,
	}
	duplicateMessages := map[string]string{
		"UQ_VENTA_PROG_ASIENTO":  MSG_VENTA_DUPLICATE_ASIENTO,
		"UQ_VENTA_PROG_PASAJERO": MSG_VENTA_DUPLICATE_PASAJERO,
	}
	return pkg.ParseDBError(err, errCode, genericMsg, fkMessages, duplicateMessages)
}

func ParsePaginationParams(r *http.Request) (int, int, error) {
	return pkg.ParsePaginationParams(r, 1, 15)
}
