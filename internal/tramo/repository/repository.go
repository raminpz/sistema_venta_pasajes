package repository

import (
	"sistema_venta_pasajes/internal/tramo/domain"

	"gorm.io/gorm"
)

type TramoRepository interface {
	Create(tramo *domain.Tramo) error
	Update(tramo *domain.Tramo) error
	Delete(id int64) error
	GetByID(id int64) (*domain.Tramo, error)
	List(offset, limit int) ([]domain.Tramo, int, error)
	ListByRuta(idRuta int64) ([]domain.Tramo, error)
	ExistsByRutaParadas(idRuta, idParadaOrigen, idParadaDestino int64) (bool, error)
}

type tramoRepository struct {
	db *gorm.DB
}

func NewTramoRepository(db *gorm.DB) TramoRepository {
	return &tramoRepository{db: db}
}

func (r *tramoRepository) Create(tramo *domain.Tramo) error {
	return r.db.Create(tramo).Error
}

func (r *tramoRepository) Update(tramo *domain.Tramo) error {
	return r.db.Save(tramo).Error
}

func (r *tramoRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.Tramo{}, "ID_TRAMO = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *tramoRepository) GetByID(id int64) (*domain.Tramo, error) {
	var t domain.Tramo
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tramoRepository) List(offset, limit int) ([]domain.Tramo, int, error) {
	var tramos []domain.Tramo
	var total int64
	db := r.db.Model(&domain.Tramo{})
	db.Count(&total)
	db.Order("ID_RUTA ASC, ID_TRAMO ASC").Offset(offset).Limit(limit).Find(&tramos)
	return tramos, int(total), db.Error
}

func (r *tramoRepository) ListByRuta(idRuta int64) ([]domain.Tramo, error) {
	var tramos []domain.Tramo
	err := r.db.Where("ID_RUTA = ?", idRuta).Order("ID_TRAMO ASC").Find(&tramos).Error
	return tramos, err
}

func (r *tramoRepository) ExistsByRutaParadas(idRuta, idParadaOrigen, idParadaDestino int64) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Tramo{}).
		Where("ID_RUTA = ? AND ID_PARADA_ORIGEN = ? AND ID_PARADA_DESTINO = ?", idRuta, idParadaOrigen, idParadaDestino).
		Count(&count).Error
	return count > 0, err
}
