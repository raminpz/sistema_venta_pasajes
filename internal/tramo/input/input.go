package input

type CreateTramoInput struct {
	IDRuta          int64 `json:"id_ruta"`
	IDParadaOrigen  int64 `json:"id_parada_origen"`
	IDParadaDestino int64 `json:"id_parada_destino"`
}

type UpdateTramoInput struct {
	IDTramo         int64  `json:"id_tramo"`
	IDRuta          *int64 `json:"id_ruta"`
	IDParadaOrigen  *int64 `json:"id_parada_origen"`
	IDParadaDestino *int64 `json:"id_parada_destino"`
}

type TramoOutput struct {
	IDTramo         int64 `json:"id_tramo"`
	IDRuta          int64 `json:"id_ruta"`
	IDParadaOrigen  int64 `json:"id_parada_origen"`
	IDParadaDestino int64 `json:"id_parada_destino"`
}
