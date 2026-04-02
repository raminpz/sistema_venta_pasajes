package input

type CreateProgramacionInput struct {
	IDRuta       int64   `json:"id_ruta"`
	IDVehiculo   int64   `json:"id_vehiculo"`
	IDConductor  int64   `json:"id_conductor"`
	FechaSalida  string  `json:"fecha_salida"`
	FechaLlegada *string `json:"fecha_llegada"`
	Estado       string  `json:"estado"`
}

type UpdateProgramacionInput struct {
	IDRuta       *int64  `json:"id_ruta"`
	IDVehiculo   *int64  `json:"id_vehiculo"`
	IDConductor  *int64  `json:"id_conductor"`
	FechaSalida  *string `json:"fecha_salida"`
	FechaLlegada *string `json:"fecha_llegada"`
	Estado       *string `json:"estado"`
}

type ProgramacionOutput struct {
	IDProgramacion int64   `json:"id_programacion"`
	IDRuta         int64   `json:"id_ruta"`
	IDVehiculo     int64   `json:"id_vehiculo"`
	IDConductor    int64   `json:"id_conductor"`
	FechaSalida    string  `json:"fecha_salida"`
	FechaLlegada   *string `json:"fecha_llegada"`
	Estado         string  `json:"estado"`
	CreatedAt      *string `json:"created_at"`
	UpdatedAt      *string `json:"updated_at"`
}
