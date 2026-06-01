package repository

import (
	"sistema_venta_pasajes/internal/tipo_vehiculo/domain"

	"gorm.io/gorm"
)

type TipoVehiculoRepository interface {
	Create(tv *domain.TipoVehiculo) error
	Update(tv *domain.TipoVehiculo) error
	Delete(id int64) error
	GetByID(id int64) (*domain.TipoVehiculo, error)
	List() ([]domain.TipoVehiculo, error)
}

type tipoVehiculoRepository struct {
	db *gorm.DB
}

func NewTipoVehiculoRepository(db *gorm.DB) TipoVehiculoRepository {
	return &tipoVehiculoRepository{db: db}
}

func (r *tipoVehiculoRepository) Create(tv *domain.TipoVehiculo) error {
	return r.db.Create(tv).Error
}

func (r *tipoVehiculoRepository) Update(tv *domain.TipoVehiculo) error {
	return r.db.Save(tv).Error
}

func (r *tipoVehiculoRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.TipoVehiculo{}, "ID_TIPO_VEHICULO = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *tipoVehiculoRepository) GetByID(id int64) (*domain.TipoVehiculo, error) {
	var tv domain.TipoVehiculo
	if err := r.db.First(&tv, "ID_TIPO_VEHICULO = ?", id).Error; err != nil {
		return nil, err
	}
	return &tv, nil
}

func (r *tipoVehiculoRepository) List() ([]domain.TipoVehiculo, error) {
	var items []domain.TipoVehiculo
	err := r.db.Model(&domain.TipoVehiculo{}).Order("ID_TIPO_VEHICULO ASC").Find(&items).Error
	return items, err
}

