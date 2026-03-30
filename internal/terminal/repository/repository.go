package repository

import (
	"sistema_venta_pasajes/internal/terminal/domain"

	"gorm.io/gorm"
)

type TerminalRepository interface {
	Create(terminal *domain.Terminal) error
	GetByID(id int64) (*domain.Terminal, error)
	Update(terminal *domain.Terminal) error
	Delete(id int64) error
	List() ([]domain.Terminal, error)
}

type terminalRepository struct {
	db *gorm.DB
}

func NewTerminalRepository(db *gorm.DB) TerminalRepository {
	return &terminalRepository{db: db}
}

func (r *terminalRepository) Create(terminal *domain.Terminal) error {
	return r.db.Create(terminal).Error
}

func (r *terminalRepository) GetByID(id int64) (*domain.Terminal, error) {
	var terminal domain.Terminal
	if err := r.db.First(&terminal, "ID_TERMINAL = ?", id).Error; err != nil {
		return nil, err
	}
	return &terminal, nil
}

func (r *terminalRepository) Update(terminal *domain.Terminal) error {
	return r.db.Save(terminal).Error
}

func (r *terminalRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Terminal{}, "ID_TERMINAL = ?", id).Error
}

func (r *terminalRepository) List() ([]domain.Terminal, error) {
	   var terminals []domain.Terminal
	   if err := r.db.Find(&terminals).Error; err != nil {
			   return nil, err
	   }
	   return terminals, nil
}
