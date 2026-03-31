package repository

import (
	"context"
	"sistema_venta_pasajes/internal/usuario/domain"
	"testing"
)

func TestUsuarioRepositoryMock(t *testing.T) {
	repo := NewUsuarioRepositoryMock()
	usuario := &domain.Usuario{IDUsuario: 1, Nombre: "Juan"}
	err := repo.Create(context.Background(), usuario)
	if err != nil {
		t.Fatalf("error al crear usuario: %v", err)
	}
	got, err := repo.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("error al obtener usuario: %v", err)
	}
	if got.Nombre != "Juan" {
		t.Errorf("esperado Juan, obtenido %s", got.Nombre)
	}
}
