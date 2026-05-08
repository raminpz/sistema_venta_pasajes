package service

import (
	"errors"
	"sistema_venta_pasajes/internal/pasajero/domain"
	"sistema_venta_pasajes/internal/pasajero/input"
	"sistema_venta_pasajes/internal/pasajero/util"
	"testing"
	"time"

	"gorm.io/gorm"
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

func TestPasajeroService_Update_List_Search_Delete(t *testing.T) {
	now := time.Now()
	repo := &fakeRepo{
		GetByIDFn: func(id int64) (*domain.Pasajero, error) {
			if id == 404 {
				return nil, gorm.ErrRecordNotFound
			}
			return &domain.Pasajero{IDPasajero: int(id), TipoDocumento: "DNI", NroDocumento: "12345678", Nombres: "Juan", Apellidos: "Perez", Telefono: "987654321", CreatedAt: now, UpdatedAt: now}, nil
		},
		UpdateFn: func(*domain.Pasajero) error { return nil },
		DeleteFn: func(id int64) error {
			if id == 404 {
				return gorm.ErrRecordNotFound
			}
			return nil
		},
		ListFn: func(page, size int) ([]domain.Pasajero, int, error) {
			return []domain.Pasajero{{IDPasajero: 1, TipoDocumento: "DNI", NroDocumento: "12345678", Nombres: "Juan", Apellidos: "Perez", Telefono: "987654321", CreatedAt: now, UpdatedAt: now}}, 1, nil
		},
		SearchFn: func(query string) ([]domain.Pasajero, int, error) {
			if query == "err" {
				return nil, 0, errors.New("db")
			}
			return []domain.Pasajero{{IDPasajero: 2, TipoDocumento: "DNI", NroDocumento: "87654321", Nombres: "Ana", Apellidos: "Lopez", Telefono: "987654321", CreatedAt: now, UpdatedAt: now}}, 1, nil
		},
	}

	svc := NewPasajeroService(repo)

	t.Run("update ok", func(t *testing.T) {
		out, err := svc.Update(1, input.UpdatePasajeroInput{Nombres: "luis", Apellidos: "diaz", Telefono: "987654321"})
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if out.Nombres != "Luis" {
			t.Fatalf("se esperaba capitalizacion, obtenido %s", out.Nombres)
		}
	})

	t.Run("update fecha invalida", func(t *testing.T) {
		f := "31-12-2027"
		_, err := svc.Update(1, input.UpdatePasajeroInput{FechaNacimiento: &f, Telefono: "987654321"})
		if err == nil || err.Error() != util.ERR_DATE_FORMAT {
			t.Fatalf("se esperaba error de fecha, obtenido %v", err)
		}
	})

	t.Run("list ok", func(t *testing.T) {
		out, meta, err := svc.List(1, 2)
		if err != nil || len(out) != 1 || meta.Total != 1 {
			t.Fatalf("resultado inesperado: len=%d total=%d err=%v", len(out), meta.Total, err)
		}
	})

	t.Run("search ok", func(t *testing.T) {
		out, err := svc.Search("ana")
		if err != nil || len(out) != 1 {
			t.Fatalf("resultado inesperado: len=%d err=%v", len(out), err)
		}
	})

	t.Run("search error", func(t *testing.T) {
		_, err := svc.Search("err")
		if err == nil {
			t.Fatal("se esperaba error")
		}
	})

	t.Run("delete not found", func(t *testing.T) {
		err := svc.Delete(404)
		if err == nil {
			t.Fatal("se esperaba not found")
		}
	})

	t.Run("delete ok", func(t *testing.T) {
		if err := svc.Delete(1); err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
	})
}
