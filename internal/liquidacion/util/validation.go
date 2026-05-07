package util
// EsEstadoValido verifica que el estado sea PENDIENTE o ENTREGADO.
func EsEstadoValido(estado string) bool {
return estado == "PENDIENTE" || estado == "ENTREGADO"
}
