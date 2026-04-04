package repository

import (
	"errors"
	"sistema_venta_pasajes/internal/venta/domain"
	"testing"
)

// ---------------------------------------------------------------------------
// Mock del repositorio
// ---------------------------------------------------------------------------

type mockVentaRepository struct {
	CreateFunc          func(venta *domain.Venta) error
	UpdateFunc          func(venta *domain.Venta) error
	DeleteFunc          func(id int64) error
	GetByIDFunc         func(id int64) (*domain.Venta, error)
	ListFunc            func(offset, limit int) ([]domain.Venta, int, error)
	NextCorrelativoFunc func(serie string) (uint, error)
}

func (m *mockVentaRepository) Create(venta *domain.Venta) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(venta)
	}
	return nil
}

func (m *mockVentaRepository) Update(venta *domain.Venta) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(venta)
	}
	return nil
}

func (m *mockVentaRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *mockVentaRepository) GetByID(id int64) (*domain.Venta, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *mockVentaRepository) List(offset, limit int) ([]domain.Venta, int, error) {
	if m.ListFunc != nil {
		return m.ListFunc(offset, limit)
	}
	return nil, 0, nil
}

func (m *mockVentaRepository) NextCorrelativo(serie string) (uint, error) {
	if m.NextCorrelativoFunc != nil {
		return m.NextCorrelativoFunc(serie)
	}
	return 1, nil
}

// ---------------------------------------------------------------------------
// Tests Create
// ---------------------------------------------------------------------------

func TestVentaRepository_Create_OK(t *testing.T) {
	venta := &domain.Venta{IDUsuario: 1, IDProgramacion: 10, IDPasajero: 20, IDAsiento: 5, Precio: 80, Serie: "B001", Correlativo: 1, Subtotal: 80}
	repo := &mockVentaRepository{
		CreateFunc: func(v *domain.Venta) error {
			if v.Serie != "B001" {
				return errors.New("serie incorrecta")
			}
			return nil
		},
	}
	if err := repo.Create(venta); err != nil {
		t.Errorf("no se esperaba error al crear: %v", err)
	}
}

func TestVentaRepository_Create_Error(t *testing.T) {
	repo := &mockVentaRepository{
		CreateFunc: func(v *domain.Venta) error {
			return errors.New("error de BD")
		},
	}
	err := repo.Create(&domain.Venta{})
	if err == nil {
		t.Error("se esperaba error al crear")
	}
}

// ---------------------------------------------------------------------------
// Tests Update
// ---------------------------------------------------------------------------

func TestVentaRepository_Update_OK(t *testing.T) {
	venta := &domain.Venta{IDVenta: 1, Nota: "actualizada"}
	repo := &mockVentaRepository{
		UpdateFunc: func(v *domain.Venta) error {
			if v.IDVenta != 1 {
				return errors.New("id incorrecto")
			}
			return nil
		},
	}
	if err := repo.Update(venta); err != nil {
		t.Errorf("no se esperaba error al actualizar: %v", err)
	}
}

func TestVentaRepository_Update_Error(t *testing.T) {
	repo := &mockVentaRepository{
		UpdateFunc: func(v *domain.Venta) error {
			return errors.New("error de BD")
		},
	}
	if err := repo.Update(&domain.Venta{}); err == nil {
		t.Error("se esperaba error al actualizar")
	}
}

// ---------------------------------------------------------------------------
// Tests Delete
// ---------------------------------------------------------------------------

func TestVentaRepository_Delete_OK(t *testing.T) {
	repo := &mockVentaRepository{
		DeleteFunc: func(id int64) error {
			if id != 5 {
				return errors.New("id incorrecto")
			}
			return nil
		},
	}
	if err := repo.Delete(5); err != nil {
		t.Errorf("no se esperaba error al eliminar: %v", err)
	}
}

func TestVentaRepository_Delete_Error(t *testing.T) {
	repo := &mockVentaRepository{
		DeleteFunc: func(id int64) error {
			return errors.New("error de BD")
		},
	}
	if err := repo.Delete(1); err == nil {
		t.Error("se esperaba error al eliminar")
	}
}

// ---------------------------------------------------------------------------
// Tests GetByID
// ---------------------------------------------------------------------------

