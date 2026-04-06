package util

const (
	// Mensajes de éxito
	MSG_LOGIN_OK   = "Sesión iniciada correctamente"
	MSG_REFRESH_OK = "Token renovado correctamente"
	MSG_LOGOUT_OK  = "Sesión cerrada correctamente"

	// Mensajes de error
	MSG_CREDENCIALES_INVALIDAS  = "Credenciales inválidas. Verifique su email y contraseña."
	MSG_USUARIO_INACTIVO        = "El usuario está inactivo. Contacte al administrador."
	MSG_TOKEN_INVALIDO          = "Token inválido o expirado."
	MSG_TOKEN_REQUERIDO         = "Se requiere token de autenticación en el encabezado Authorization."
	MSG_REFRESH_TOKEN_INVALIDO  = "El token de refresco es inválido o ya fue revocado."
	MSG_ACCESO_DENEGADO         = "No tienes permisos para realizar esta acción."
	MSG_BODY_REQUERIDO          = "El cuerpo de la solicitud es obligatorio."
	MSG_EMAIL_REQUERIDO         = "El campo email es obligatorio."
	MSG_PASSWORD_REQUERIDO      = "El campo password es obligatorio."
	MSG_REFRESH_TOKEN_REQUERIDO = "El campo refresh_token es obligatorio."

	// Códigos de error
	ERR_CODE_CREDENCIALES   = "credenciales_invalidas"
	ERR_CODE_USUARIO_INACT  = "usuario_inactivo"
	ERR_CODE_TOKEN_INVALIDO = "token_invalido"
	ERR_CODE_TOKEN_REQUERID = "token_requerido"
	ERR_CODE_ACCESO_DENEG   = "acceso_denegado"
	ERR_CODE_BODY_INVALIDO  = "body_invalido"
	ERR_CODE_NO_AUTENTICADO = "no_autenticado"
)
