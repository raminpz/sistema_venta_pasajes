package repository

import (
	"sistema_venta_pasajes/internal/venta/domain"

	"gorm.io/gorm"
)

type VentaRepository interface {
	Create(venta *domain.Venta) error
	Update(venta *domain.Venta) error
	Delete(id int64) error
	GetByID(id int64) (*domain.Venta, error)
	List(offset, limit int) ([]domain.Venta, int, error)
	NextCorrelativo(serie string) (uint, error)
}

type ventaRepository struct {
	db *gorm.DB
}

func NewVentaRepository(db *gorm.DB) VentaRepository {
	return &ventaRepository{db: db}
}

func (r *ventaRepository) Create(venta *domain.Venta) error {
	return r.db.Create(venta).Error
}

func (r *ventaRepository) Update(venta *domain.Venta) error {
	return r.db.Save(venta).Error
}

func (r *ventaRepository) Delete(id int64) error {
	return r.db.Delete(&domain.Venta{}, id).Error
}

func (r *ventaRepository) GetByID(id int64) (*domain.Venta, error) {
	var venta domain.Venta
	err := r.db.First(&venta, id).Error
	if err != nil {
		return nil, err
	}
	return &venta, nil
}

func (r *ventaRepository) List(offset, limit int) ([]domain.Venta, int, error) {
	var ventas []domain.Venta
	var total int64
	db := r.db.Model(&domain.Venta{})
	db.Count(&total)
	db = db.Order("ID_VENTA ASC").Offset(offset).Limit(limit)
	err := db.Find(&ventas).Error
	return ventas, int(total), err
}

// NextCorrelativo obtiene el siguiente número correlativo para una serie dada.
// Consulta el máximo existente y retorna max+1 (empieza en 1 si no hay registros).
func (r *ventaRepository) NextCorrelativo(serie string) (uint, error) {
	var maxCorrelativo uint
	err := r.db.Model(&domain.Venta{}).
		Where("SERIE = ?", serie).
		Select("COALESCE(MAX(CORRELATIVO), 0)").
		Scan(&maxCorrelativo).Error
	if err != nil {
		return 0, err
	}
	return maxCorrelativo + 1, nil
}
