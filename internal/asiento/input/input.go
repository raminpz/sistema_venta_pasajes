package input

import (
	"sistema_venta_pasajes/internal/asiento/util"
)

type CreateAsientoInput struct {
	IDVehiculo    int    `json:"id_vehiculo" binding:"required"`
	NumeroAsiento string `json:"numero_asiento" binding:"required"`
	Estado        string `json:"estado" binding:"required"`
}

type UpdateAsientoInput struct {
	NumeroAsiento string `json:"numero_asiento" binding:"required"`
	Estado        string `json:"estado" binding:"required"`
}

type CambiarEstadoAsientoInput struct {
	Estado string `json:"estado" binding:"required"`
}

func (in *CambiarEstadoAsientoInput) Validate() error {
	return util.ValidateEstadoAsiento(in.Estado)
}
