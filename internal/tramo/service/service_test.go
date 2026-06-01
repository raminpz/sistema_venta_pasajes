package service

import (
	"sistema_venta_pasajes/internal/tramo/domain"
	"sistema_venta_pasajes/internal/tramo/input"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockTramoRepo struct {
	tramos map[int64]*domain.Tramo
	nextID  int64
}

func newMockRepo() *mockTramoRepo {
	return &mockTramoRepo{tramos: make(map[int64]*domain.Tramo), nextID: 1}
}

func (m *mockTramoRepo) Create(t *domain.Tramo) error {
	t.IDTramo = m.nextID
	m.nextID++
	m.tramos[t.IDTramo] = t
	return nil
}
func (m *mockTramoRepo) Update(t *domain.Tramo) error {
	if _, ok := m.tramos[t.IDTramo]; !ok {
		return gorm.ErrRecordNotFound
	}
	m.tramos[t.IDTramo] = t
	return nil
}
func (m *mockTramoRepo) Delete(id int64) error {
	if _, ok := m.tramos[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.tramos, id)
	return nil
}
func (m *mockTramoRepo) GetByID(id int64) (*domain.Tramo, error) {
	t, ok := m.tramos[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return t, nil
}
func (m *mockTramoRepo) List(offset, limit int) ([]domain.Tramo, int, error) {
	var list []domain.Tramo
	for _, t := range m.tramos {
		list = append(list, *t)
	}
	return list, len(list), nil
}
func (m *mockTramoRepo) ListByRuta(idRuta int64) ([]domain.Tramo, error) {
	var list []domain.Tramo
	for _, t := range m.tramos {
		if t.IDRuta == idRuta {
			list = append(list, *t)
		}
	}
	return list, nil
}
func (m *mockTramoRepo) ExistsByRutaParadas(idRuta, idOrigen, idDestino int64) (bool, error) {
	for _, t := range m.tramos {
		if t.IDRuta == idRuta && t.IDParadaOrigen == idOrigen && t.IDParadaDestino == idDestino {
			return true, nil
		}
	}
	return false, nil
}

type mockTramoRepoFn struct {
	createFn            func(*domain.Tramo) error
	updateFn            func(*domain.Tramo) error
	deleteFn            func(int64) error
	getByIDFn           func(int64) (*domain.Tramo, error)
	listFn              func(int, int) ([]domain.Tramo, int, error)
	listByRutaFn        func(int64) ([]domain.Tramo, error)
	existsByRutaParFn   func(int64, int64, int64) (bool, error)
}

func (m *mockTramoRepoFn) Create(t *domain.Tramo) error {
	if m.createFn != nil {
		return m.createFn(t)
	}
	return nil
}
func (m *mockTramoRepoFn) Update(t *domain.Tramo) error {
	if m.updateFn != nil {
		return m.updateFn(t)
	}
	return nil
}
func (m *mockTramoRepoFn) Delete(id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}
func (m *mockTramoRepoFn) GetByID(id int64) (*domain.Tramo, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockTramoRepoFn) List(offset, limit int) ([]domain.Tramo, int, error) {
	if m.listFn != nil {
		return m.listFn(offset, limit)
	}
	return []domain.Tramo{}, 0, nil
}
func (m *mockTramoRepoFn) ListByRuta(idRuta int64) ([]domain.Tramo, error) {
	if m.listByRutaFn != nil {
		return m.listByRutaFn(idRuta)
	}
	return []domain.Tramo{}, nil
}
func (m *mockTramoRepoFn) ExistsByRutaParadas(idRuta, idOrigen, idDestino int64) (bool, error) {
	if m.existsByRutaParFn != nil {
		return m.existsByRutaParFn(idRuta, idOrigen, idDestino)
	}
	return false, nil
}

func TestCreate_OK(t *testing.T) {
	svc := NewTramoService(newMockRepo())
	out, err := svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 2})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), out.IDRuta)
}

func TestCreate_SinRuta(t *testing.T) {
	svc := NewTramoService(newMockRepo())
	_, err := svc.Create(input.CreateTramoInput{IDRuta: 0, IDParadaOrigen: 1, IDParadaDestino: 2})
	assert.Error(t, err)
}

func TestCreate_ParadasIguales(t *testing.T) {
	svc := NewTramoService(newMockRepo())
	_, err := svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 2, IDParadaDestino: 2})
	assert.Error(t, err)
}

func TestCreate_ParadaOrigenInvalida(t *testing.T) {
	svc := NewTramoService(newMockRepo())
	_, err := svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 0, IDParadaDestino: 2})
	assert.Error(t, err)
}

func TestCreate_Duplicado(t *testing.T) {
	repo := newMockRepo()
	svc := NewTramoService(repo)
	in := input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 2}
	_, _ = svc.Create(in)
	_, err := svc.Create(in)
	assert.Error(t, err)
}

