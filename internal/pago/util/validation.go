package util

import (
	"errors"
	"net/http"
	"strings"

	"sistema_venta_pasajes/pkg"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

func IsValidEstado(estado string) bool {
	switch estado {
	case STATUS_REGISTRADA, STATUS_PARCIAL, STATUS_PAGADA, STATUS_ANULADA:
		return true
	default:
		return false
	}
}

// ParseDBError convierte un error de base de datos en un AppError con mensaje claro.
// Identifica el tipo exacto de error MySQL (FK, ENUM truncado, etc.) para pago.
func ParseDBError(err error, errCode string, genericMsg string) error {
	if err == nil {
		return nil
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) {
		// Error 1452: No matching row in referenced table (Foreign Key constraint)
		// Error 1451: Cannot delete or update a parent row (Foreign Key constraint)
		if mysqlErr.Number == 1452 || mysqlErr.Number == 1451 {
			errText := strings.ToUpper(mysqlErr.Message)

			// Buscar constraint name en el mensaje de error
			if strings.Contains(errText, "FK_PAGO_VENTA") {
				return pkg.Conflict(errCode, MSG_PAGO_VENTA_NOT_FOUND)
			}
			if strings.Contains(errText, "FK_PAGO_METODO") {
				return pkg.Conflict(errCode, MSG_PAGO_METODO_NOT_FOUND)
			}

			// Si es FK pero no coincide con nuestras constraints conocidas
			return pkg.Conflict(errCode, MSG_PAGO_FOREIGN_KEY_CONFLICT)
		}

		// Error 1265: Data truncated - ENUM value not in allowed list
		if mysqlErr.Number == 1265 {
			return pkg.BadRequest(errCode, MSG_PAGO_ENUM_DB_ERROR)
		}
	}

	// Si ya es AppError, pasarlo tal cual
	appErr := pkg.AsAppError(err)
	if appErr != nil {
		return appErr
	}

	// Si es un error desconocido, retornar el genérico
	return pkg.NewAppError(http.StatusInternalServerError, errCode, genericMsg)
}

func ParsePaginationParams(r *http.Request) (int, int, error) {
	return pkg.ParsePaginationParams(r, 1, 15)
}
