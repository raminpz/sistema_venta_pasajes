package domain

import (
	"time"
)

type Terminal struct {
	IDTerminal   int64      `json:"id_terminal" gorm:"column:ID_TERMINAL;primaryKey;autoIncrement"`
	NOMBRE       string     `json:"nombre" gorm:"column:NOMBRE;not null"`
	CIUDAD       string     `json:"ciudad" gorm:"column:CIUDAD;not null"`
	DEPARTAMENTO string     `json:"departamento" gorm:"column:DEPARTAMENTO;not null"`
	DIRECCION    string     `json:"direccion" gorm:"column:DIRECCION;not null"`
	ESTADO       string     `json:"estado" gorm:"column:ESTADO"`
	CreatedAt    *time.Time `json:"-" gorm:"column:CREATED_AT"`
	UpdatedAt    *time.Time `json:"-" gorm:"column:UPDATED_AT"`
}

func (Terminal) TableName() string {
	return "TERMINAL"
}
