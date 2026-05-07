package service

import (
	"sistema_venta_pasajes/internal/parada/domain"
	"sistema_venta_pasajes/internal/parada/input"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// ── mock ─────────────────────────────────────────────────────────────────────

type mockParadaRepo struct {
	paradas map[int64]*domain.Parada
	nextID  int64
}

func newMock() *mockParadaRepo {
	return &mockParadaRepo{paradas: make(map[int64]*domain.Parada), nextID: 1}
}

func (m *mockParadaRepo) Create(p *domain.Parada) error {
	p.IDParada = m.nextID
	m.nextID++
	m.paradas[p.IDParada] = p
	return nil
}
func (m *mockParadaRepo) Update(p *domain.Parada) error {
	if _, ok := m.paradas[p.IDParada]; !ok {
		return gorm.ErrRecordNotFound
	}
	m.paradas[p.IDParada] = p
	return nil
}
func (m *mockParadaRepo) Delete(id int64) error {
	if _, ok := m.paradas[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.paradas, id)
	return nil
}
func (m *mockParadaRepo) GetByID(id int64) (*domain.Parada, error) {
	p, ok := m.paradas[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return p, nil
}
func (m *mockParadaRepo) ListByRuta(idRuta int64) ([]domain.Parada, error) {
	var list []domain.Parada
	for _, p := range m.paradas {
		if p.IDRuta == idRuta {
			list = append(list, *p)
		}
	}
	return list, nil
}
func (m *mockParadaRepo) ExistsByRutaNombre(idRuta int64, nombreParada string) (bool, error) {
	for _, p := range m.paradas {
		if p.IDRuta == idRuta && p.NombreParada == nombreParada {
			return true, nil
		}
	}
	return false, nil
}
func (m *mockParadaRepo) ExistsByRutaOrden(idRuta int64, orden int) (bool, error) {
	for _, p := range m.paradas {
		if p.IDRuta == idRuta && p.Orden == orden {
			return true, nil
		}
	}
	return false, nil
}
func (m *mockParadaRepo) GetOrdenByID(idParada int64) (int, error) {
	p, ok := m.paradas[idParada]
	if !ok {
		return 0, gorm.ErrRecordNotFound
	}
	return p.Orden, nil
}

// ── tests ─────────────────────────────────────────────────────────────────────

func TestParada_Create_OK(t *testing.T) {
	svc := NewParadaService(newMock())
	out, err := svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Huanta", Orden: 1})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), out.IDParada)
}

func TestParada_Create_SinRuta(t *testing.T) {
	svc := NewParadaService(newMock())
	_, err := svc.Create(input.CreateParadaInput{IDRuta: 0, NombreParada: "Huanta", Orden: 1})
	assert.Error(t, err)
}

func TestParada_Create_SinNombreParada(t *testing.T) {
	svc := NewParadaService(newMock())
	_, err := svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "", Orden: 1})
	assert.Error(t, err)
}

func TestParada_Create_OrdenCero(t *testing.T) {
	svc := NewParadaService(newMock())
	_, err := svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Huanta", Orden: 0})
	assert.Error(t, err)
}

func TestParada_Create_DuplicadoNombre(t *testing.T) {
	repo := newMock()
	svc := NewParadaService(repo)
	_, _ = svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Huanta", Orden: 1})
	_, err := svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Huanta", Orden: 2})
	assert.Error(t, err)
}

func TestParada_Create_DuplicadoOrden(t *testing.T) {
	repo := newMock()
	svc := NewParadaService(repo)
	_, _ = svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Huanta", Orden: 1})
	_, err := svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Izcuchaca", Orden: 1})
	assert.Error(t, err)
}

func TestParada_GetByID_OK(t *testing.T) {
	repo := newMock()
	svc := NewParadaService(repo)
	_, _ = svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Huanta", Orden: 1})
	out, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), out.IDParada)
	assert.Equal(t, 1, out.Orden)
}

func TestParada_GetByID_NotFound(t *testing.T) {
	svc := NewParadaService(newMock())
	_, err := svc.GetByID(999)
	assert.Error(t, err)
}

func TestParada_Update_OK(t *testing.T) {
	repo := newMock()
	svc := NewParadaService(repo)
	_, _ = svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Huanta", Orden: 1})
	nombreParada := "Huamanga"
	out, err := svc.Update(input.UpdateParadaInput{IDParada: 1, NombreParada: &nombreParada})
	assert.NoError(t, err)
	assert.Equal(t, "Huamanga", out.NombreParada)
}

func TestParada_Update_EmptyFields(t *testing.T) {
	svc := NewParadaService(newMock())
	_, err := svc.Update(input.UpdateParadaInput{IDParada: 1})
	assert.Error(t, err)
}

func TestParada_Delete_OK(t *testing.T) {
	repo := newMock()
	svc := NewParadaService(repo)
	_, _ = svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Huanta", Orden: 1})
	err := svc.Delete(1)
	assert.NoError(t, err)
}

func TestParada_Delete_NotFound(t *testing.T) {
	svc := NewParadaService(newMock())
	err := svc.Delete(999)
	assert.Error(t, err)
}

func TestParada_ListByRuta_OK(t *testing.T) {
	repo := newMock()
	svc := NewParadaService(repo)
	_, _ = svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Huanta", Orden: 1})
	_, _ = svc.Create(input.CreateParadaInput{IDRuta: 1, NombreParada: "Izcuchaca", Orden: 2})
	_, _ = svc.Create(input.CreateParadaInput{IDRuta: 2, NombreParada: "Huancayo", Orden: 1})
	list, err := svc.ListByRuta(1)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestParada_ListByRuta_Vacia(t *testing.T) {
	svc := NewParadaService(newMock())
	list, err := svc.ListByRuta(99)
	assert.NoError(t, err)
	assert.Empty(t, list)
}
