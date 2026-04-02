package util

const (
	MsgVentaCreada               = "Venta registrada correctamente"
	MsgVentaActualizada          = "Venta actualizada correctamente"
	MsgVentaEliminada            = "Venta eliminada correctamente"
	MsgVentaListada              = "Lista de ventas"
	MsgVentaObtener              = "Venta obtenida correctamente"
	MsgVentaNoEncontrada         = "Venta no encontrada"
	MsgVentaErrorCrear           = "Error al registrar venta"
	MsgVentaErrorActualizar      = "Error al actualizar venta"
	MsgVentaErrorEliminar        = "Error al eliminar venta"
	MsgVentaErrorListar          = "Error al listar ventas"
	MsgVentaErrorValidacion      = "Datos de venta inválidos"
	MsgVentaErrorTipoComprob     = "Tipo de comprobante inválido. Use: 1=BOLETA, 2=FACTURA, 3=TICKET"
	MsgVentaErrorSerie           = "No se pudo determinar la serie del comprobante"
	MsgVentaErrorCorrelativo     = "No se pudo obtener el correlativo automático"
	MsgVentaSubtotalRequerido    = "El subtotal debe ser mayor a cero"
	MsgVentaUsuarioRequerido     = "El usuario es obligatorio"
	MsgVentaComprobanteRequerido = "El tipo de comprobante es obligatorio"
	MsgVentaErrorQR              = "Error al generar el código QR"
	MsgVentaPaginacionInvalida   = "Parámetros de paginación inválidos"

	VentaEstadoRegistrada = "REGISTRADA"

	ErrCodeInvalidBody = "invalid_body"
	ErrCodeInvalidID   = "invalid_id"
	ErrCodeInvalidPage = "invalid_pagination"
	ErrCodeCreateError = "create_error"
	ErrCodeUpdateError = "update_error"
	ErrCodeDeleteError = "delete_error"
	ErrCodeListError   = "list_error"
	ErrCodeNotFound    = "not_found"
)
