package service

import (
	"errors"
	"testing"

	"sistema_venta_pasajes/internal/asiento/domain"
	"sistema_venta_pasajes/internal/asiento/input"
)

type mockRepo struct {
	CreateFn         func(*domain.Asiento) error
	GetByIDFn        func(int64) (*domain.Asiento, error)
	ListByVehiculoFn func(int64) ([]*domain.Asiento, error)
	UpdateFn         func(*domain.Asiento) error
	DeleteFn         func(int64) error
	CambiarEstadoFn  func(int64, string) error
}

func (m *mockRepo) Create(a *domain.Asiento) error                     { return m.CreateFn(a) }
func (m *mockRepo) GetByID(id int64) (*domain.Asiento, error)          { return m.GetByIDFn(id) }
func (m *mockRepo) ListByVehiculo(id int64) ([]*domain.Asiento, error) { return m.ListByVehiculoFn(id) }
func (m *mockRepo) Update(a *domain.Asiento) error                     { return m.UpdateFn(a) }
func (m *mockRepo) Delete(id int64) error                              { return m.DeleteFn(id) }
func (m *mockRepo) CambiarEstado(id int64, estado string) error        { return m.CambiarEstadoFn(id, estado) }

func TestAsientoService_Create(t *testing.T) {
	repo := &mockRepo{
		CreateFn: func(a *domain.Asiento) error {
			a.IDAsiento = 1
			return nil
		},
	}
	svc := NewAsientoService(repo)
	in := input.CreateAsientoInput{IDVehiculo: 2, NumeroAsiento: "A1"}
	asiento, err := svc.Create(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if asiento.IDAsiento != 1 || asiento.IDVehiculo != 2 || asiento.NumeroAsiento != "A1" {
		t.Errorf("unexpected asiento: %+v", asiento)
	}
}

func TestAsientoService_GetByID(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(id int64) (*domain.Asiento, error) {
			if id == 1 {
				return &domain.Asiento{IDAsiento: 1, IDVehiculo: 2, NumeroAsiento: "A1"}, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := NewAsientoService(repo)
	asiento, err := svc.GetByID(1)
	if err != nil || asiento.IDAsiento != 1 {
		t.Errorf("expected asiento, got %v, err %v", asiento, err)
	}
	_, err = svc.GetByID(2)
	if err == nil {
		t.Error("expected error for not found")
	}
}

func TestAsientoService_ListByVehiculo(t *testing.T) {
	repo := &mockRepo{
		ListByVehiculoFn: func(id int64) ([]*domain.Asiento, error) {
			if id == 2 {
				return []*domain.Asiento{{IDAsiento: 1, IDVehiculo: 2, NumeroAsiento: "A1"}}, nil
			}
			return nil, nil
		},
	}
	svc := NewAsientoService(repo)
	asientos, err := svc.ListByVehiculo(2)
	if err != nil || len(asientos) != 1 {
		t.Errorf("expected 1 asiento, got %v, err %v", asientos, err)
	}
}

func TestAsientoService_Update(t *testing.T) {
	updated := false
	repo := &mockRepo{
		GetByIDFn: func(id int64) (*domain.Asiento, error) {
			return &domain.Asiento{IDAsiento: int(id), NumeroAsiento: "A1"}, nil
		},
		UpdateFn: func(a *domain.Asiento) error {
			if a.NumeroAsiento == "B2" {
				updated = true
			}
			return nil
		},
	}
	svc := NewAsientoService(repo)
	err := svc.Update(1, input.UpdateAsientoInput{NumeroAsiento: "B2"})
	if err != nil || !updated {
		t.Errorf("expected update, got err %v, updated %v", err, updated)
	}
}

func TestAsientoService_Delete(t *testing.T) {
	called := false
	repo := &mockRepo{
		DeleteFn: func(id int64) error {
			called = true
			if id != 1 {
				return errors.New("not found")
			}
			return nil
		},
	}
	svc := NewAsientoService(repo)
	err := svc.Delete(1)
	if err != nil || !called {
		t.Errorf("expected delete, got err %v, called %v", err, called)
	}
}
