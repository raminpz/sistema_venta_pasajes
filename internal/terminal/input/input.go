package input

type CreateTerminalInput struct {
	Nombre       string `json:"nombre" binding:"required"`
	Ciudad       string `json:"ciudad" binding:"required"`
	Departamento string `json:"departamento" binding:"required"`
	Direccion    string `json:"direccion" binding:"required"`
	Estado       string `json:"estado"`
}

type UpdateTerminalInput struct {
	Nombre       string `json:"nombre"`
	Ciudad       string `json:"ciudad"`
	Departamento string `json:"departamento"`
	Direccion    string `json:"direccion"`
	Estado       string `json:"estado"`
}

type TerminalOutput struct {
	IDTerminal   int64  `json:"id_terminal"`
	Nombre       string `json:"nombre"`
	Ciudad       string `json:"ciudad"`
	Departamento string `json:"departamento"`
	Direccion    string `json:"direccion"`
	Estado       string `json:"estado"`
}
