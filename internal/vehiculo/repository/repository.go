package repository

import (
	"sistema_venta_pasajes/internal/vehiculo/domain"

	"gorm.io/gorm"
)

type VehiculoRepository interface {
	Create(vehiculo *domain.Vehiculo) error
	Update(vehiculo *domain.Vehiculo) error
	Delete(id int64) error
	GetByID(id int64) (*domain.Vehiculo, error)
	List(offset, limit int) ([]domain.Vehiculo, int, error)
	ExistsByPlaca(placa string) (bool, error)
}

type vehiculoRepository struct {
	db *gorm.DB
}

func NewVehiculoRepository(db *gorm.DB) VehiculoRepository {
	return &vehiculoRepository{db: db}
}

func (r *vehiculoRepository) Create(vehiculo *domain.Vehiculo) error {
	return r.db.Create(vehiculo).Error
}

func (r *vehiculoRepository) Update(vehiculo *domain.Vehiculo) error {
	return r.db.Save(vehiculo).Error
}

func (r *vehiculoRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.Vehiculo{}, "ID_VEHICULO = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *vehiculoRepository) GetByID(id int64) (*domain.Vehiculo, error) {
	var v domain.Vehiculo
	err := r.db.First(&v, id).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *vehiculoRepository) List(offset, limit int) ([]domain.Vehiculo, int, error) {
	var vehiculos []domain.Vehiculo
	var total int64
	db := r.db.Model(&domain.Vehiculo{})
	db.Count(&total)
	db.Order("ID_VEHICULO ASC").Offset(offset).Limit(limit).Find(&vehiculos)
	return vehiculos, int(total), db.Error
}

func (r *vehiculoRepository) ExistsByPlaca(placa string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Vehiculo{}).
		Where("NRO_PLACA = ?", placa).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
