package util

import (
	"regexp"
)

func ValidarPlaca(placa string) bool {
	// Ejemplo: ABC-123 o similar
	reg := regexp.MustCompile(`^[A-Z0-9-]{6,10}$`)
	return reg.MatchString(placa)
}
