package input

import (
	"sistema_venta_pasajes/internal/conductor/util"
	"testing"
)

func TestIsValidNumeroLicenciaAlfanumerico(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"ABC123456", true},
		{"123456789", true},
		{"abcdefghi", true},
		{"abc12345", false},   // menos de 9
		{"abc1234567", false}, // más de 9
		{"abc12345!", false},  // carácter especial
		{"", false},
		{"S44106817", true}, // caso explícito alfanumérico
	}
	for _, c := range cases {
		if got := util.IsValidNumeroLicencia(c.input); got != c.expected {
			t.Errorf("input: %q, expected %v, got %v", c.input, c.expected, got)
		}
	}
}
