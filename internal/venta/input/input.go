package input

type VentaCreateInput struct {
	IDUsuario         int64   `json:"id_usuario"`
	IDTipoComprobante int64   `json:"id_tipo_comprobante"`
	Nota              string  `json:"nota"`
	Observaciones     string  `json:"observaciones"`
	Subtotal          float64 `json:"subtotal"`
}

type VentaUpdateInput struct {
	Nota          string `json:"nota"`
	Observaciones string `json:"observaciones"`
	Estado        string `json:"estado"`
}

type VentaOutput struct {
	IDVenta           int64   `json:"id_venta"`
	IDUsuario         int64   `json:"id_usuario"`
	IDTipoComprobante int64   `json:"id_tipo_comprobante"`
	Serie             string  `json:"serie"`
	Correlativo       uint    `json:"correlativo"`
	NumeroComprobante string  `json:"numero_comprobante"`
	Nota              string  `json:"nota"`
	Observaciones     string  `json:"observaciones"`
	Subtotal          float64 `json:"subtotal"`
	IGV               float64 `json:"igv"`
	Total             float64 `json:"total"`
	QRCode            string  `json:"qr_code"`
	Estado            string  `json:"estado"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}
