package repository

import (
	"sistema_venta_pasajes/internal/pago/domain"

	"gorm.io/gorm"
)

type PagoRepository interface {
	Create(pago *domain.Pago) error
	Update(pago *domain.Pago) error
	Delete(id int64) error
	GetByID(id int64) (*domain.Pago, error)
	List(offset, limit int, idVenta *int64) ([]domain.Pago, int, error)
}

type pagoRepository struct {
	db *gorm.DB
}

func NewPagoRepository(db *gorm.DB) PagoRepository {
	return &pagoRepository{db: db}
}

func (r *pagoRepository) Create(pago *domain.Pago) error {
	return r.db.Create(pago).Error
}

func (r *pagoRepository) Update(pago *domain.Pago) error {
	return r.db.Save(pago).Error
}

func (r *pagoRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.Pago{}, "ID_PAGO = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *pagoRepository) GetByID(id int64) (*domain.Pago, error) {
	var pago domain.Pago
	err := r.db.First(&pago, "ID_PAGO = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pago, nil
}

func (r *pagoRepository) List(offset, limit int, idVenta *int64) ([]domain.Pago, int, error) {
	var pagos []domain.Pago
	var total int64
	db := r.db.Model(&domain.Pago{})
	if idVenta != nil {
		db = db.Where("ID_VENTA = ?", *idVenta)
	}
	db.Count(&total)
	db = db.Order("ID_PAGO ASC").Offset(offset).Limit(limit)
	err := db.Find(&pagos).Error
	return pagos, int(total), err
}
