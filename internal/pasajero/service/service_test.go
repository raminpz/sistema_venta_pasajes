package service

import (
	"errors"
	"sistema_venta_pasajes/internal/pasajero/domain"
	"sistema_venta_pasajes/internal/pasajero/input"
	"testing"
)

type fakeRepo struct {
	CreateFn  func(*domain.Pasajero) error
	GetByIDFn func(int64) (*domain.Pasajero, error)
	UpdateFn  func(*domain.Pasajero) error
	DeleteFn  func(int64) error
	ListFn    func(page, size int) ([]domain.Pasajero, int, error)
	SearchFn  func(string) ([]domain.Pasajero, int, error)
}

func (f *fakeRepo) Create(p *domain.Pasajero) error            { return f.CreateFn(p) }
func (f *fakeRepo) GetByID(id int64) (*domain.Pasajero, error) { return f.GetByIDFn(id) }
func (f *fakeRepo) Update(p *domain.Pasajero) error            { return f.UpdateFn(p) }
func (f *fakeRepo) Delete(id int64) error                      { return f.DeleteFn(id) }
func (f *fakeRepo) List(page, size int) ([]domain.Pasajero, int, error) {
	if f.ListFn != nil {
		return f.ListFn(page, size)
	}
	return nil, 0, nil
}
func (f *fakeRepo) Search(query string) ([]domain.Pasajero, int, error) {
	if f.SearchFn != nil {
		return f.SearchFn(query)
	}
	return nil, 0, nil
}

func TestCreatePasajero_Valid(t *testing.T) {
	repo := &fakeRepo{
		CreateFn: func(p *domain.Pasajero) error { return nil },
	}
	svc := NewPasajeroService(repo)
	in := input.CreatePasajeroInput{
		TipoDocumento: "DNI",
		NroDocumento:  "12345678",
		Nombres:       "Juan",
		Apellidos:     "Perez",
		Telefono:      "987654321",
	}
	_, err := svc.Create(in)
	if err != nil {
		t.Fatalf("esperaba nil, obtuve: %v", err)
	}
}

func TestCreatePasajero_InvalidTelefono(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewPasajeroService(repo)
	in := input.CreatePasajeroInput{
		TipoDocumento: "DNI",
		NroDocumento:  "12345678",
		Nombres:       "Juan",
		Apellidos:     "Perez",
		Telefono:      "1234",
	}
	_, err := svc.Create(in)
	if err == nil {
		t.Fatal("esperaba error por teléfono inválido")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &fakeRepo{
		GetByIDFn: func(id int64) (*domain.Pasajero, error) { return nil, errors.New("not found") },
	}
	svc := NewPasajeroService(repo)
	_, err := svc.GetByID(1)
	if err == nil {
		t.Fatal("esperaba error por pasajero no encontrado")
	}
}
