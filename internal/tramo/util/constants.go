package util

const (
	// Mensajes de éxito
	MSG_CREATED = "Tramo creado correctamente"
	MSG_UPDATED = "Tramo actualizado correctamente"
	MSG_DELETED = "Tramo eliminado correctamente"
	MSG_GET     = "Tramo obtenido correctamente"
	MSG_LIST    = "Lista de tramos"

	// Errores de validación de campos
	ERR_REQUIRED_RUTA           = "El ID de ruta es obligatorio"
	ERR_REQUIRED_PARADA_ORIGEN  = "La parada de origen es obligatoria"
	ERR_REQUIRED_PARADA_DESTINO = "La parada de destino es obligatoria"
	ERR_PARADAS_IGUALES         = "La parada de origen y destino no pueden ser iguales"
	ERR_EMPTY_UPDATE            = "Debe enviar al menos un campo para actualizar"
	ERR_DUPLICATE               = "Ya existe un tramo con esas paradas para esta ruta"

	// Errores genéricos
	ERR_INVALID_JSON = "El cuerpo JSON no tiene un formato válido"
	ERR_INVALID_ID   = "El ID proporcionado no es válido"
	ERR_NOT_FOUND    = "No se encontró el tramo"
	ERR_INVALID_PAGE = "Parámetros de paginación inválidos"
	ERR_DELETE       = "Error al eliminar el tramo"

	// Códigos de error
	ERR_CODE_NOT_FOUND = "not_found"
	ERR_CODE_DELETE    = "delete_error"
)
