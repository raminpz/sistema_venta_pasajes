package input

type CreateParadaInput struct {
	IDRuta       int64  `json:"id_ruta"`
	NombreParada string `json:"nombre_parada"`
	Orden        int    `json:"orden"`
}

type UpdateParadaInput struct {
	IDParada     int64   `json:"id_parada"`
	NombreParada *string `json:"nombre_parada"`
	Orden        *int    `json:"orden"`
}

type ParadaOutput struct {
	IDParada     int64  `json:"id_parada"`
	IDRuta       int64  `json:"id_ruta"`
	NombreParada string `json:"nombre_parada"`
	Orden        int    `json:"orden"`
}
