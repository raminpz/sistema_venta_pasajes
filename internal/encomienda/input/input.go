package input

type CreateEncomiendaInput struct {
	IDVenta            int64    `json:"id_venta"`
	IDProgramacion     int64    `json:"id_programacion"`
	Descripcion        string   `json:"descripcion"`
	PesoKg             *float64 `json:"peso_kg"`
	Costo              float64  `json:"costo"`
	RemitenteNombre    string   `json:"remitente_nombre"`
	RemitenteDoc       string   `json:"remitente_doc"`
	DestinatarioNombre string   `json:"destinatario_nombre"`
	DestinatarioDoc    *string  `json:"destinatario_doc"`
	DestinatarioTel    string   `json:"destinatario_tel"`
	Estado             string   `json:"estado"`
}

type UpdateEncomiendaInput struct {
	IDVenta            int64    `json:"id_venta"`
	IDProgramacion     int64    `json:"id_programacion"`
	Descripcion        string   `json:"descripcion"`
	PesoKg             *float64 `json:"peso_kg"`
	Costo              float64  `json:"costo"`
	RemitenteNombre    string   `json:"remitente_nombre"`
	RemitenteDoc       string   `json:"remitente_doc"`
	DestinatarioNombre string   `json:"destinatario_nombre"`
	DestinatarioDoc    *string  `json:"destinatario_doc"`
	DestinatarioTel    string   `json:"destinatario_tel"`
	Estado             string   `json:"estado"`
}

type EncomiendaOutput struct {
	IDEncomienda       int64    `json:"id_encomienda"`
	IDVenta            int64    `json:"id_venta"`
	IDProgramacion     int64    `json:"id_programacion"`
	Descripcion        string   `json:"descripcion"`
	PesoKg             *float64 `json:"peso_kg"`
	Costo              float64  `json:"costo"`
	RemitenteNombre    string   `json:"remitente_nombre"`
	RemitenteDoc       string   `json:"remitente_doc"`
	DestinatarioNombre string   `json:"destinatario_nombre"`
	DestinatarioDoc    *string  `json:"destinatario_doc"`
	DestinatarioTel    string   `json:"destinatario_tel"`
	Estado             string   `json:"estado"`
	CreatedAt          *string  `json:"created_at,omitempty"`
	UpdatedAt          *string  `json:"updated_at,omitempty"`
}
