package input

import "time"

// GenerarLiquidacionInput es el input para generar una liquidación de viaje.
type GenerarLiquidacionInput struct {
	IDProgramacion int64  `json:"id_programacion"`
	Observaciones  string `json:"observaciones"`
}

// ActualizarEstadoInput permite cambiar el estado de la liquidación.
type ActualizarEstadoInput struct {
	Estado        string `json:"estado"`
	Observaciones string `json:"observaciones"`
}

// LiquidacionOutput es la respuesta de una liquidación persistida.
type LiquidacionOutput struct {
	IDLiquidacion       int64      `json:"id_liquidacion"`
	IDProgramacion      int64      `json:"id_programacion"`
	IDConductor         int64      `json:"id_conductor"`
	TotalPasajes        float64    `json:"total_pasajes"`
	TotalEncomiendas    float64    `json:"total_encomiendas"`
	TotalCaja           float64    `json:"total_caja"`
	Estado              string     `json:"estado"`
	FechaLiquidacion    *time.Time `json:"fecha_liquidacion"`
	Observaciones       string     `json:"observaciones"`
	CantidadPasajes     int        `json:"cantidad_pasajes"`
	CantidadEncomiendas int        `json:"cantidad_encomiendas"`
}

// ResumenCajaOutput es la previsualización del total de caja sin persistir.
type ResumenCajaOutput struct {
	IDProgramacion      int64   `json:"id_programacion"`
	TotalPasajes        float64 `json:"total_pasajes"`
	TotalEncomiendas    float64 `json:"total_encomiendas"`
	TotalCaja           float64 `json:"total_caja"`
	CantidadPasajes     int     `json:"cantidad_pasajes"`
	CantidadEncomiendas int     `json:"cantidad_encomiendas"`
}
