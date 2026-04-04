package input

type VentaCreateInput struct {
	IDUsuario         int64    `json:"id_usuario"`
	IDTipoComprobante int64    `json:"id_tipo_comprobante"`
	IDProgramacion    int64    `json:"id_programacion"`
	IDPasajero        int64    `json:"id_pasajero"`
	IDAsiento         int64    `json:"id_asiento"`
	Precio            float64  `json:"precio"`
	Descuento         *float64 `json:"descuento"`
	Nota              string   `json:"nota"`
	Observaciones     string   `json:"observaciones"`
}

type VentaUpdateInput struct {
	Nota          string `json:"nota"`
	Observaciones string `json:"observaciones"`
}

type VentaOutput struct {
	IDVenta           int64    `json:"id_venta"`
	IDUsuario         int64    `json:"id_usuario"`
	IDTipoComprobante int64    `json:"id_tipo_comprobante"`
	IDProgramacion    int64    `json:"id_programacion"`
	IDPasajero        int64    `json:"id_pasajero"`
	IDAsiento         int64    `json:"id_asiento"`
	Precio            float64  `json:"precio"`
	Descuento         *float64 `json:"descuento"`
	Serie             string   `json:"serie"`
	Correlativo       uint     `json:"correlativo"`
	NumeroComprobante string   `json:"numero_comprobante"`
	Nota              string   `json:"nota"`
	Observaciones     string   `json:"observaciones"`
	Subtotal          float64  `json:"subtotal"`
	IGV               float64  `json:"igv"`
	Total             float64  `json:"total"`
	QRCode            string   `json:"qr_code"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}
