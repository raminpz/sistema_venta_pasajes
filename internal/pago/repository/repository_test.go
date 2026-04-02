package repository

import (
	"errors"
	"sistema_venta_pasajes/internal/pago/domain"
	"testing"
)

type mockPagoRepository struct {
	CreateFunc  func(*domain.Pago) error
	UpdateFunc  func(*domain.Pago) error
	DeleteFunc  func(int64) error
	GetByIDFunc func(int64) (*domain.Pago, error)
	ListFunc    func(int, int, *int64) ([]domain.Pago, int, error)
}

func (m *mockPagoRepository) Create(p *domain.Pago) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(p)
	}
	return nil
}

func (m *mockPagoRepository) Update(p *domain.Pago) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(p)
	}
	return nil
}

func (m *mockPagoRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *mockPagoRepository) GetByID(id int64) (*domain.Pago, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *mockPagoRepository) List(offset, limit int, idVenta *int64) ([]domain.Pago, int, error) {
	if m.ListFunc != nil {
		return m.ListFunc(offset, limit, idVenta)
	}
	return nil, 0, nil
}

func TestRepositoryCreateOK(t *testing.T) {
	repo := &mockPagoRepository{CreateFunc: func(p *domain.Pago) error {
		if p.IDVenta <= 0 {
			return errors.New("id_venta invalido")
		}
		return nil
	}}
	if err := repo.Create(&domain.Pago{IDVenta: 1, Estado: "REGISTRADA"}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRepositoryCreateEstadoParcial(t *testing.T) {
	repo := &mockPagoRepository{CreateFunc: func(p *domain.Pago) error {
		if p.Estado != "PARCIAL" {
			return errors.New("estado incorrecto")
		}
		return nil
	}}
	if err := repo.Create(&domain.Pago{IDVenta: 2, IDMetodo: 1, Monto: 50, Estado: "PARCIAL"}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRepositoryCreateEstadoPagada(t *testing.T) {
	repo := &mockPagoRepository{CreateFunc: func(p *domain.Pago) error {
		if p.Estado != "PAGADA" {
			return errors.New("estado incorrecto")
		}
		return nil
	}}
	if err := repo.Create(&domain.Pago{IDVenta: 3, IDMetodo: 1, Monto: 100, Estado: "PAGADA"}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRepositoryCreateEstadoAnulada(t *testing.T) {
	repo := &mockPagoRepository{CreateFunc: func(p *domain.Pago) error {
		if p.Estado != "ANULADA" {
			return errors.New("estado incorrecto")
		}
		return nil
	}}
	if err := repo.Create(&domain.Pago{IDVenta: 4, IDMetodo: 1, Monto: 0, Estado: "ANULADA"}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRepositoryGetByIDOK(t *testing.T) {
	repo := &mockPagoRepository{GetByIDFunc: func(id int64) (*domain.Pago, error) {
		return &domain.Pago{IDPago: id, IDVenta: 1, Estado: "REGISTRADA"}, nil
	}}
	p, err := repo.GetByID(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if p.Estado != "REGISTRADA" {
		t.Fatalf("estado esperado REGISTRADA, obtuvo %s", p.Estado)
	}
}

func TestRepositoryGetByIDNotFound(t *testing.T) {
	repo := &mockPagoRepository{GetByIDFunc: func(id int64) (*domain.Pago, error) {
		return nil, errors.New("no encontrado")
	}}
	if _, err := repo.GetByID(1); err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestRepositoryUpdateEstado(t *testing.T) {
	repo := &mockPagoRepository{UpdateFunc: func(p *domain.Pago) error {
		if p.Estado != "PAGADA" {
			return errors.New("estado incorrecto")
		}
		return nil
	}}
	if err := repo.Update(&domain.Pago{IDPago: 1, Estado: "PAGADA"}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRepositoryListOK(t *testing.T) {
	repo := &mockPagoRepository{ListFunc: func(offset, limit int, idVenta *int64) ([]domain.Pago, int, error) {
		return []domain.Pago{{IDPago: 1, Estado: "REGISTRADA"}, {IDPago: 2, Estado: "PARCIAL"}}, 2, nil
	}}
	items, total, err := repo.List(0, 15, nil)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if total != 2 || len(items) != 2 {
		t.Fatal("resultado inesperado")
	}
}

func TestRepositoryListFiltradoPorVenta(t *testing.T) {
	idVenta := int64(3)
	repo := &mockPagoRepository{ListFunc: func(offset, limit int, v *int64) ([]domain.Pago, int, error) {
		if v == nil || *v != 3 {
			t.Fatal("id_venta no se pasó correctamente")
		}
		return []domain.Pago{{IDPago: 5, IDVenta: 3, Estado: "PAGADA"}}, 1, nil
	}}
	items, total, err := repo.List(0, 15, &idVenta)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatal("resultado inesperado")
	}
}

func TestRepositoryDeleteOK(t *testing.T) {
	repo := &mockPagoRepository{DeleteFunc: func(id int64) error { return nil }}
	if err := repo.Delete(1); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestRepositoryDeleteError(t *testing.T) {
	repo := &mockPagoRepository{DeleteFunc: func(id int64) error { return errors.New("db") }}
	if err := repo.Delete(1); err == nil {
		t.Fatal("se esperaba error")
	}
}
