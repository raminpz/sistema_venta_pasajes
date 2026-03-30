package repository

import (
	"sistema_venta_pasajes/internal/conductor/domain"
	"gorm.io/gorm"
)

type ConductorRepository interface {
	Create(conductor *domain.Conductor) error
	GetByID(id int64) (*domain.Conductor, error)
	Update(conductor *domain.Conductor) error
	Delete(id int64) error
	List() ([]domain.Conductor, error)
}

type conductorRepository struct {
	db *gorm.DB
}

func NewConductorRepository(db *gorm.DB) ConductorRepository {
	return &conductorRepository{db: db}
}

func (r *conductorRepository) Create(conductor *domain.Conductor) error {
	return r.db.Create(conductor).Error
}

func (r *conductorRepository) GetByID(id int64) (*domain.Conductor, error) {
	var conductor domain.Conductor
	if err := r.db.First(&conductor, "id_conductor = ?", id).Error; err != nil {
		return nil, err
	}
	return &conductor, nil
}

func (r *conductorRepository) Update(conductor *domain.Conductor) error {
	return r.db.Save(conductor).Error
}

func (r *conductorRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.Conductor{}, "id_conductor = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *conductorRepository) List() ([]domain.Conductor, error) {
	var conductores []domain.Conductor
	err := r.db.Model(&domain.Conductor{}).Order("id_conductor ASC").Find(&conductores).Error
	return conductores, err
}
