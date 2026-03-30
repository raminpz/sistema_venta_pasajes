package input

type CreateInput struct {
	ClaveLicencia   string `json:"clave_licencia"`
	FechaActivacion string `json:"fecha_activacion"`
	FechaExpiracion string `json:"fecha_expiracion"`
	Estado          string `json:"estado"`
}

type UpdateInput struct {
	FechaActivacion string `json:"fecha_activacion"`
	FechaExpiracion string `json:"fecha_expiracion"`
	Estado          string `json:"estado"`
}

