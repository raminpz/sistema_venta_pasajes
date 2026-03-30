package input

type CreateRutaInput struct {
	IDOrigenTerminal  int     `json:"id_origen_terminal" binding:"required"`
	IDDestinoTerminal int     `json:"id_destino_terminal" binding:"required"`
	DuracionHoras     float64 `json:"duracion_horas" binding:"required"`
}

type UpdateRutaInput struct {
	IDOrigenTerminal  *int     `json:"id_origen_terminal"`
	IDDestinoTerminal *int     `json:"id_destino_terminal"`
	DuracionHoras     *float64 `json:"duracion_horas"`
}

type RutaOutput struct {
	IDRuta            int     `json:"id_ruta"`
	IDOrigenTerminal  int     `json:"id_origen_terminal"`
	IDDestinoTerminal int     `json:"id_destino_terminal"`
	DuracionHoras     float64 `json:"duracion_horas"`
}

