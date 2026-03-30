package input

import (
	"fmt"
	"time"
)

// DateOnly permite parsear fechas en formato YYYY-MM-DD
type DateOnly struct {
	time.Time
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // quitar comillas
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("la fecha debe tener formato YYYY-MM-DD: %w", err)
	}
	d.Time = t
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte("\"" + d.Time.Format("2006-01-02") + "\""), nil
}

type CreateVehiculoInput struct {
	IDTipoVehiculo       int64     `json:"id_tipo_vehiculo"`
	NroPlaca             string    `json:"nro_placa"`
	Marca                string    `json:"marca"`
	Modelo               string    `json:"modelo"`
	AnioFabricacion      int       `json:"anio_fabricacion"`
	NumeroChasis         string    `json:"numero_chasis"`
	Capacidad            int       `json:"capacidad"`
	NroSoat              string    `json:"nro_soat"`
	FechaVencSoat        *DateOnly `json:"fecha_venc_soat"`
	NroRevisionTecnica   string    `json:"nro_revision_tecnica"`
	FechaVencRevisionTec *DateOnly `json:"fecha_venc_revision_tecnica"`
	Estado               string    `json:"estado"`
}

type UpdateVehiculoInput struct {
	IDVehiculo           int64     `json:"id_vehiculo"`
	IDTipoVehiculo       int64     `json:"id_tipo_vehiculo"`
	NroPlaca             string    `json:"nro_placa"`
	Marca                string    `json:"marca"`
	Modelo               string    `json:"modelo"`
	AnioFabricacion      int       `json:"anio_fabricacion"`
	NumeroChasis         string    `json:"numero_chasis"`
	Capacidad            int       `json:"capacidad"`
	NroSoat              string    `json:"nro_soat"`
	FechaVencSoat        *DateOnly `json:"fecha_venc_soat"`
	NroRevisionTecnica   string    `json:"nro_revision_tecnica"`
	FechaVencRevisionTec *DateOnly `json:"fecha_venc_revision_tecnica"`
	Estado               string    `json:"estado"`
}

type VehiculoOutput struct {
	IDVehiculo           int64     `json:"id_vehiculo"`
	IDTipoVehiculo       int64     `json:"id_tipo_vehiculo"`
	NroPlaca             string    `json:"nro_placa"`
	Marca                string    `json:"marca"`
	Modelo               string    `json:"modelo"`
	AnioFabricacion      int       `json:"anio_fabricacion"`
	NumeroChasis         string    `json:"numero_chasis"`
	Capacidad            int       `json:"capacidad"`
	NroSoat              string    `json:"nro_soat"`
	FechaVencSoat        *DateOnly `json:"fecha_venc_soat"`
	NroRevisionTecnica   string    `json:"nro_revision_tecnica"`
	FechaVencRevisionTec *DateOnly `json:"fecha_venc_revision_tecnica"`
	Estado               string    `json:"estado"`
}
