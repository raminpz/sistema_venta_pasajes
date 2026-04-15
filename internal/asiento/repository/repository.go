package repository

import (
	"sistema_venta_pasajes/internal/asiento/domain"

	"gorm.io/gorm"
)

type AsientoRepository interface {
	Create(asiento *domain.Asiento) error
	GetByID(id int64) (*domain.Asiento, error)
	ListByVehiculo(idVehiculo int64) ([]*domain.Asiento, error)
	Update(asiento *domain.Asiento) error
	Delete(id int64) error
	CambiarEstado(id int64, estado string) error
}

type asientoRepository struct {
	db *gorm.DB
}

func NewAsientoRepository(db *gorm.DB) AsientoRepository {
	return &asientoRepository{db: db}
}

func (r *asientoRepository) Create(asiento *domain.Asiento) error {
	return r.db.Create(asiento).Error
}

func (r *asientoRepository) GetByID(id int64) (*domain.Asiento, error) {
	var a domain.Asiento
	if err := r.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *asientoRepository) ListByVehiculo(idVehiculo int64) ([]*domain.Asiento, error) {
	var asientos []*domain.Asiento
	if err := r.db.Where("id_vehiculo = ?", idVehiculo).Find(&asientos).Error; err != nil {
		return nil, err
	}
	return asientos, nil
}

func (r *asientoRepository) Update(asiento *domain.Asiento) error {
	return r.db.Save(asiento).Error
}

func (r *asientoRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.Asiento{}, "ID_ASIENTO = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *asientoRepository) CambiarEstado(id int64, estado string) error {
	return r.db.Model(&domain.Asiento{}).
		Where("ID_ASIENTO = ?", id).
		Update("ESTADO", estado).Error
}

