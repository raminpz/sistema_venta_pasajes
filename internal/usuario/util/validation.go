package util

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
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
	page := 1
	size := 15
	if p := r.URL.Query().Get("page"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 {
			return 0, 0, errors.New("Parámetro 'page' inválido")
		}
		page = v
	}
	if s := r.URL.Query().Get("size"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 {
			return 0, 0, errors.New("Parámetro 'size' inválido")
		}
		size = v
	}
	return page, size, nil
}
