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
	IsAsientoDisponible(idProgramacion, idAsiento, idTramo int64) (bool, error)
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
	res := r.db.Delete(&domain.Venta{}, "ID_VENTA = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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

// IsAsientoDisponible verifica que el asiento no esté ocupado en un tramo solapado.
// Compara los ORDEN de las PARADA: hay solapamiento si la venta existente empieza
// antes del destino nuevo Y termina después del origen nuevo.
func (r *ventaRepository) IsAsientoDisponible(idProgramacion, idAsiento, idTramo int64) (bool, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM VENTA v
		JOIN TRAMO t_nuevo   ON t_nuevo.ID_TRAMO = ?
		JOIN PARADA po_nuevo ON po_nuevo.ID_PARADA = t_nuevo.ID_PARADA_ORIGEN
		JOIN PARADA pd_nuevo ON pd_nuevo.ID_PARADA = t_nuevo.ID_PARADA_DESTINO
		JOIN TRAMO t_exist   ON t_exist.ID_TRAMO = v.ID_TRAMO
		JOIN PARADA po_exist ON po_exist.ID_PARADA = t_exist.ID_PARADA_ORIGEN
		JOIN PARADA pd_exist ON pd_exist.ID_PARADA = t_exist.ID_PARADA_DESTINO
		WHERE v.ID_PROGRAMACION = ?
		  AND v.ID_ASIENTO = ?
		  AND po_exist.ORDEN < pd_nuevo.ORDEN
		  AND pd_exist.ORDEN > po_nuevo.ORDEN
	`, idTramo, idProgramacion, idAsiento).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
