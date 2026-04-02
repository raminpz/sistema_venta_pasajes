package util

import (
	"net/http"
	"regexp"
	"sistema_venta_pasajes/pkg"
	"strings"
)

func ValidarEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func ValidarDNI(dni string) bool {
	re := regexp.MustCompile(`^[0-9]{8}$`)
	return re.MatchString(dni)
}

func ValidarCamposObligatorios(nombre, apellidos, email, password, telefono string) bool {
	return strings.TrimSpace(nombre) != "" && strings.TrimSpace(apellidos) != "" &&
		strings.TrimSpace(email) != "" && strings.TrimSpace(password) != "" && strings.TrimSpace(telefono) != ""
}

// ParsePaginationParams extrae y valida los parámetros de paginación desde la request
func ParsePaginationParams(r *http.Request) (int, int, error) {
	return pkg.ParsePaginationParams(r, 1, 15)
}
