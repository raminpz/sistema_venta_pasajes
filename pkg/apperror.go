package pkg

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *AppError) WithDetails(details any) *AppError {
	if e == nil {
		return nil
	}
	e.Details = details
	return e
}

func (e *AppError) WithCause(err error) *AppError {
	if e == nil {
		return nil
	}
	e.Err = err
	return e
}

func NewAppError(status int, code, message string) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func BadRequest(code, message string) *AppError {
	return NewAppError(http.StatusBadRequest, code, message)
}

func Validation(message string, details any) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, "validation_error", message).WithDetails(details)
}

func NotFound(code, message string) *AppError {
	return NewAppError(http.StatusNotFound, code, message)
}

func Conflict(code, message string) *AppError {
	return NewAppError(http.StatusConflict, code, message)
}

func MethodNotAllowed(message string) *AppError {
	return NewAppError(http.StatusMethodNotAllowed, "method_not_allowed", message)
}

func ServiceUnavailable(code, message string) *AppError {
	return NewAppError(http.StatusServiceUnavailable, code, message)
}

func Forbidden(code, message string) *AppError {
	return NewAppError(http.StatusForbidden, code, message)
}

func Unauthorized(code, message string) *AppError {
	return NewAppError(http.StatusUnauthorized, code, message)
}

func Internal(message string, _ ...interface{}) *AppError {
	return NewAppError(http.StatusInternalServerError, "internal_error", message)
}

func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	if errors.Is(err, sql.ErrNoRows) {
		return NotFound("resource_not_found", "recurso no encontrado").WithCause(err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound("resource_not_found", "recurso no encontrado").WithCause(err)
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062:
			return Conflict("duplicate_resource", "el recurso ya existe o ya fue registrado").WithCause(err)
		case 1451, 1452:
			return Conflict("foreign_key_conflict", "la operación viola restricciones de integridad referencial").WithCause(err)
		default:
			return Internal("Error en base de datos: " + mysqlErr.Message).WithDetails(mysqlErr.Message).WithCause(err)
		}
	}

	return Internal("ocurrió un error interno en el servidor").WithCause(err)
}

// Manejo de errores de decodificación JSON para handlers
func HandleDecodeError(w http.ResponseWriter, err error) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.Is(err, io.EOF):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "empty_body", "message": "el cuerpo de la solicitud es obligatorio"})
	case errors.As(err, &syntaxError):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_json", "message": "el cuerpo JSON no tiene un formato válido"})
	case errors.As(err, &unmarshalTypeError):
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_json_type", "message": "uno o más campos del JSON tienen un tipo inválido"})
	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_json", "message": "el cuerpo JSON es inválido"})
	}
}
