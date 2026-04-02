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
			if strings.Contains(errText, "FK_ENCOMIENDA_VENTA") {
				return pkg.Conflict(errCode, MSG_ENCOMIENDA_VENTA_NOT_FOUND)
			}
			if strings.Contains(errText, "FK_ENCOMIENDA_PROGRAMACION") {
				return pkg.Conflict(errCode, MSG_ENCOMIENDA_PROG_NOT_FOUND)
			}

			// Si es FK pero no coincide con nuestras constraints conocidas, retornar mensaje genérico
			return pkg.Conflict(errCode, MSG_ENCOMIENDA_FOREIGN_KEY_ERROR)
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
