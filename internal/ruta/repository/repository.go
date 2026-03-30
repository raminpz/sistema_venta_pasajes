package repository

import (
	"sistema_venta_pasajes/internal/ruta/domain"

	"gorm.io/gorm"
)

type RutaRepository interface {
	Create(ruta *domain.Ruta) error
	GetByID(id int) (*domain.Ruta, error)
	Update(ruta *domain.Ruta) error
	Delete(id int) error
	List() ([]domain.Ruta, error)
}

type rutaRepository struct {
	db *gorm.DB
}

func NewRutaRepository(db *gorm.DB) RutaRepository {
	return &rutaRepository{db: db}
}

func (r *rutaRepository) Create(ruta *domain.Ruta) error {
	return r.db.Create(ruta).Error
}

func (r *rutaRepository) GetByID(id int) (*domain.Ruta, error) {
	var ruta domain.Ruta
	if err := r.db.First(&ruta, "id_ruta = ?", id).Error; err != nil {
		return nil, err
	}
	return &ruta, nil
}

func (r *rutaRepository) Update(ruta *domain.Ruta) error {
	return r.db.Save(ruta).Error
}

func (r *rutaRepository) Delete(id int) error {
	res := r.db.Delete(&domain.Ruta{}, "id_ruta = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *rutaRepository) List() ([]domain.Ruta, error) {
	var rutas []domain.Ruta
	err := r.db.Model(&domain.Ruta{}).Order("id_ruta ASC").Find(&rutas).Error
	return rutas, err
}
