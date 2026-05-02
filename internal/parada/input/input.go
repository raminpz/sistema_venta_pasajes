package input

type CreateParadaInput struct {
	IDRuta     int64 `json:"id_ruta"`
	IDTerminal int64 `json:"id_terminal"`
	Orden      int   `json:"orden"`
}

type UpdateParadaInput struct {
	IDParada   int64  `json:"id_parada"`
	IDTerminal *int64 `json:"id_terminal"`
	Orden      *int   `json:"orden"`
}

type ParadaOutput struct {
	IDParada   int64 `json:"id_parada"`
	IDRuta     int64 `json:"id_ruta"`
	IDTerminal int64 `json:"id_terminal"`
	Orden      int   `json:"orden"`
}
