package util

import "strings"

// ValidarLoginInput valida los campos obligatorios del login.
func ValidarLoginInput(email, password string) (string, bool) {
	if strings.TrimSpace(email) == "" {
		return MSG_EMAIL_REQUERIDO, false
	}
	if strings.TrimSpace(password) == "" {
		return MSG_PASSWORD_REQUERIDO, false
	}
	return "", true
}
