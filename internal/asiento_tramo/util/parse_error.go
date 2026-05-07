package util

import (
	"sistema_venta_pasajes/pkg"
	"strings"
)

// ParseDBError convierte errores de BD en respuestas amigables
func ParseDBError(err error, defaultCode, defaultMsg string) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Duplicate entry (unique constraint)
	if strings.Contains(errStr, "Duplicate entry") || strings.Contains(errStr, "1062") {
		return pkg.BadRequest(ERR_CODE_ASIENTO_TRAMO_DUPLICATE, MSG_ASIENTO_TRAMO_DUPLICATE)
	}

	// Foreign key constraint
	if strings.Contains(errStr, "foreign key constraint") || strings.Contains(errStr, "1452") {
		return pkg.BadRequest(ERR_CODE_ASIENTO_TRAMO_CREATE, MSG_ASIENTO_TRAMO_CREATE_ERROR)
	}

	// Data truncated
	if strings.Contains(errStr, "Data truncated") {
		return pkg.BadRequest(defaultCode, "El valor proporcionado no es válido para este campo")
	}

	// Default error
	return pkg.NewAppError(500, defaultCode, defaultMsg).WithCause(err)
}
