package service

import (
	"context"
	"errors"
	"sistema_venta_pasajes/internal/conductor/domain"
	"sistema_venta_pasajes/internal/conductor/input"
	"testing"
)

type mockRepo struct{}

func (m *mockRepo) List() ([]domain.Conductor, error) {
	return []domain.Conductor{{IDConductor: 1, Nombres: "Juan", Apellidos: "Perez", NumeroLicencia: "ABC123456", Telefono: "987654321"}}, nil
}
func (m *mockRepo) GetByID(id int64) (*domain.Conductor, error) {
	if id == 1 {
		return &domain.Conductor{IDConductor: 1, Nombres: "Juan", Apellidos: "Perez", NumeroLicencia: "ABC123456", Telefono: "987654321"}, nil
	}
	return nil, errors.New("not found")
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
