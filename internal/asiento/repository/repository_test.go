package repository

import (
	"errors"
	"sistema_venta_pasajes/internal/asiento/domain"
	"testing"
)

// mockAsientoDB simula un repositorio en memoria para pruebas
type mockAsientoDB struct {
	asientos map[int64]*domain.Asiento
	nextID   int64
}

func newMockAsientoDB() *mockAsientoDB {
	return &mockAsientoDB{
		asientos: make(map[int64]*domain.Asiento),
		nextID:   1,
	}
}

func (m *mockAsientoDB) Create(a *domain.Asiento) error {
	a.IDAsiento = int(m.nextID)
	cp := *a
	m.asientos[m.nextID] = &cp
	m.nextID++
	return nil
}

func (m *mockAsientoDB) GetByID(id int64) (*domain.Asiento, error) {
	a, ok := m.asientos[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return a, nil
}

func (m *mockAsientoDB) ListByVehiculo(idVehiculo int64) ([]*domain.Asiento, error) {
	var result []*domain.Asiento
	for _, a := range m.asientos {
		if int64(a.IDVehiculo) == idVehiculo {
			cp := *a
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (m *mockAsientoDB) Update(a *domain.Asiento) error {
	if _, ok := m.asientos[int64(a.IDAsiento)]; !ok {
		return errors.New("record not found")
	}
	cp := *a
	m.asientos[int64(a.IDAsiento)] = &cp
	return nil
}

func (m *mockAsientoDB) Delete(id int64) error {
	if _, ok := m.asientos[id]; !ok {
		return errors.New("record not found")
	}
	delete(m.asientos, id)
	return nil
}

func (m *mockAsientoDB) CambiarEstado(id int64, estado string) error {
	a, ok := m.asientos[id]
	if !ok {
		return errors.New("record not found")
	}
	a.Estado = estado
	return nil
}

// ---- Tests ----

func TestAsientoRepository_Create(t *testing.T) {
	db := newMockAsientoDB()
	a := &domain.Asiento{IDVehiculo: 1, NumeroAsiento: "A1", Estado: "ACTIVO"}
	err := db.Create(a)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if a.IDAsiento == 0 {
		t.Error("IDAsiento debe ser asignado después de Create")
	}
}

func TestAsientoRepository_GetByID(t *testing.T) {
	db := newMockAsientoDB()
	a := &domain.Asiento{IDVehiculo: 2, NumeroAsiento: "B2", Estado: "RESERVADO"}
	if err := db.Create(a); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := db.GetByID(int64(a.IDAsiento))
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.IDAsiento != a.IDAsiento {
		t.Errorf("GetByID() ID esperado %d, obtenido %d", a.IDAsiento, got.IDAsiento)
	}
	if got.Estado != "RESERVADO" {
		t.Errorf("GetByID() Estado esperado RESERVADO, obtenido %s", got.Estado)
	}

	_, err = db.GetByID(9999)
	if err == nil {
		t.Error("GetByID() debería retornar error para ID inexistente")
	}
}

func TestAsientoRepository_ListByVehiculo(t *testing.T) {
	db := newMockAsientoDB()
	if err := db.Create(&domain.Asiento{IDVehiculo: 2, NumeroAsiento: "A1", Estado: "ACTIVO"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := db.Create(&domain.Asiento{IDVehiculo: 2, NumeroAsiento: "A2", Estado: "OCUPADO"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := db.Create(&domain.Asiento{IDVehiculo: 5, NumeroAsiento: "C1", Estado: "ACTIVO"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	asientos, err := db.ListByVehiculo(2)
	if err != nil {
		t.Fatalf("ListByVehiculo() error = %v", err)
	}
	if len(asientos) != 2 {
		t.Errorf("ListByVehiculo() esperaba 2 asientos, obtuvo %d", len(asientos))
	}
}

func TestAsientoRepository_Update(t *testing.T) {
	db := newMockAsientoDB()
	a := &domain.Asiento{IDVehiculo: 1, NumeroAsiento: "C3", Estado: "ACTIVO"}
	if err := db.Create(a); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	a.NumeroAsiento = "C4"
	a.Estado = "OCUPADO"
	err := db.Update(a)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, _ := db.GetByID(int64(a.IDAsiento))
	if got.NumeroAsiento != "C4" {
		t.Errorf("Update() NumeroAsiento esperado C4, obtenido %s", got.NumeroAsiento)
	}
	if got.Estado != "OCUPADO" {
		t.Errorf("Update() Estado esperado OCUPADO, obtenido %s", got.Estado)
	}
}

func TestAsientoRepository_Delete(t *testing.T) {
	db := newMockAsientoDB()
	a := &domain.Asiento{IDVehiculo: 1, NumeroAsiento: "D1", Estado: "ACTIVO"}
	if err := db.Create(a); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err := db.Delete(int64(a.IDAsiento))
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = db.GetByID(int64(a.IDAsiento))
	if err == nil {
		t.Error("GetByID() debería retornar error después de Delete")
	}

	err = db.Delete(9999)
	if err == nil {
		t.Error("Delete() debería retornar error para ID inexistente")
	}
}

func TestAsientoRepository_CambiarEstado(t *testing.T) {
	db := newMockAsientoDB()
	a := &domain.Asiento{IDVehiculo: 1, NumeroAsiento: "F1", Estado: "ACTIVO"}
	if err := db.Create(a); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	estados := []string{"RESERVADO", "OCUPADO", "ACTIVO"}
	for _, estado := range estados {
		err := db.CambiarEstado(int64(a.IDAsiento), estado)
		if err != nil {
			t.Errorf("CambiarEstado(%s) error = %v", estado, err)
		}
		got, _ := db.GetByID(int64(a.IDAsiento))
		if got.Estado != estado {
			t.Errorf("CambiarEstado() Estado esperado %s, obtenido %s", estado, got.Estado)
		}
	}

	err := db.CambiarEstado(9999, "ACTIVO")
	if err == nil {
		t.Error("CambiarEstado() debería retornar error para ID inexistente")
	}
}

