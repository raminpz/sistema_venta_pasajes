package repository

import (
	"sistema_venta_pasajes/internal/asiento_tramo/domain"

	"gorm.io/gorm"
)

type AsientoTramoRepository interface {
	Create(at *domain.AsientoTramo) error
	Update(at *domain.AsientoTramo) error
	Delete(id int64) error
	GetByID(id int64) (*domain.AsientoTramo, error)
	GetByAsientoTramo(idAsiento, idTramo int64) (*domain.AsientoTramo, error)
	GetDisponiblesEnTramo(idTramo int64) ([]domain.AsientoTramo, error)
	MarkAsOccupied(idAsiento, idTramo int64, idVenta *int64) error
	MarkAsAvailable(idAsiento, idTramo int64) error
	DeleteByVenta(idVenta int64) error
}

type asientoTramoRepository struct {
	db *gorm.DB
}

func NewAsientoTramoRepository(db *gorm.DB) AsientoTramoRepository {
	return &asientoTramoRepository{db: db}
}

func (r *asientoTramoRepository) Create(at *domain.AsientoTramo) error {
	return r.db.Create(at).Error
}

func (r *asientoTramoRepository) Update(at *domain.AsientoTramo) error {
	return r.db.Save(at).Error
}

func (r *asientoTramoRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.AsientoTramo{}, "ID_ASIENTO_TRAMO = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *asientoTramoRepository) GetByID(id int64) (*domain.AsientoTramo, error) {
	var at domain.AsientoTramo
	err := r.db.First(&at, id).Error
	if err != nil {
		return nil, err
	}
	return &at, nil
}

func (r *asientoTramoRepository) GetByAsientoTramo(idAsiento, idTramo int64) (*domain.AsientoTramo, error) {
	var at domain.AsientoTramo
	err := r.db.Where("ID_ASIENTO = ? AND ID_TRAMO = ?", idAsiento, idTramo).First(&at).Error
	if err != nil {
		return nil, err
	}
	return &at, nil
}

func (r *asientoTramoRepository) GetDisponiblesEnTramo(idTramo int64) ([]domain.AsientoTramo, error) {
	var asientos []domain.AsientoTramo
	err := r.db.Where("ID_TRAMO = ? AND ESTADO = ?", idTramo, "DISPONIBLE").Find(&asientos).Error
	return asientos, err
}

func (r *asientoTramoRepository) MarkAsOccupied(idAsiento, idTramo int64, idVenta *int64) error {
	return r.db.Model(&domain.AsientoTramo{}).
		Where("ID_ASIENTO = ? AND ID_TRAMO = ?", idAsiento, idTramo).
		Updates(map[string]interface{}{
			"ESTADO":   "OCUPADO",
			"ID_VENTA": idVenta,
		}).Error
}

func (r *asientoTramoRepository) MarkAsAvailable(idAsiento, idTramo int64) error {
	return r.db.Model(&domain.AsientoTramo{}).
		Where("ID_ASIENTO = ? AND ID_TRAMO = ?", idAsiento, idTramo).
		Updates(map[string]interface{}{
			"ESTADO":   "DISPONIBLE",
			"ID_VENTA": nil,
		}).Error
}

func (r *asientoTramoRepository) DeleteByVenta(idVenta int64) error {
	res := r.db.Delete(&domain.AsientoTramo{}, "ID_VENTA = ?", idVenta)
	return res.Error
}
