package util

import (
	"regexp"
)

// IsValidNumeroLicencia valida que el número de licencia tenga exactamente 9 caracteres alfanuméricos (letras y números)
func IsValidNumeroLicencia(numero string) bool {
	   matched, _ := regexp.MatchString(`^[A-Za-z0-9]{9}$`, numero)
	   return matched
}
