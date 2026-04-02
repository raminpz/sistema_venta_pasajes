package util

// Reusable error messages for passenger and other modules
const (
	ERR_PHONE_FORMAT   = "El teléfono debe tener exactamente 9 dígitos numéricos"
	ERR_EMAIL_FORMAT   = "El email no tiene un formato válido"
	ERR_DATE_FORMAT    = "La fecha debe tener formato YYYY-MM-DD"
	ERR_REQUIRED_FIELD = "El campo es requerido"
)

// Mensajes para handler
const (
	MSG_JSON_INVALID   = "JSON inválido"
	MSG_CREATE_ERROR   = "Error al crear pasajero"
	MSG_MISSING_ID     = "Falta id en la URL"
	MSG_INVALID_ID     = "id inválido"
	MSG_UPDATE_ERROR   = "Error al actualizar pasajero"
	MSG_DELETE_ERROR   = "Error al eliminar pasajero"
	MSG_DELETE_SUCCESS = "Pasajero eliminado correctamente"
	MSG_NOT_FOUND      = "Pasajero no encontrado"
	MSG_LIST_ERROR     = "Error al listar pasajeros"
	MSG_LIST_SUCCESS   = "Lista de pasajeros"
	MSG_CREATE_SUCCESS = "Pasajero creado"
	MSG_UPDATE_SUCCESS = "Pasajero actualizado"
	MSG_FOUND_SUCCESS  = "Pasajero encontrado"
	MSG_MISSING_QUERY  = "El parámetro de búsqueda 'q' es obligatorio"
	MSG_SEARCH_ERROR   = "Error al buscar pasajeros"
	MSG_DUPLICATE_DNI  = "El DNI ingresado ya fue registrado."
	MSG_SEARCH_SUCCESS = "Búsqueda de pasajeros exitosa"
	MSG_INVALID_PAGE   = "Parámetros de paginación inválidos"
)
