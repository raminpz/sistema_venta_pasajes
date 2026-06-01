package service

import (
	"errors"
	"sistema_venta_pasajes/internal/tipo_vehiculo/domain"
	"sistema_venta_pasajes/internal/tipo_vehiculo/input"
	"testing"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type fakeTipoVehiculoRepo struct {
	createFn  func(*domain.TipoVehiculo) error
	updateFn  func(*domain.TipoVehiculo) error
	deleteFn  func(int64) error
	getByIDFn func(int64) (*domain.TipoVehiculo, error)
	listFn    func() ([]domain.TipoVehiculo, error)
}

func (f *fakeTipoVehiculoRepo) Create(tv *domain.TipoVehiculo) error {
	if f.createFn != nil {
		return f.createFn(tv)
	}
	return nil
}
func (f *fakeTipoVehiculoRepo) Update(tv *domain.TipoVehiculo) error {
	if f.updateFn != nil {
		return f.updateFn(tv)
	}
	return nil
}
func (f *fakeTipoVehiculoRepo) Delete(id int64) error {
	if f.deleteFn != nil {
		return f.deleteFn(id)
	}
	return nil
}
func (f *fakeTipoVehiculoRepo) GetByID(id int64) (*domain.TipoVehiculo, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeTipoVehiculoRepo) List() ([]domain.TipoVehiculo, error) {
	if f.listFn != nil {
		return f.listFn()
	}
	return []domain.TipoVehiculo{}, nil
}

func TestServiceCreateOK(t *testing.T) {
	repo := &fakeTipoVehiculoRepo{createFn: func(tv *domain.TipoVehiculo) error {
		tv.IDTipoVehiculo = 1
		return nil
	}}
	s := NewTipoVehiculoService(repo)
	out, err := s.Create(input.CreateTipoVehiculoInput{Nombre: "camioneta", Descripcion: "ruta corta"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), out.IDTipoVehiculo)
	assert.Equal(t, "Camioneta", out.Nombre)
}

func TestServiceCreateValidationError(t *testing.T) {
	s := NewTipoVehiculoService(&fakeTipoVehiculoRepo{})
	_, err := s.Create(input.CreateTipoVehiculoInput{Nombre: "", Descripcion: "desc"})
	assert.Error(t, err)
}

func TestServiceCreateDuplicate(t *testing.T) {
	repo := &fakeTipoVehiculoRepo{createFn: func(tv *domain.TipoVehiculo) error {
		return &mysqlDriver.MySQLError{Number: 1062, Message: "Duplicate entry 'AUTO' for key 'NOMBRE'"}
	}}
	s := NewTipoVehiculoService(repo)
	_, err := s.Create(input.CreateTipoVehiculoInput{Nombre: "AUTO", Descripcion: "desc"})
	assert.Error(t, err)
}

func TestServiceCreateDataTruncatedExactMessage(t *testing.T) {
	repo := &fakeTipoVehiculoRepo{createFn: func(tv *domain.TipoVehiculo) error {
		return &mysqlDriver.MySQLError{Number: 1265, Message: "Data truncated for column 'NOMBRE' at row 1"}
	}}
	s := NewTipoVehiculoService(repo)
	_, err := s.Create(input.CreateTipoVehiculoInput{Nombre: "Bus", Descripcion: "Bus de alto tonelaje"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Data truncated for column 'NOMBRE' at row 1")
}

func TestServiceUpdateOK(t *testing.T) {
	nombre := "Mini bus"
	desc := "interprovincial"
	repo := &fakeTipoVehiculoRepo{
		getByIDFn: func(id int64) (*domain.TipoVehiculo, error) {
			return &domain.TipoVehiculo{IDTipoVehiculo: id, Nombre: "Combi", Descripcion: "ant"}, nil
		},
		updateFn: func(tv *domain.TipoVehiculo) error {
			assert.Equal(t, "Mini Bus", tv.Nombre)
			assert.Equal(t, "Interprovincial", tv.Descripcion)
			return nil
		},
	}
	s := NewTipoVehiculoService(repo)
	out, err := s.Update(1, input.UpdateTipoVehiculoInput{Nombre: &nombre, Descripcion: &desc})
	assert.NoError(t, err)
	assert.Equal(t, "Mini Bus", out.Nombre)
}

func TestServiceUpdateInvalidID(t *testing.T) {
	s := NewTipoVehiculoService(&fakeTipoVehiculoRepo{})
	_, err := s.Update(0, input.UpdateTipoVehiculoInput{})
	assert.Error(t, err)
}

func TestServiceUpdateValidation(t *testing.T) {
	s := NewTipoVehiculoService(&fakeTipoVehiculoRepo{})
	_, err := s.Update(1, input.UpdateTipoVehiculoInput{})
	assert.Error(t, err)
}

func TestServiceUpdateNotFound(t *testing.T) {
	nombre := "Auto"
	repo := &fakeTipoVehiculoRepo{getByIDFn: func(int64) (*domain.TipoVehiculo, error) {
		return nil, gorm.ErrRecordNotFound
	}}
	s := NewTipoVehiculoService(repo)
	_, err := s.Update(10, input.UpdateTipoVehiculoInput{Nombre: &nombre})
	assert.Error(t, err)
}

func TestServiceDeleteOK(t *testing.T) {
	repo := &fakeTipoVehiculoRepo{
		getByIDFn: func(id int64) (*domain.TipoVehiculo, error) {
			return &domain.TipoVehiculo{IDTipoVehiculo: id, Nombre: "Auto", Descripcion: "desc"}, nil
		},
		deleteFn: func(int64) error { return nil },
	}
	s := NewTipoVehiculoService(repo)
	err := s.Delete(1)
	assert.NoError(t, err)
}

func TestServiceDeleteConflict(t *testing.T) {
	repo := &fakeTipoVehiculoRepo{
		getByIDFn: func(id int64) (*domain.TipoVehiculo, error) {
			return &domain.TipoVehiculo{IDTipoVehiculo: id, Nombre: "Auto", Descripcion: "desc"}, nil
		},
		deleteFn: func(int64) error {
			return &mysqlDriver.MySQLError{Number: 1451, Message: "Cannot delete or update a parent row"}
		},
	}
	s := NewTipoVehiculoService(repo)
	err := s.Delete(1)
	assert.Error(t, err)
}

func TestServiceDeleteGetByIDError(t *testing.T) {
	repo := &fakeTipoVehiculoRepo{getByIDFn: func(int64) (*domain.TipoVehiculo, error) {
		return nil, errors.New("db down")
	}}
	s := NewTipoVehiculoService(repo)
	err := s.Delete(1)
	assert.Error(t, err)
}

func TestServiceGetByIDOK(t *testing.T) {
	repo := &fakeTipoVehiculoRepo{getByIDFn: func(id int64) (*domain.TipoVehiculo, error) {
		return &domain.TipoVehiculo{IDTipoVehiculo: id, Nombre: "Auto", Descripcion: "desc"}, nil
	}}
	s := NewTipoVehiculoService(repo)
	out, err := s.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), out.IDTipoVehiculo)
}

func TestServiceGetByIDInvalid(t *testing.T) {
	s := NewTipoVehiculoService(&fakeTipoVehiculoRepo{})
	_, err := s.GetByID(0)
	assert.Error(t, err)
}

func TestServiceListOK(t *testing.T) {
	repo := &fakeTipoVehiculoRepo{listFn: func() ([]domain.TipoVehiculo, error) {
		return []domain.TipoVehiculo{
			{IDTipoVehiculo: 1, Nombre: "Auto", Descripcion: "desc"},
			{IDTipoVehiculo: 2, Nombre: "Combi", Descripcion: "desc"},
		}, nil
	}}
	s := NewTipoVehiculoService(repo)
	out, err := s.List()
	assert.NoError(t, err)
	assert.Len(t, out, 2)
}

func TestServiceListError(t *testing.T) {
	repo := &fakeTipoVehiculoRepo{listFn: func() ([]domain.TipoVehiculo, error) {
		return nil, errors.New("db")
	}}
	s := NewTipoVehiculoService(repo)
	_, err := s.List()
	assert.Error(t, err)
}
