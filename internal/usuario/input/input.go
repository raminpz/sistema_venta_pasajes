package input

// UsuarioCreateInput representa los datos necesarios para crear un usuario.
type UsuarioCreateInput struct {
	IDRol     int    `json:"id_rol"`
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	DNI       string `json:"dni"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Telefono  string `json:"telefono"`
}

// UsuarioUpdateInput representa los datos para actualizar un usuario.
type UsuarioUpdateInput struct {
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Telefono  string `json:"telefono"`
	Estado    string `json:"estado"`
}

// UsuarioOutput representa la respuesta de usuario para la API.
type UsuarioOutput struct {
	IDUsuario int    `json:"id_usuario"`
	IDRol     int    `json:"id_rol"`
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	DNI       string `json:"dni"`
	Email     string `json:"email"`
	Telefono  string `json:"telefono"`
	Estado    string `json:"estado"`
}
