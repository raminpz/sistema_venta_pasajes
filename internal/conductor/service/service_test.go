package service

import (
	"context"
	"errors"
	"sistema_venta_pasajes/internal/conductor/domain"
	"sistema_venta_pasajes/internal/conductor/input"
	"testing"

	"gorm.io/gorm"
)

type mockRepo struct{}

func (m *mockRepo) List() ([]domain.Conductor, error) {
	return []domain.Conductor{{IDConductor: 1, Nombres: "Juan", Apellidos: "Perez", NumeroLicencia: "ABC123456", Telefono: "987654321"}}, nil
}
func (m *mockRepo) GetByID(id int64) (*domain.Conductor, error) {
	if id == 1 {
		return &domain.Conductor{IDConductor: 1, Nombres: "Juan", Apellidos: "Perez", NumeroLicencia: "ABC123456", Telefono: "987654321"}, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockRepo) Create(c *domain.Conductor) error {
	c.IDConductor = 2
	return nil
}
func (m *mockRepo) Update(c *domain.Conductor) error {
	if c.IDConductor == 0 {
		return errors.New("not found")
	}
	return nil
}
func (m *mockRepo) Delete(id int64) error {
	if id == 1 {
		return nil
	}
	return errors.New("not found")
}

func TestService_Create(t *testing.T) {
	s := New(&mockRepo{})
	in := input.CreateConductorInput{
		Nombres:           "Juan",
		Apellidos:         "Perez",
		DNI:               "12345678",
		NumeroLicencia:    "ABC123456",
		Telefono:          "987654321",
		FechaVencLicencia: "2027-03-29",
	}
	c, err := s.Create(context.Background(), in)
	if err != nil || c == nil || c.NumeroLicencia != "ABC123456" {
		t.Errorf("Create failed: %v", err)
	}
}

func TestService_Create_invalid(t *testing.T) {
	s := New(&mockRepo{})
	in := input.CreateConductorInput{
		Nombres:           "Juan",
		Apellidos:         "Perez",
		DNI:               "12345678",
		NumeroLicencia:    "1234",
		Telefono:          "987654321",
		FechaVencLicencia: "2027-03-29",
	}
	_, err := s.Create(context.Background(), in)
	if err == nil {
		t.Errorf("Expected error for invalid control_acceso")
	}
}

func TestService_GetByID(t *testing.T) {
	s := New(&mockRepo{})
	if _, err := s.GetByID(context.Background(), 0); err == nil {
		t.Fatal("se esperaba error por id invalido")
	}
	if _, err := s.GetByID(context.Background(), 999); err == nil {
		t.Fatal("se esperaba not found")
	}
	if _, err := s.GetByID(context.Background(), 1); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestService_UpdateAndDelete(t *testing.T) {
	s := New(&mockRepo{})

	if _, err := s.Update(context.Background(), 0, input.UpdateConductorInput{}); err == nil {
		t.Fatal("se esperaba error por id invalido")
	}

	if _, err := s.Update(context.Background(), 999, input.UpdateConductorInput{}); err == nil {
		t.Fatal("se esperaba not found en update")
	}

	nombres := "mario"
	up := input.UpdateConductorInput{Nombres: &nombres}
	if _, err := s.Update(context.Background(), 1, up); err != nil {
		t.Fatalf("no se esperaba error en update: %v", err)
	}

	if err := s.Delete(context.Background(), 0); err == nil {
		t.Fatal("se esperaba error por id invalido")
	}
	if err := s.Delete(context.Background(), 999); err == nil {
		t.Fatal("se esperaba not found en delete")
	}
	if err := s.Delete(context.Background(), 1); err != nil {
		t.Fatalf("no se esperaba error en delete: %v", err)
	}
}

func TestService_Create_invalidDate(t *testing.T) {
	s := New(&mockRepo{})
	in := input.CreateConductorInput{
		Nombres:           "Juan",
		Apellidos:         "Perez",
		DNI:               "12345678",
		NumeroLicencia:    "ABC123456",
		Telefono:          "987654321",
		FechaVencLicencia: "29-03-2027",
	}
	if _, err := s.Create(context.Background(), in); err == nil {
		t.Fatal("se esperaba error por fecha invalida")
	}
}
