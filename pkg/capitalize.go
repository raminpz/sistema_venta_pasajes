package pkg

import (
	"strings"
)

// Capitaliza cada palabra: "ramiro nuñez perez" -> "Ramiro Nuñez Perez"
func CapitalizeWords(s string) string {
	words := strings.Fields(strings.ToLower(s))
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

