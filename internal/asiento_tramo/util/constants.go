package util

// Error codes
const (
	ERR_CODE_ASIENTO_TRAMO_CREATE    = "ASIENTO_TRAMO_001"
	ERR_CODE_ASIENTO_TRAMO_UPDATE    = "ASIENTO_TRAMO_002"
	ERR_CODE_ASIENTO_TRAMO_DELETE    = "ASIENTO_TRAMO_003"
	ERR_CODE_ASIENTO_TRAMO_NOT_FOUND = "ASIENTO_TRAMO_004"
	ERR_CODE_ASIENTO_TRAMO_DUPLICATE = "ASIENTO_TRAMO_005"
)

// Messages
const (
	MSG_ASIENTO_TRAMO_CREATE_ERROR = "Error al registrar asiento-tramo"
	MSG_ASIENTO_TRAMO_UPDATE_ERROR = "Error al actualizar asiento-tramo"
	MSG_ASIENTO_TRAMO_DELETE_ERROR = "Error al eliminar asiento-tramo"
	MSG_ASIENTO_TRAMO_NOT_FOUND    = "El asiento-tramo no fue encontrado"
	MSG_ASIENTO_TRAMO_DUPLICATE    = "El asiento ya está registrado en este tramo"
)

// Estados
const (
	ESTADO_DISPONIBLE = "DISPONIBLE"
	ESTADO_OCUPADO    = "OCUPADO"
)
