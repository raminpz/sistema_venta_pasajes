package pkg

import (
	"errors"
	"net/http"
	"testing"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestParseDBError_DataTruncated_PropagaMensajeExacto(t *testing.T) {
	err := &mysqlDriver.MySQLError{Number: 1265, Message: "Data truncated for column 'NOMBRE' at row 1"}

	appErr := ParseDBError(err, "tipo_vehiculo_create_error", "Error al crear tipo de vehiculo", nil, nil)
	assert.Error(t, appErr)

	parsed := AsAppError(appErr)
	assert.Equal(t, http.StatusBadRequest, parsed.Status)
	assert.Equal(t, "tipo_vehiculo_create_error", parsed.Code)
	assert.Contains(t, parsed.Message, "Data truncated for column 'NOMBRE' at row 1")
	assert.Equal(t, "Data truncated for column 'NOMBRE' at row 1", parsed.Details)
}

func TestAsAppError_PreservesExistingAppError(t *testing.T) {
	input := BadRequest("x", "y").WithCause(errors.New("cause"))
	out := AsAppError(input)
	assert.Equal(t, input, out)
}

