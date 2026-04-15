package pkg

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var mysqlConstraintNameRegex = regexp.MustCompile("(?i)constraint `([^`]+)`")

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

	if dbErr := ParseDBError(err, "database_error", "Error interno de base de datos", nil, nil); dbErr != nil {
		var parsed *AppError
		if errors.As(dbErr, &parsed) {
			return parsed
		}
	}

	return Internal("ocurrió un error interno en el servidor").WithCause(err)
}

// ParseDBError transforma errores MySQL en AppError estandar con mensajes claros.
// fkMessages y duplicateMessages usan como clave el nombre del constraint/index en MySQL.
func ParseDBError(err error, errCode, genericMsg string, fkMessages map[string]string, duplicateMessages map[string]string) error {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	var mysqlErr *mysqlDriver.MySQLError
	if !errors.As(err, &mysqlErr) {
		return NewAppError(http.StatusInternalServerError, errCode, genericMsg).WithCause(err)
	}

	constraint := extractMySQLConstraint(mysqlErr.Message)
	constraintUpper := strings.ToUpper(constraint)

	switch mysqlErr.Number {
	case 1062:
		if msg := lookupConstraintMessage(duplicateMessages, constraintUpper); msg != "" {
			return Conflict("duplicate_resource", msg).WithCause(err)
		}
		return Conflict("duplicate_resource", "El registro ya existe en la base de datos").WithCause(err)
	case 1451, 1452:
		if msg := lookupConstraintMessage(fkMessages, constraintUpper); msg != "" {
			return Conflict("fk_conflict", msg).WithCause(err)
		}
		if fallbackMsg, ok := fkMessages["*"]; ok && fallbackMsg != "" {
			return Conflict("fk_conflict", fallbackMsg).WithCause(err)
		}
		return Conflict("fk_conflict", "No se pudo completar la operación por restricción de integridad referencial").WithCause(err)
	case 1265:
		return BadRequest(errCode, "Uno o más campos tienen un valor inválido para la base de datos").WithCause(err)
	default:
		return NewAppError(http.StatusInternalServerError, errCode, genericMsg).WithDetails(mysqlErr.Message).WithCause(err)
	}
}

func extractMySQLConstraint(message string) string {
	matches := mysqlConstraintNameRegex.FindStringSubmatch(message)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func lookupConstraintMessage(messages map[string]string, constraintUpper string) string {
	if len(messages) == 0 || constraintUpper == "" {
		return ""
	}
	if msg, ok := messages[constraintUpper]; ok {
		return msg
	}
	for key, msg := range messages {
		if strings.Contains(constraintUpper, strings.ToUpper(key)) {
			return msg
		}
	}
	return ""
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
