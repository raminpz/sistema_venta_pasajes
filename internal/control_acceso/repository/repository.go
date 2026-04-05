package repository

import (
	"time"

	"gorm.io/gorm"

	"sistema_venta_pasajes/internal/control_acceso/domain"
)

type ControlAccesoRepository interface {
	GetLatest() (*domain.ControlAcceso, error)
	GetByID(id int64) (*domain.ControlAcceso, error)
	Create(acceso *domain.ControlAcceso) error
	SetEstado(id int64, estado string) error
	Renovar(id int64, nuevaFecha time.Time) error
}

type controlAccesoRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) ControlAccesoRepository {
	return &controlAccesoRepository{db: db}
}

func (r *controlAccesoRepository) GetLatest() (*domain.ControlAcceso, error) {
	var a domain.ControlAcceso
	if err := r.db.Order("ID_ACCESO DESC").First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *controlAccesoRepository) GetByID(id int64) (*domain.ControlAcceso, error) {
	var a domain.ControlAcceso
	if err := r.db.First(&a, "ID_ACCESO = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *controlAccesoRepository) Create(acceso *domain.ControlAcceso) error {
	return r.db.Create(acceso).Error
}

func (r *controlAccesoRepository) SetEstado(id int64, estado string) error {
	result := r.db.Model(&domain.ControlAcceso{}).
		Where("ID_ACCESO = ?", id).
		Update("ESTADO", estado)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *controlAccesoRepository) Renovar(id int64, nuevaFecha time.Time) error {
	result := r.db.Model(&domain.ControlAcceso{}).
		Where("ID_ACCESO = ?", id).
		Updates(map[string]interface{}{
			"FECHA_EXPIRACION": nuevaFecha,
			"ESTADO":           "OPERATIVO",
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
