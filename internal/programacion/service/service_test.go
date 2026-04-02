package service

import (
	"errors"
	"sistema_venta_pasajes/internal/programacion/domain"
	"sistema_venta_pasajes/internal/programacion/input"
	"testing"
	"time"
)

type fakeProgramacionRepo struct {
	createFn  func(*domain.Programacion) error
	updateFn  func(*domain.Programacion) error
	deleteFn  func(int64) error
	getByIDFn func(int64) (*domain.Programacion, error)
	listFn    func(int, int) ([]domain.Programacion, int, error)
}

func (f *fakeProgramacionRepo) Create(p *domain.Programacion) error {
	if f.createFn != nil {
		return f.createFn(p)
	}
	return nil
}

func (f *fakeProgramacionRepo) Update(p *domain.Programacion) error {
	if f.updateFn != nil {
		return f.updateFn(p)
	}
	return nil
}

func (f *fakeProgramacionRepo) Delete(id int64) error {
	if f.deleteFn != nil {
		return f.deleteFn(id)
	}
	return nil
}

func (f *fakeProgramacionRepo) GetByID(id int64) (*domain.Programacion, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(id)
	}
	return nil, nil
}

func (f *fakeProgramacionRepo) List(offset, limit int) ([]domain.Programacion, int, error) {
	if f.listFn != nil {
		return f.listFn(offset, limit)
	}
	return nil, 0, nil
}

func TestServiceCreateOK(t *testing.T) {
	repo := &fakeProgramacionRepo{createFn: func(p *domain.Programacion) error {
		if p.Estado != "PROGRAMADO" {
			t.Fatalf("estado esperado PROGRAMADO, obtuvo %s", p.Estado)
		}
		return nil
	}}
	s := NewProgramacionService(repo)

	out, err := s.Create(input.CreateProgramacionInput{
		IDRuta:      1,
		IDVehiculo:  2,
		IDConductor: 3,
		FechaSalida: "2026-04-05 02:00:00",
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out == nil || out.IDRuta != 1 {
		t.Fatal("salida invalida")
	}
}

func TestServiceCreateInvalidDates(t *testing.T) {
	repo := &fakeProgramacionRepo{}
	s := NewProgramacionService(repo)
	llegada := "2026-04-05 01:00:00"
	_, err := s.Create(input.CreateProgramacionInput{
		IDRuta:       1,
		IDVehiculo:   2,
		IDConductor:  3,
		FechaSalida:  "2026-04-05 02:00:00",
		FechaLlegada: &llegada,
	})
	if err == nil {
		t.Fatal("se esperaba error de fechas")
	}
}

func TestServiceUpdateNotFound(t *testing.T) {
	repo := &fakeProgramacionRepo{getByIDFn: func(int64) (*domain.Programacion, error) {
		return nil, errors.New("not found")
	}}
	s := NewProgramacionService(repo)
	_, err := s.Update(1, input.UpdateProgramacionInput{})
	if err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestServiceListOK(t *testing.T) {
	now := time.Now()
	repo := &fakeProgramacionRepo{listFn: func(offset, limit int) ([]domain.Programacion, int, error) {
		if offset != 5 || limit != 5 {
			t.Fatalf("offset/limit esperados 5/5, obtuvo %d/%d", offset, limit)
		}
		return []domain.Programacion{{
			IDProgramacion: 1,
			IDRuta:         1,
			IDVehiculo:     1,
			IDConductor:    1,
			FechaSalida:    now,
			Estado:         "PROGRAMADO",
		}}, 1, nil
	}}
	s := NewProgramacionService(repo)
	outs, total, err := s.List(2, 5)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if total != 1 || len(outs) != 1 {
		t.Fatalf("resultado inesperado total=%d len=%d", total, len(outs))
	}
}

func TestServiceDeleteError(t *testing.T) {
	repo := &fakeProgramacionRepo{deleteFn: func(int64) error { return errors.New("db") }}
	s := NewProgramacionService(repo)
	if err := s.Delete(1); err == nil {
		t.Fatal("se esperaba error")
	}
}
