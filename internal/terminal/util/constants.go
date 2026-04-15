package util

const (
	MSG_TERMINAL_CREATED = "Terminal creada correctamente"
	MSG_TERMINAL_UPDATED = "Terminal actualizada correctamente"
	MSG_TERMINAL_DELETED = "Terminal eliminada correctamente"
	MSG_TERMINAL_GET     = "Terminal obtenida correctamente"
	MSG_TERMINAL_LIST    = "Terminales obtenidas correctamente"

	MSG_TERMINAL_CREATE_ERROR    = "Error al crear terminal"
	MSG_TERMINAL_DUPLICATE       = "Ya existe una terminal con el mismo nombre, ciudad y departamento"
	MSG_TERMINAL_VALIDATION      = "Existen errores de validacion"
	MSG_TERMINAL_INVALID_ID      = "El id de la terminal no es valido"
	MSG_TERMINAL_SINGLE_JSON_OBJ = "El cuerpo JSON debe contener un unico objeto"
	MSG_TERMINAL_NOT_FOUND       = "Terminal no encontrada"
	MSG_TERMINAL_DELETE_ERROR    = "Error al eliminar terminal"

	ERR_CODE_DUPLICATE_RESOURCE = "duplicate_resource"
	ERR_CODE_INVALID_ID         = "invalid_terminal_id"
	ERR_CODE_INVALID_JSON       = "invalid_json"
	ERR_CODE_EMPTY_BODY         = "empty_body"
	ERR_CODE_INVALID_JSON_TYPE  = "invalid_json_type"
	ERR_CODE_NOT_FOUND          = "not_found"
)
