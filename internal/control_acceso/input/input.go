package input

// ActivarControlAccesoInput es el input para crear un nuevo registro de control de acceso.
// Si clave_licencia se omite, el sistema genera una automáticamente.
type ActivarControlAccesoInput struct {
	FechaActivacion string `json:"fecha_activacion"` // YYYY-MM-DD
	FechaExpiracion string `json:"fecha_expiracion"` // YYYY-MM-DD
}

// RenovarControlAccesoInput es el input para renovar el control de acceso.
// Si clave_licencia se omite, el sistema genera una nueva automáticamente.
type RenovarControlAccesoInput struct {
	FechaExpiracion string `json:"fecha_expiracion"` // YYYY-MM-DD
}

// ControlAccesoOutput es la respuesta completa del control de acceso (solo para proveedor).
type ControlAccesoOutput struct {
	IDAcceso        int64  `json:"id_acceso"`
	FechaActivacion string `json:"fecha_activacion"`
	FechaExpiracion string `json:"fecha_expiracion"`
	EstadoDB        string `json:"estado_db"`
	EstadoEfectivo  string `json:"estado_efectivo"`
	DiasParaVencer  int    `json:"dias_para_vencer"`
	EnAlerta        bool   `json:"en_alerta"` // true = faltan ≤30 días, mostrar aviso
	EnGracia        bool   `json:"en_gracia"` // true = dentro del período de gracia (solo lectura)
}

// ControlAccesoStatusOutput es la respuesta pública del estado del sistema.
// El frontend consulta este endpoint para decidir qué mostrar al usuario:
//   - en_alerta=true  → banner de aviso de renovación (sistema aún operativo)
//   - SOLO_LECTURA    → formularios deshabilitados, solo consultas
//   - BLOQUEADO       → solo login visible, dashboard con mensaje de bloqueo
type ControlAccesoStatusOutput struct {
	EstadoEfectivo string `json:"estado_efectivo"`
	Mensaje        string `json:"mensaje"`
	DiasParaVencer int    `json:"dias_para_vencer"`
	EnAlerta       bool   `json:"en_alerta"` // true = próximo a vencer (≤30 días), sistema operativo
	EnGracia       bool   `json:"en_gracia"` // true = vencida pero en período de gracia (solo lectura)
}
