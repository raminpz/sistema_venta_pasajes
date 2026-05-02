package util

const (
	ERR_INVALID_JSON    = "El cuerpo JSON no tiene un formato valido"
	ERR_INVALID_ID      = "El ID proporcionado no es valido"
	MSG_DELETED         = "Vehiculo eliminado correctamente"
	ERR_NOT_FOUND       = "No se encontro el vehiculo"
	ERR_DUPLICATE_PLATE = "El numero de placa ya esta registrado"
	MSG_CREATED         = "Vehiculo creado correctamente"
	MSG_UPDATED         = "Vehiculo actualizado correctamente"
	MSG_GET             = "Vehiculo obtenido correctamente"
	MSG_LIST            = "Lista de vehiculos"
	ERR_INVALID_BODY    = "El cuerpo de la solicitud es invalido"
	ERR_INVALID_PAGE    = "Parametros de paginacion invalidos"
	ERR_DELETE          = "Error al eliminar vehiculo"

	ERR_REQUIRED_TIPO_VEHICULO    = "El tipo de vehiculo es obligatorio"
	ERR_REQUIRED_PLACA            = "El numero de placa es obligatorio"
	ERR_REQUIRED_MARCA            = "La marca es obligatoria"
	ERR_REQUIRED_MODELO           = "El modelo es obligatorio"
	ERR_REQUIRED_ANIO             = "El año de fabricacion es obligatorio"
	ERR_REQUIRED_CHASIS           = "El numero de chasis es obligatorio"
	ERR_REQUIRED_CAPACIDAD        = "La capacidad es obligatoria y debe ser mayor a cero"
	ERR_REQUIRED_SOAT             = "El numero de SOAT es obligatorio"
	ERR_REQUIRED_FECHA_VENC_SOAT  = "La fecha de vencimiento del SOAT es obligatoria"
	ERR_REQUIRED_REVISION_TECNICA = "El numero de revision tecnica es obligatorio"
	ERR_REQUIRED_FECHA_VENC_REV   = "La fecha de vencimiento de revision tecnica es obligatoria"
	ERR_REQUIRED_ESTADO           = "El estado es obligatorio"
	ERR_INVALID_ESTADO            = "El estado debe ser ACTIVO o INACTIVO"
	ERR_EMPTY_UPDATE              = "Debe enviar al menos un campo para actualizar"

	ERR_CODE_NOT_FOUND = "not_found"
	ERR_CODE_DELETE    = "delete_error"
)
