package util

const (
	MSG_CREATED = "Tipo de vehiculo creado correctamente"
	MSG_UPDATED = "Tipo de vehiculo actualizado correctamente"
	MSG_DELETED = "Tipo de vehiculo eliminado correctamente"
	MSG_GET     = "Tipo de vehiculo obtenido correctamente"
	MSG_LIST    = "Tipos de vehiculo obtenidos correctamente"
)

const (
	ERR_INVALID_ID      = "El ID es invalido"
	ERR_REQUIRED_NAME   = "El nombre es obligatorio"
	ERR_REQUIRED_DESC   = "La descripcion es obligatoria"
	ERR_EMPTY_UPDATE    = "Debe enviar al menos un campo para actualizar"
	ERR_NOT_FOUND       = "No se encontro el tipo de vehiculo"
	ERR_CREATE          = "No se pudo crear el tipo de vehiculo"
	ERR_UPDATE          = "No se pudo actualizar el tipo de vehiculo"
	ERR_DELETE          = "No se pudo eliminar el tipo de vehiculo"
	ERR_LIST            = "No se pudo listar los tipos de vehiculo"
	ERR_DUPLICATE_NAME  = "El tipo de vehiculo ya existe"
	ERR_DELETE_CONFLICT = "No se puede eliminar: hay vehiculos asociados a este tipo"
)

const (
	ERR_CODE_NOT_FOUND = "tipo_vehiculo_not_found"
	ERR_CODE_INVALIDID = "invalid_id"
	ERR_CODE_CREATE    = "tipo_vehiculo_create_error"
	ERR_CODE_UPDATE    = "tipo_vehiculo_update_error"
	ERR_CODE_DELETE    = "tipo_vehiculo_delete_error"
	ERR_CODE_LIST      = "tipo_vehiculo_list_error"
)

