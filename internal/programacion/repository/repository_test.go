package repository

import (
	"errors"
	"sistema_venta_pasajes/internal/programacion/domain"
	"testing"
)

type mockProgramacionRepository struct {
	CreateFunc  func(*domain.Programacion) error
	UpdateFunc  func(*domain.Programacion) error
	DeleteFunc  func(int64) error
	GetByIDFunc func(int64) (*domain.Programacion, error)
	ListFunc    func(int, int) ([]domain.Programacion, int, error)
}

func (m *mockProgramacionRepository) Create(p *domain.Programacion) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(p)
	}
	return nil
}

func (m *mockProgramacionRepository) Update(p *domain.Programacion) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(p)
	}
	return nil
}

func (m *mockProgramacionRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *mockProgramacionRepository) GetByID(id int64) (*domain.Programacion, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *mockProgramacionRepository) List(offset, limit int) ([]domain.Programacion, int, error) {
	if m.ListFunc != nil {
		return m.ListFunc(offset, limit)
	}
	return nil, 0, nil
}

func TestRepositoryCreateOK(t *testing.T) {
	repo := &mockProgramacionRepository{CreateFunc: func(p *domain.Programacion) error {
		if p.IDRuta <= 0 {
			return errors.New("id_ruta invalido")
		}
		return nil
	}}
	if err := repo.Create(&domain.Programacion{IDRuta: 1}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRepositoryGetByIDNotFound(t *testing.T) {
	repo := &mockProgramacionRepository{GetByIDFunc: func(id int64) (*domain.Programacion, error) {
		return nil, errors.New("no encontrado")
	}}
	if _, err := repo.GetByID(1); err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestRepositoryListOK(t *testing.T) {
	repo := &mockProgramacionRepository{ListFunc: func(offset, limit int) ([]domain.Programacion, int, error) {
		return []domain.Programacion{{IDProgramacion: 1}}, 1, nil
	}}
	items, total, err := repo.List(0, 15)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("resultado inesperado")
	}
}

func TestRepositoryDeleteError(t *testing.T) {
	repo := &mockProgramacionRepository{DeleteFunc: func(id int64) error { return errors.New("db") }}
	if err := repo.Delete(1); err == nil {
		t.Fatal("se esperaba error")
	}
}
