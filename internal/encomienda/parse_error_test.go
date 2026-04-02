package encomienda

import (
	"testing"

	"sistema_venta_pasajes/internal/encomienda/util"
	"sistema_venta_pasajes/pkg"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

// Test: ParseDBError detects FK violations
func TestParseDBErrorFKVenta(t *testing.T) {
	// Simular error MySQL 1452 (FK violation) para ID_VENTA
	mysqlErr := &mysqlDriver.MySQLError{
		Number:  1452,
		Message: "Error 1452: Cannot add or update a child row: a foreign key constraint fails (`test`.`ENCOMIENDA`, CONSTRAINT `FK_ENCOMIENDA_VENTA` FOREIGN KEY (`ID_VENTA`) REFERENCES `VENTA` (`ID_VENTA`))",
	}

	err := util.ParseDBError(mysqlErr, util.ERR_CODE_CREATE, "Error genérico")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	appErr, ok := err.(*pkg.AppError)
	if !ok {
		t.Fatalf("Expected *pkg.AppError, got %T", err)
	}

	if appErr.Message != util.MSG_ENCOMIENDA_VENTA_NOT_FOUND {
		t.Errorf("Expected message '%s', got '%s'", util.MSG_ENCOMIENDA_VENTA_NOT_FOUND, appErr.Message)
	}

	// Debe retornar 409 Conflict por FK, no 400
	if appErr.Status != 409 {
		t.Errorf("Expected status 409, got %d", appErr.Status)
	}
}

// Test: ParseDBError detects FK violations for Programacion
func TestParseDBErrorFKProgramacion(t *testing.T) {
	mysqlErr := &mysqlDriver.MySQLError{
		Number:  1452,
		Message: "Error 1452: Cannot add or update a child row: a foreign key constraint fails (`test`.`ENCOMIENDA`, CONSTRAINT `FK_ENCOMIENDA_PROGRAMACION` FOREIGN KEY (`ID_PROGRAMACION`) REFERENCES `PROGRAMACION` (`ID_PROGRAMACION`))",
	}

	err := util.ParseDBError(mysqlErr, util.ERR_CODE_CREATE, "Error genérico")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	appErr, ok := err.(*pkg.AppError)
	if !ok {
		t.Fatalf("Expected *pkg.AppError, got %T", err)
	}

	if appErr.Message != util.MSG_ENCOMIENDA_PROG_NOT_FOUND {
		t.Errorf("Expected message '%s', got '%s'", util.MSG_ENCOMIENDA_PROG_NOT_FOUND, appErr.Message)
	}

	if appErr.Status != 409 {
		t.Errorf("Expected status 409, got %d", appErr.Status)
	}
}

// Test: ParseDBError detects unknown FK violations
func TestParseDBErrorUnknownFK(t *testing.T) {
	mysqlErr := &mysqlDriver.MySQLError{
		Number:  1452,
		Message: "Error 1452: Cannot add or update a child row: a foreign key constraint fails",
	}

	err := util.ParseDBError(mysqlErr, util.ERR_CODE_CREATE, "Error genérico")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	appErr, ok := err.(*pkg.AppError)
	if !ok {
		t.Fatalf("Expected *pkg.AppError, got %T", err)
	}

	if appErr.Message != util.MSG_ENCOMIENDA_FOREIGN_KEY_ERROR {
		t.Errorf("Expected message '%s', got '%s'", util.MSG_ENCOMIENDA_FOREIGN_KEY_ERROR, appErr.Message)
	}

	if appErr.Status != 409 {
		t.Errorf("Expected status 409, got %d", appErr.Status)
	}
}

// Test: ParseDBError detects 1451 FK violations (delete/update parent)
func TestParseDBErrorFK1451(t *testing.T) {
	mysqlErr := &mysqlDriver.MySQLError{
		Number:  1451,
		Message: "Error 1451: Cannot delete or update a parent row",
	}

	err := util.ParseDBError(mysqlErr, util.ERR_CODE_DELETE, "Error genérico")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	appErr, ok := err.(*pkg.AppError)
	if !ok {
		t.Fatalf("Expected *pkg.AppError, got %T", err)
	}

	if appErr.Status != 409 {
		t.Errorf("Expected status 409, got %d", appErr.Status)
	}
}

// Test: Validation error tiene estructura correcta
func TestValidationErrorStructure(t *testing.T) {
	details := map[string]string{
		"id_venta": util.MSG_ENCOMIENDA_ID_VENTA_REQUIRED,
		"costo":    util.MSG_ENCOMIENDA_COST_REQUIRED,
	}

	appErr := pkg.Validation(util.MSG_ENCOMIENDA_VALIDATION, details)
	if appErr == nil {
		t.Fatal("Expected error, got nil")
	}

	if appErr.Status != 422 {
		t.Errorf("Expected status 422, got %d", appErr.Status)
	}

	if appErr.Code != "validation_error" {
		t.Errorf("Expected code 'validation_error', got '%s'", appErr.Code)
	}

	if appErr.Details == nil {
		t.Fatal("Expected details, got nil")
	}
}
