package input

type CreateTipoVehiculoInput struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

type UpdateTipoVehiculoInput struct {
	Nombre      *string `json:"nombre"`
	Descripcion *string `json:"descripcion"`
}

type TipoVehiculoOutput struct {
	IDTipoVehiculo int64  `json:"id_tipo_vehiculo"`
	Nombre         string `json:"nombre"`
	Descripcion    string `json:"descripcion"`
}

