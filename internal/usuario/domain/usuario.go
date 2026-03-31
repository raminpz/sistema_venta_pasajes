package domain

import "time"

// Usuario representa la entidad de usuario del sistema.
type Usuario struct {
	IDUsuario int       `json:"id_usuario" gorm:"primaryKey;autoIncrement"`
	IDRol     int       `json:"id_rol"`
	Nombre    string    `json:"nombre"`
	Apellidos string    `json:"apellidos" gorm:"column:APELLIDOS"`
	DNI       string    `json:"dni"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Telefono  string    `json:"telefono"`
	Estado    string    `json:"estado"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Usuario) TableName() string {
	return "USUARIO"
}
