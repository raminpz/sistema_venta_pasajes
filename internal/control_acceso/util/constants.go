package util

const (
	// DIAS_ALERTA son los días previos al vencimiento en que el frontend muestra alerta de renovación
	DIAS_ALERTA = 30

	// DIAS_GRACIA son los días posteriores al vencimiento en que el sistema opera en modo SOLO_LECTURA
	DIAS_GRACIA = 30

	// PROVEEDOR_TELEFONO es el número de contacto del proveedor para renovaciones
	PROVEEDOR_TELEFONO = "961501468"

	// Mensajes de éxito
	MSG_ACCESO_STATUS    = "Estado del sistema obtenido correctamente."
	MSG_ACCESO_DETALLES  = "Detalles del control de acceso obtenidos correctamente."
	MSG_ACCESO_CREADO    = "Control de acceso creado y activado correctamente."
	MSG_ACCESO_ACTIVADO  = "Sistema reactivado correctamente."
	MSG_ACCESO_BLOQUEADO = "Sistema bloqueado correctamente."
	MSG_ACCESO_RENOVADO  = "Control de acceso renovado correctamente."

	// Mensajes de error
	ERR_ACCESO_FECHA_REQUERIDA  = "Las fechas de activación y expiración son obligatorias."
	ERR_ACCESO_FECHA_FORMATO    = "El formato de fecha es inválido. Use YYYY-MM-DD."
	ERR_ACCESO_FECHA_EXPIRACION = "La fecha de expiración debe ser igual o posterior a la de activación."
	ERR_ACCESO_NO_ENCONTRADO    = "El registro de control de acceso no existe."
	ERR_ACCESO_EXP_REQUERIDA    = "La nueva fecha de expiración es obligatoria."

	// Códigos de error
	CODE_ACCESO_NO_ENCONTRADO  = "acceso_no_encontrado"
	CODE_ACCESO_FECHA_INVALIDA = "fecha_invalida"
	CODE_ACCESO_FECHA_VALID    = "fecha_validacion_error"
)
