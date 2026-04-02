package repository

import (
	"errors"
	"sistema_venta_pasajes/internal/encomienda/domain"
	"testing"
)

type mockEncomiendaRepository struct {
	CreateFunc  func(*domain.Encomienda) error
	UpdateFunc  func(*domain.Encomienda) error
	DeleteFunc  func(int64) error
	GetByIDFunc func(int64) (*domain.Encomienda, error)
	ListFunc    func(int, int) ([]domain.Encomienda, int, error)
}

func (m *mockEncomiendaRepository) Create(e *domain.Encomienda) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(e)
	}
	return nil
}

func (m *mockEncomiendaRepository) Update(e *domain.Encomienda) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(e)
	}
	return nil
}

func (m *mockEncomiendaRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *mockEncomiendaRepository) GetByID(id int64) (*domain.Encomienda, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *mockEncomiendaRepository) List(offset, limit int) ([]domain.Encomienda, int, error) {
	if m.ListFunc != nil {
		return m.ListFunc(offset, limit)
	}
	return nil, 0, nil
}

func TestRepositoryCreateOK(t *testing.T) {
	repo := &mockEncomiendaRepository{CreateFunc: func(e *domain.Encomienda) error {
		if e.IDVenta <= 0 {
			return errors.New("id_venta invalido")
		}
		return nil
	}}
	if err := repo.Create(&domain.Encomienda{IDVenta: 1}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRepositoryGetByIDNotFound(t *testing.T) {
	repo := &mockEncomiendaRepository{GetByIDFunc: func(id int64) (*domain.Encomienda, error) {
		return nil, errors.New("no encontrado")
	}}
	if _, err := repo.GetByID(1); err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestRepositoryListOK(t *testing.T) {
	repo := &mockEncomiendaRepository{ListFunc: func(offset, limit int) ([]domain.Encomienda, int, error) {
		return []domain.Encomienda{{IDEncomienda: 1}}, 1, nil
	}}
	items, total, err := repo.List(0, 15)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatal("resultado inesperado")
	}
}

func TestRepositoryDeleteError(t *testing.T) {
	repo := &mockEncomiendaRepository{DeleteFunc: func(id int64) error { return errors.New("db") }}
	if err := repo.Delete(1); err == nil {
		t.Fatal("se esperaba error")
	}
}
