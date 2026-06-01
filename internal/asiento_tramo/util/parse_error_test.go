package util

import (
	"errors"
	"net/http"
	"testing"

	"sistema_venta_pasajes/pkg"

	"github.com/stretchr/testify/assert"
)

func TestParseDBError_Nil(t *testing.T) {
	err := ParseDBError(nil, "code", "msg")
	assert.NoError(t, err)
}

func TestParseDBError_DuplicateEntry(t *testing.T) {
	err := ParseDBError(errors.New("Duplicate entry '12345' for key 'ASIENTO_TRAMO.UQ_ASIENTO_TRAMO'"), "code", "msg")
	assert.Error(t, err)
	appErr := pkg.AsAppError(err)
	assert.Equal(t, http.StatusBadRequest, appErr.Status)
	assert.Equal(t, ERR_CODE_ASIENTO_TRAMO_DUPLICATE, appErr.Code)
	assert.Equal(t, MSG_ASIENTO_TRAMO_DUPLICATE, appErr.Message)
}

func TestParseDBError_DuplicateEntry_ByCode1062(t *testing.T) {
	err := ParseDBError(errors.New("ERROR 1062: duplicate entry"), "code", "msg")
	assert.Error(t, err)
	appErr := pkg.AsAppError(err)
	assert.Equal(t, http.StatusBadRequest, appErr.Status)
	assert.Equal(t, ERR_CODE_ASIENTO_TRAMO_DUPLICATE, appErr.Code)
}

func TestParseDBError_ForeignKeyConstraint(t *testing.T) {
	err := ParseDBError(errors.New("foreign key constraint fails"), "code", "msg")
	assert.Error(t, err)
	appErr := pkg.AsAppError(err)
	assert.Equal(t, http.StatusBadRequest, appErr.Status)
	assert.Equal(t, ERR_CODE_ASIENTO_TRAMO_CREATE, appErr.Code)
}

func TestParseDBError_ForeignKey_ByCode1452(t *testing.T) {
	err := ParseDBError(errors.New("error 1452: cannot add child row"), "code", "msg")
	assert.Error(t, err)
	appErr := pkg.AsAppError(err)
	assert.Equal(t, http.StatusBadRequest, appErr.Status)
	assert.Equal(t, ERR_CODE_ASIENTO_TRAMO_CREATE, appErr.Code)
}

func TestParseDBError_DataTruncated(t *testing.T) {
	err := ParseDBError(errors.New("Data truncated for column 'ESTADO'"), "code", "msg")
	assert.Error(t, err)
	appErr := pkg.AsAppError(err)
	assert.Equal(t, http.StatusBadRequest, appErr.Status)
	assert.Equal(t, "code", appErr.Code)
	assert.Contains(t, appErr.Message, "válido")
}

func TestParseDBError_Generic(t *testing.T) {
	err := ParseDBError(errors.New("connection refused"), "my_code", "mi mensaje de error")
	assert.Error(t, err)
	appErr := pkg.AsAppError(err)
	assert.Equal(t, http.StatusInternalServerError, appErr.Status)
	assert.Equal(t, "my_code", appErr.Code)
	assert.Equal(t, "mi mensaje de error", appErr.Message)
}

