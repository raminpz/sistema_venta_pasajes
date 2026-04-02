package input

type CreatePagoInput struct {
	IDVenta  int64   `json:"id_venta"`
	IDMetodo int64   `json:"id_metodo"`
	Monto    float64 `json:"monto"`
	Estado   string  `json:"estado"`
}

type UpdatePagoInput struct {
	IDMetodo *int64   `json:"id_metodo"`
	Monto    *float64 `json:"monto"`
	Estado   *string  `json:"estado"`
}

type PagoOutput struct {
	IDPago    int64   `json:"id_pago"`
	IDVenta   int64   `json:"id_venta"`
	IDMetodo  int64   `json:"id_metodo"`
	Monto     float64 `json:"monto"`
	Estado    string  `json:"estado"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}
