package repository

import (
	"sistema_venta_pasajes/internal/programacion/domain"

	"gorm.io/gorm"
)

type ProgramacionRepository interface {
	Create(programacion *domain.Programacion) error
	Update(programacion *domain.Programacion) error
	Delete(id int64) error
	GetByID(id int64) (*domain.Programacion, error)
	List(offset, limit int) ([]domain.Programacion, int, error)
}

type programacionRepository struct {
	db *gorm.DB
}

func NewProgramacionRepository(db *gorm.DB) ProgramacionRepository {
	return &programacionRepository{db: db}
}

func (r *programacionRepository) Create(programacion *domain.Programacion) error {
	return r.db.Create(programacion).Error
}

func (r *programacionRepository) Update(programacion *domain.Programacion) error {
	return r.db.Save(programacion).Error
}

func (r *programacionRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.Programacion{}, "ID_PROGRAMACION = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *programacionRepository) GetByID(id int64) (*domain.Programacion, error) {
	var programacion domain.Programacion
	err := r.db.First(&programacion, "ID_PROGRAMACION = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &programacion, nil
}

func (r *programacionRepository) List(offset, limit int) ([]domain.Programacion, int, error) {
	var programaciones []domain.Programacion
	var total int64
	db := r.db.Model(&domain.Programacion{})
	db.Count(&total)
	db = db.Order("ID_PROGRAMACION ASC").Offset(offset).Limit(limit)
	err := db.Find(&programaciones).Error
	return programaciones, int(total), err
}