func TestGetByID_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewTramoService(repo)
	_, _ = svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 3})
	out, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), out.IDTramo)
}

func TestGetByID_NotFound(t *testing.T) {
	svc := NewTramoService(newMockRepo())
	_, err := svc.GetByID(999)
	assert.Error(t, err)
}

func TestUpdate_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewTramoService(repo)
	_, _ = svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 2})
	nuevaParadaDestino := int64(3)
	out, err := svc.Update(input.UpdateTramoInput{IDTramo: 1, IDParadaDestino: &nuevaParadaDestino})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), out.IDParadaDestino)
}

func TestUpdate_EmptyFields(t *testing.T) {
	svc := NewTramoService(newMockRepo())
	_, err := svc.Update(input.UpdateTramoInput{IDTramo: 1})
	assert.Error(t, err)
}

func TestUpdate_ExistsCheckError(t *testing.T) {
	svc := NewTramoService(&mockTramoRepoFn{existsByRutaParFn: func(int64, int64, int64) (bool, error) {
		return false, assert.AnError
	}})
	_, err := svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 2})
	assert.Error(t, err)
}

func TestUpdate_NotFound(t *testing.T) {
	svc := NewTramoService(&mockTramoRepoFn{getByIDFn: func(int64) (*domain.Tramo, error) {
		return nil, gorm.ErrRecordNotFound
	}})
	ruta := int64(1)
	origen := int64(1)
	destino := int64(2)
	_, err := svc.Update(input.UpdateTramoInput{IDTramo: 1, IDRuta: &ruta, IDParadaOrigen: &origen, IDParadaDestino: &destino})
	assert.Error(t, err)
}

func TestUpdate_GetByIDInternalError(t *testing.T) {
	svc := NewTramoService(&mockTramoRepoFn{getByIDFn: func(int64) (*domain.Tramo, error) {
		return nil, assert.AnError
	}})
	ruta := int64(1)
	origen := int64(1)
	destino := int64(2)
	_, err := svc.Update(input.UpdateTramoInput{IDTramo: 1, IDRuta: &ruta, IDParadaOrigen: &origen, IDParadaDestino: &destino})
	assert.Error(t, err)
}

func TestUpdate_ParadasIgualesDespuesDeMerge(t *testing.T) {
	svc := NewTramoService(&mockTramoRepoFn{getByIDFn: func(int64) (*domain.Tramo, error) {
		return &domain.Tramo{IDTramo: 1, IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 2}, nil
	}})
	nueva := int64(1)
	_, err := svc.Update(input.UpdateTramoInput{IDTramo: 1, IDParadaDestino: &nueva})
	assert.Error(t, err)
}

func TestDelete_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewTramoService(repo)
	_, _ = svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 2})
	err := svc.Delete(1)
	assert.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	svc := NewTramoService(newMockRepo())
	err := svc.Delete(999)
	assert.Error(t, err)
}

func TestDelete_GetByIDInternalError(t *testing.T) {
	svc := NewTramoService(&mockTramoRepoFn{getByIDFn: func(int64) (*domain.Tramo, error) {
		return nil, assert.AnError
	}})
	err := svc.Delete(1)
	assert.Error(t, err)
}

func TestDelete_DeleteNotFound(t *testing.T) {
	svc := NewTramoService(&mockTramoRepoFn{
		getByIDFn: func(int64) (*domain.Tramo, error) { return &domain.Tramo{IDTramo: 1}, nil },
		deleteFn:  func(int64) error { return gorm.ErrRecordNotFound },
	})
	err := svc.Delete(1)
	assert.Error(t, err)
}

func TestListByRuta_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewTramoService(repo)
	_, _ = svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 2})
	_, _ = svc.Create(input.CreateTramoInput{IDRuta: 2, IDParadaOrigen: 3, IDParadaDestino: 4})
	list, err := svc.ListByRuta(1)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestList_OK(t *testing.T) {
	repo := newMockRepo()
	svc := NewTramoService(repo)
	_, _ = svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 1, IDParadaDestino: 2})
	_, _ = svc.Create(input.CreateTramoInput{IDRuta: 1, IDParadaOrigen: 2, IDParadaDestino: 3})
	out, total, err := svc.List(1, 15)
	assert.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, 2, total)
}

func TestList_Error(t *testing.T) {
	svc := NewTramoService(&mockTramoRepoFn{listFn: func(int, int) ([]domain.Tramo, int, error) {
		return nil, 0, assert.AnError
	}})
	_, _, err := svc.List(1, 10)
	assert.Error(t, err)
}

func TestListByRuta_Error(t *testing.T) {
	svc := NewTramoService(&mockTramoRepoFn{listByRutaFn: func(int64) ([]domain.Tramo, error) {
		return nil, assert.AnError
	}})
	_, err := svc.ListByRuta(1)
	assert.Error(t, err)
}
