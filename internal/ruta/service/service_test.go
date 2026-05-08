package service

import (
	"context"
	"errors"
	"reflect"
	"sistema_venta_pasajes/internal/ruta/domain"
	"sistema_venta_pasajes/internal/ruta/input"
	"sistema_venta_pasajes/pkg"
	"testing"

	"gorm.io/gorm"
)

type mockRepo struct {
	CreateFn  func(ruta *domain.Ruta) error
	GetByIDFn func(id int) (*domain.Ruta, error)
	UpdateFn  func(ruta *domain.Ruta) error
	DeleteFn  func(id int) error
	ListFn    func() ([]domain.Ruta, error)
}

func (m *mockRepo) Create(ruta *domain.Ruta) error       { return m.CreateFn(ruta) }
func (m *mockRepo) GetByID(id int) (*domain.Ruta, error) { return m.GetByIDFn(id) }
func (m *mockRepo) Update(ruta *domain.Ruta) error       { return m.UpdateFn(ruta) }
func (m *mockRepo) Delete(id int) error                  { return m.DeleteFn(id) }
func (m *mockRepo) List() ([]domain.Ruta, error)         { return m.ListFn() }

func TestService_Create(t *testing.T) {
	repo := &mockRepo{
		CreateFn: func(ruta *domain.Ruta) error {
			ruta.IDRuta = 1
			return nil
		},
	}
	svc := New(repo)
	in := input.CreateRutaInput{IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5}
	ruta, err := svc.Create(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ruta.IDRuta != 1 || ruta.IDOrigenTerminal != 1 || ruta.IDDestinoTerminal != 2 || ruta.DuracionHoras != 5.5 {
		t.Errorf("unexpected ruta: %+v", ruta)
	}
}

func TestService_GetByID(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(id int) (*domain.Ruta, error) {
			if id == 1 {
				return &domain.Ruta{IDRuta: 1, IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5}, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := New(repo)
	ruta, err := svc.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ruta.IDRuta != 1 {
		t.Errorf("unexpected ruta: %+v", ruta)
	}
}

func TestService_Update(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(id int) (*domain.Ruta, error) {
			return &domain.Ruta{IDRuta: id, IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 5.5}, nil
		},
		UpdateFn: func(ruta *domain.Ruta) error {
			if ruta.DuracionHoras != 7.0 {
				return errors.New("update failed")
			}
			return nil
		},
	}
	svc := New(repo)
	in := input.UpdateRutaInput{DuracionHoras: ptrFloat(7.0)}
	ruta, err := svc.Update(context.Background(), 1, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ruta.DuracionHoras != 7.0 {
		t.Errorf("unexpected ruta: %+v", ruta)
	}
}

func TestService_Delete(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(id int) (*domain.Ruta, error) {
			return &domain.Ruta{IDRuta: id}, nil
		},
		DeleteFn: func(id int) error {
			if id != 1 {
				return errors.New("delete failed")
			}
			return nil
		},
	}
	svc := New(repo)
	if err := svc.Delete(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_Delete_DeleteError(t *testing.T) {
	repo := &mockRepo{
		GetByIDFn: func(id int) (*domain.Ruta, error) {
			return &domain.Ruta{IDRuta: id}, nil
		},
		DeleteFn: func(id int) error {
			return errors.New("fk constraint")
		},
	}
	svc := New(repo)
	err := svc.Delete(context.Background(), 1)
	if err == nil {
		t.Fatal("expected delete error")
	}
}

func TestService_List(t *testing.T) {
	repo := &mockRepo{
		ListFn: func() ([]domain.Ruta, error) {
			return []domain.Ruta{{IDRuta: 1}, {IDRuta: 2}}, nil
		},
	}
	svc := New(repo)
	rutas, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(rutas, []domain.Ruta{{IDRuta: 1}, {IDRuta: 2}}) {
		t.Errorf("unexpected rutas: %+v", rutas)
	}
}

func ptrFloat(f float64) *float64 { return &f }

func TestService_Create_DuplicateAndInternalError(t *testing.T) {
	svcDuplicate := New(&mockRepo{
		CreateFn: func(ruta *domain.Ruta) error {
			return pkg.Conflict("duplicate_resource", "duplicado")
		},
	})
	_, err := svcDuplicate.Create(context.Background(), input.CreateRutaInput{IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 1})
	if err == nil {
		t.Fatal("se esperaba error por duplicado")
	}

	svcInternal := New(&mockRepo{
		CreateFn: func(ruta *domain.Ruta) error {
			return errors.New("db down")
		},
	})
	_, err = svcInternal.Create(context.Background(), input.CreateRutaInput{IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 1})
	if err == nil {
		t.Fatal("se esperaba error interno")
	}
}

func TestService_UpdateAndListErrorBranches(t *testing.T) {
	svc := New(&mockRepo{
		GetByIDFn: func(id int) (*domain.Ruta, error) {
			return &domain.Ruta{IDRuta: id, IDOrigenTerminal: 1, IDDestinoTerminal: 2, DuracionHoras: 2.5}, nil
		},
		UpdateFn: func(ruta *domain.Ruta) error {
			return errors.New("update failed")
		},
		ListFn: func() ([]domain.Ruta, error) {
			return nil, errors.New("db")
		},
	})

	_, err := svc.Update(context.Background(), 1, input.UpdateRutaInput{DuracionHoras: ptrFloat(3.2)})
	if err == nil {
		t.Fatal("se esperaba error de update")
	}

	_, err = svc.List(context.Background())
	if err == nil {
		t.Fatal("se esperaba error en list")
	}
}

func TestService_Delete_NotFoundByGorm(t *testing.T) {
	svc := New(&mockRepo{
		GetByIDFn: func(id int) (*domain.Ruta, error) {
			return nil, gorm.ErrRecordNotFound
		},
		DeleteFn: func(id int) error { return nil },
	})

	err := svc.Delete(context.Background(), 1)
	if err == nil {
		t.Fatal("se esperaba not found")
	}
}
