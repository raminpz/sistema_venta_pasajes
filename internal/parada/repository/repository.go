package repository

import (
	"sistema_venta_pasajes/internal/parada/domain"

	"gorm.io/gorm"
)

type ParadaRepository interface {
	Create(p *domain.Parada) error
	Update(p *domain.Parada) error
	Delete(id int64) error
	GetByID(id int64) (*domain.Parada, error)
	ListByRuta(idRuta int64) ([]domain.Parada, error)
	ExistsByRutaTerminal(idRuta, idTerminal int64) (bool, error)
	ExistsByRutaOrden(idRuta int64, orden int) (bool, error)
	GetOrdenByID(idParada int64) (int, error)
}

type paradaRepository struct {
	db *gorm.DB
}

func NewParadaRepository(db *gorm.DB) ParadaRepository {
	return &paradaRepository{db: db}
}

func (r *paradaRepository) Create(p *domain.Parada) error {
	return r.db.Create(p).Error
}

func (r *paradaRepository) Update(p *domain.Parada) error {
	return r.db.Save(p).Error
}

func (r *paradaRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.Parada{}, "ID_PARADA = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *paradaRepository) GetByID(id int64) (*domain.Parada, error) {
	var p domain.Parada
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paradaRepository) ListByRuta(idRuta int64) ([]domain.Parada, error) {
	var paradas []domain.Parada
	err := r.db.Where("ID_RUTA = ?", idRuta).Order("ORDEN ASC").Find(&paradas).Error
	return paradas, err
}

func (r *paradaRepository) ExistsByRutaTerminal(idRuta, idTerminal int64) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Parada{}).
		Where("ID_RUTA = ? AND ID_TERMINAL = ?", idRuta, idTerminal).
		Count(&count).Error
	return count > 0, err
}

func (r *paradaRepository) ExistsByRutaOrden(idRuta int64, orden int) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Parada{}).
		Where("ID_RUTA = ? AND ORDEN = ?", idRuta, orden).
		Count(&count).Error
	return count > 0, err
}

func (r *paradaRepository) GetOrdenByID(idParada int64) (int, error) {
	var p domain.Parada
	if err := r.db.Select("ORDEN").First(&p, idParada).Error; err != nil {
		return 0, err
	}
	return p.Orden, nil
}

