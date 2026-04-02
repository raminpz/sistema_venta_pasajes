package repository

import (
	"sistema_venta_pasajes/internal/encomienda/domain"

	"gorm.io/gorm"
)

type EncomiendaRepository interface {
	Create(encomienda *domain.Encomienda) error
	Update(encomienda *domain.Encomienda) error
	Delete(id int64) error
	GetByID(id int64) (*domain.Encomienda, error)
	List(offset, limit int) ([]domain.Encomienda, int, error)
}

type encomiendaRepository struct {
	db *gorm.DB
}

func NewEncomiendaRepository(db *gorm.DB) EncomiendaRepository {
	return &encomiendaRepository{db: db}
}

func (r *encomiendaRepository) Create(encomienda *domain.Encomienda) error {
	return r.db.Create(encomienda).Error
}

func (r *encomiendaRepository) Update(encomienda *domain.Encomienda) error {
	return r.db.Save(encomienda).Error
}

func (r *encomiendaRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.Encomienda{}, "ID_ENCOMIENDA = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *encomiendaRepository) GetByID(id int64) (*domain.Encomienda, error) {
	var encomienda domain.Encomienda
	err := r.db.First(&encomienda, "ID_ENCOMIENDA = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &encomienda, nil
}

func (r *encomiendaRepository) List(offset, limit int) ([]domain.Encomienda, int, error) {
	var encomiendas []domain.Encomienda
	var total int64

	db := r.db.Model(&domain.Encomienda{})
	db.Count(&total)
	db = db.Order("ID_ENCOMIENDA ASC").Offset(offset).Limit(limit)
	err := db.Find(&encomiendas).Error
	return encomiendas, int(total), err
}
