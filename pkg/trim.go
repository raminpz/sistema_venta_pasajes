package pkg

import (
	"reflect"
	"strings"
)

// TrimSpacesOnStruct recorre todos los campos string de un struct (por puntero) y les aplica TrimSpace.
// Solo afecta campos string y *string.
func TrimSpacesOnStruct(ptr interface{}) {
	// Usamos reflexión para recorrer los campos
	// Si el campo es string o *string, aplicamos strings.TrimSpace
	// Si es *string y es nil, no hace nada
	// Si es *string y no es nil, lo actualiza
	// Si es string, lo actualiza directamente
	// Si el campo es un struct anidado, no lo recorre (solo primer nivel)
	   // (import redundante eliminado)

	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
			   field := v.Field(i)
			   if !field.CanSet() {
					   continue
			   }
			   if field.Kind() == reflect.String {
					   field.SetString(strings.TrimSpace(field.String()))
			   } else if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.String {
					   if !field.IsNil() {
							   str := field.Elem().String()
							   field.Elem().SetString(strings.TrimSpace(str))
					   }
			   }
	}
}

