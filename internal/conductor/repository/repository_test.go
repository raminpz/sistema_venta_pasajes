
package repository

import (
	"errors"
	"sistema_venta_pasajes/internal/conductor/domain"
	"testing"
)

type mockDB struct {
	conductores map[int64]*domain.Conductor
	nextID     int64
}

func newMockDB() *mockDB {
	return &mockDB{
		conductores: make(map[int64]*domain.Conductor),
		nextID:     1,
	}
}

func (m *mockDB) Create(conductor *domain.Conductor) error {
	conductor.IDConductor = m.nextID
	m.conductores[m.nextID] = conductor
	m.nextID++
	return nil
}

func (m *mockDB) GetByID(id int64) (*domain.Conductor, error) {
	c, ok := m.conductores[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (m *mockDB) Update(conductor *domain.Conductor) error {
	if conductor.IDConductor == 0 {
		return errors.New("not found")
	}
	m.conductores[conductor.IDConductor] = conductor
	return nil
}

func (m *mockDB) Delete(id int64) error {
	if _, ok := m.conductores[id]; !ok {
		return errors.New("not found")
	}
	delete(m.conductores, id)
	return nil
}

func (m *mockDB) List() ([]domain.Conductor, error) {
	var list []domain.Conductor
	for _, c := range m.conductores {
		list = append(list, *c)
	}
	return list, nil
}

func TestMockConductorRepository_CRUD(t *testing.T) {
	db := newMockDB()
	c := &domain.Conductor{
		Nombres:        "Juan",
		Apellidos:      "Perez",
		NumeroLicencia: "ABC123456",
		Telefono:       "987654321",
	}
	// Create
	err := db.Create(c)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if c.IDConductor == 0 {
		t.Error("IDConductor not set after create")
	}
	// GetByID
	got, err := db.GetByID(c.IDConductor)
	if err != nil || got == nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	// Update
	c.Nombres = "Juan Carlos"
	err = db.Update(c)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	got, _ = db.GetByID(c.IDConductor)
	if got.Nombres != "Juan Carlos" {
		t.Error("Update did not persist changes")
	}
	// List
	list, err := db.List()
	if err != nil || len(list) == 0 {
		t.Errorf("List failed: %v", err)
	}
	// Delete
	err = db.Delete(c.IDConductor)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	got, err = db.GetByID(c.IDConductor)
	if got != nil {
		t.Error("Delete did not remove record")
	}
}
