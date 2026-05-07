package input

type AsientoTramoOutput struct {
	IDAsientoTramo int64  `json:"id_asiento_tramo"`
	IDVenta        *int64 `json:"id_venta,omitempty"`
	IDAsiento      int64  `json:"id_asiento"`
	IDTramo        int64  `json:"id_tramo"`
	Estado         string `json:"estado"`
}

type CreateAsientoTramoInput struct {
	IDVenta   *int64
	IDAsiento int64
	IDTramo   int64
	Estado    string
}
