package input

type CreateInput struct {
	RUC             string `json:"ruc"`
	RazonSocial     string `json:"razon_social"`
	NombreComercial string `json:"nombre_comercial"`
	Direccion       string `json:"direccion"`
	Telefono        string `json:"telefono"`
	Email           string `json:"email"`
	Web             string `json:"web"`
}

type UpdateInput struct {
	RUC             string `json:"ruc"`
	RazonSocial     string `json:"razon_social"`
	NombreComercial string `json:"nombre_comercial"`
	Direccion       string `json:"direccion"`
	Telefono        string `json:"telefono"`
	Email           string `json:"email"`
	Web             string `json:"web"`
}

type ProveedorOutput struct {
	ID              int64  `json:"id_proveedor"`
	RUC             string `json:"ruc"`
	RazonSocial     string `json:"razon_social"`
	NombreComercial string `json:"nombre_comercial"`
	Direccion       string `json:"direccion"`
	Telefono        string `json:"telefono"`
	Email           string `json:"email"`
	Web             string `json:"web"`
}