func TestVentaRepository_GetByID_OK(t *testing.T) {
	esperada := &domain.Venta{IDVenta: 3, IDPasajero: 20, IDAsiento: 5, Serie: "F001", Correlativo: 10}
	repo := &mockVentaRepository{
		GetByIDFunc: func(id int64) (*domain.Venta, error) {
			if id == 3 {
				return esperada, nil
			}
			return nil, errors.New("no encontrado")
		},
	}
	venta, err := repo.GetByID(3)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if venta.Serie != "F001" || venta.Correlativo != 10 {
		t.Errorf("datos incorrectos: %+v", venta)
	}
}

func TestVentaRepository_GetByID_NotFound(t *testing.T) {
	repo := &mockVentaRepository{
		GetByIDFunc: func(id int64) (*domain.Venta, error) {
			return nil, errors.New("no encontrado")
		},
	}
	_, err := repo.GetByID(99)
	if err == nil {
		t.Error("se esperaba error al buscar id inexistente")
	}
}

// ---------------------------------------------------------------------------
// Tests List
// ---------------------------------------------------------------------------

func TestVentaRepository_List_OK(t *testing.T) {
	ventas := []domain.Venta{{IDVenta: 1}, {IDVenta: 2}, {IDVenta: 3}}
	repo := &mockVentaRepository{
		ListFunc: func(offset, limit int) ([]domain.Venta, int, error) {
			return ventas, len(ventas), nil
		},
	}
	result, total, err := repo.List(0, 15)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("esperado 3 ventas, obtenido %d", len(result))
	}
	if total != 3 {
		t.Errorf("esperado total 3, obtenido %d", total)
	}
}

func TestVentaRepository_List_Error(t *testing.T) {
	repo := &mockVentaRepository{
		ListFunc: func(offset, limit int) ([]domain.Venta, int, error) {
			return nil, 0, errors.New("error de BD")
		},
	}
	_, _, err := repo.List(0, 15)
	if err == nil {
		t.Error("se esperaba error al listar")
	}
}

// ---------------------------------------------------------------------------
// Tests NextCorrelativo
// ---------------------------------------------------------------------------

func TestVentaRepository_NextCorrelativo_PrimeroEnSerie(t *testing.T) {
	repo := &mockVentaRepository{
		NextCorrelativoFunc: func(serie string) (uint, error) {
			if serie == "B001" {
				return 1, nil
			}
			return 0, errors.New("serie no reconocida")
		},
	}
	correlativo, err := repo.NextCorrelativo("B001")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if correlativo != 1 {
		t.Errorf("esperado 1, obtenido %d", correlativo)
	}
}

func TestVentaRepository_NextCorrelativo_Secuencial(t *testing.T) {
	repo := &mockVentaRepository{
		NextCorrelativoFunc: func(serie string) (uint, error) {
			return 6, nil
		},
	}
	correlativo, err := repo.NextCorrelativo("B001")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if correlativo != 6 {
		t.Errorf("esperado 6, obtenido %d", correlativo)
	}
}

func TestVentaRepository_NextCorrelativo_IndependientePorSerie(t *testing.T) {
	contadores := map[string]uint{"B001": 3, "F001": 10, "T001": 1}
	repo := &mockVentaRepository{
		NextCorrelativoFunc: func(serie string) (uint, error) {
			v, ok := contadores[serie]
			if !ok {
				return 0, errors.New("serie inválida")
			}
			return v, nil
		},
	}
	casos := []struct {
		serie    string
		esperado uint
	}{
		{"B001", 3},
		{"F001", 10},
		{"T001", 1},
	}
	for _, c := range casos {
		got, err := repo.NextCorrelativo(c.serie)
		if err != nil {
			t.Errorf("serie %s: no se esperaba error: %v", c.serie, err)
		}
		if got != c.esperado {
			t.Errorf("serie %s: esperado %d, obtenido %d", c.serie, c.esperado, got)
		}
	}
}

func TestVentaRepository_NextCorrelativo_Error(t *testing.T) {
	repo := &mockVentaRepository{
		NextCorrelativoFunc: func(serie string) (uint, error) {
			return 0, errors.New("error de BD")
		},
	}
	_, err := repo.NextCorrelativo("B001")
	if err == nil {
		t.Error("se esperaba error en NextCorrelativo")
	}
}
