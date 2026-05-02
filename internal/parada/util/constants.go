package util

const (
	MSG_CREATED = "Parada creada correctamente"
	MSG_UPDATED = "Parada actualizada correctamente"
	MSG_DELETED = "Parada eliminada correctamente"
	MSG_GET     = "Parada obtenida correctamente"
	MSG_LIST    = "Lista de paradas"

	ERR_REQUIRED_RUTA      = "El ID de ruta es obligatorio"
	ERR_REQUIRED_TERMINAL  = "El ID de terminal es obligatorio"
	ERR_REQUIRED_ORDEN     = "El orden es obligatorio y debe ser mayor a cero"
	ERR_EMPTY_UPDATE       = "Debe enviar al menos un campo para actualizar"
	ERR_DUPLICATE          = "Ya existe una parada con ese terminal en esta ruta"
	ERR_DUPLICATE_ORDEN    = "Ya existe una parada con ese orden en esta ruta"

	ERR_INVALID_JSON   = "El cuerpo JSON no tiene un formato válido"
	ERR_INVALID_ID     = "El ID proporcionado no es válido"
	ERR_NOT_FOUND      = "No se encontró la parada"
	ERR_INVALID_PAGE   = "Parámetros de paginación inválidos"
	ERR_DELETE         = "Error al eliminar la parada"

	ERR_CODE_NOT_FOUND = "not_found"
	ERR_CODE_DELETE    = "delete_error"
)

