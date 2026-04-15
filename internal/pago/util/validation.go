package util

import (
	"errors"
	"net/http"

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
	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1265 {
		return pkg.BadRequest(errCode, MSG_PAGO_ENUM_DB_ERROR).WithCause(err)
	}

	fkMessages := map[string]string{
		"FK_PAGO_VENTA":  MSG_PAGO_VENTA_NOT_FOUND,
		"FK_PAGO_METODO": MSG_PAGO_METODO_NOT_FOUND,
		"*":              MSG_PAGO_FOREIGN_KEY_CONFLICT,
	}
	return pkg.ParseDBError(err, errCode, genericMsg, fkMessages, nil)
}

func ParsePaginationParams(r *http.Request) (int, int, error) {
	return pkg.ParsePaginationParams(r, 1, 15)
}
