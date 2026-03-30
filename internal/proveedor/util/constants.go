package util

const (
	MSG_PROVIDER_ID_GT_ZERO     = "el id del proveedor del sistema debe ser mayor a 0"
	MSG_RUC_FORMAT              = "el RUC debe contener exactamente 11 dígitos numéricos"
	MSG_RAZON_SOCIAL_REQUIRED   = "la razón social es obligatoria"
	MSG_RAZON_SOCIAL_LENGTH     = "la razón social no puede exceder 150 caracteres"
	MSG_NOMBRE_COMERCIAL_LENGTH = "el nombre comercial no puede exceder 150 caracteres"
	MSG_DIRECCION_LENGTH        = "la dirección no puede exceder 200 caracteres"
	MSG_TELEFONO_LENGTH         = "el teléfono no puede exceder 20 caracteres"
	MSG_EMAIL_LENGTH            = "el email no puede exceder 100 caracteres"
	MSG_EMAIL_FORMAT            = "el email no tiene un formato válido"
	MSG_WEB_LENGTH              = "la web no puede exceder 150 caracteres"
	MSG_WEB_URL                 = "la web debe tener una URL válida"
	MSG_WEB_SCHEMA_HOST         = "la web debe incluir esquema y host válidos"
	MSG_WEB_HTTP_HTTPS          = "la web debe usar http o https"
	MSG_VALIDATION              = "los datos enviados no son válidos"
	MSG_PROVIDER_NOT_FOUND      = "proveedor del sistema no encontrado"
	MSG_UPDATE_PROVIDER_ERROR   = "actualizar proveedor del sistema: %w"
	MSG_DELETE_PROVIDER_ERROR   = "eliminar proveedor del sistema: %w"
	ERR_RUC_DUPLICADO           = "El RUC ya está registrado"
	ERR_EMAIL_DUPLICADO         = "El email ya está registrado"
)
