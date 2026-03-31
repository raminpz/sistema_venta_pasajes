package service

import (
	"sistema_venta_pasajes/internal/usuario/input"
	"sistema_venta_pasajes/internal/usuario/repository"
	"testing"
)

func TestCrearUsuario(t *testing.T) {
	repo := repository.NewUsuarioRepositoryMock()
	serv := NewUsuarioService(repo)
	in := &input.UsuarioCreateInput{
		IDRol:     1,
		Nombre:    "Juan",
		Apellidos: "Perez",
		DNI:       "12345678",
		Email:     "juan@mail.com",
		Password:  "1234",
		Telefono:  "999999999",
	}
	usuario, err := serv.Create(*in)
	if err != nil {
		t.Fatalf("error al crear usuario: %v", err)
	}
	if usuario.Nombre != "Juan" {
		t.Errorf("esperado nombre Juan, obtenido %s", usuario.Nombre)
	}
	if usuario.Apellidos != "Perez" {
		t.Errorf("esperado apellidos Perez, obtenido %s", usuario.Apellidos)
	}
}
