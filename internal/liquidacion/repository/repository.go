package repository

import (
	"sistema_venta_pasajes/internal/liquidacion/domain"

	"gorm.io/gorm"
)

// LiquidacionRepository define las operaciones de acceso a datos.
type LiquidacionRepository interface {
	Create(liq *domain.LiquidacionViaje) error
	Update(liq *domain.LiquidacionViaje) error
	Delete(id int64) error
	GetByID(id int64) (*domain.LiquidacionViaje, error)
	GetByProgramacion(idProgramacion int64) (*domain.LiquidacionViaje, error)
	List(offset, limit int) ([]domain.LiquidacionViaje, int, error)
	ExistsByProgramacion(idProgramacion int64) (bool, error)
	GetConductorByProgramacion(idProgramacion int64) (int64, error)
	SumarVentas(idProgramacion int64) (total float64, cantidad int, err error)
	SumarEncomiendas(idProgramacion int64) (total float64, cantidad int, err error)
}
type liquidacionRepository struct {
	db *gorm.DB
}

// NewLiquidacionRepository crea una nueva instancia del repositorio.
func NewLiquidacionRepository(db *gorm.DB) LiquidacionRepository {
	return &liquidacionRepository{db: db}
}
func (r *liquidacionRepository) Create(liq *domain.LiquidacionViaje) error {
	return r.db.Create(liq).Error
}
func (r *liquidacionRepository) Update(liq *domain.LiquidacionViaje) error {
	return r.db.Save(liq).Error
}
func (r *liquidacionRepository) Delete(id int64) error {
	return r.db.Delete(&domain.LiquidacionViaje{}, id).Error
}
func (r *liquidacionRepository) GetByID(id int64) (*domain.LiquidacionViaje, error) {
	var liq domain.LiquidacionViaje
	err := r.db.First(&liq, id).Error
	if err != nil {
		return nil, err
	}
	return &liq, nil
}
func (r *liquidacionRepository) GetByProgramacion(idProgramacion int64) (*domain.LiquidacionViaje, error) {
	var liq domain.LiquidacionViaje
	err := r.db.Where("ID_PROGRAMACION = ?", idProgramacion).First(&liq).Error
	if err != nil {
		return nil, err
	}
	return &liq, nil
}
func (r *liquidacionRepository) List(offset, limit int) ([]domain.LiquidacionViaje, int, error) {
	var liqs []domain.LiquidacionViaje
	var total int64
	db := r.db.Model(&domain.LiquidacionViaje{})
	db.Count(&total)
	db.Order("ID_LIQUIDACION ASC").Offset(offset).Limit(limit).Find(&liqs)
	return liqs, int(total), db.Error
}
func (r *liquidacionRepository) ExistsByProgramacion(idProgramacion int64) (bool, error) {
	var count int64
	err := r.db.Model(&domain.LiquidacionViaje{}).
		Where("ID_PROGRAMACION = ?", idProgramacion).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *liquidacionRepository) GetConductorByProgramacion(idProgramacion int64) (int64, error) {
	var row struct {
		IDConductor int64 `gorm:"column:ID_CONDUCTOR"`
	}
	err := r.db.Table("PROGRAMACION").
		Select("ID_CONDUCTOR").
		Where("ID_PROGRAMACION = ?", idProgramacion).
		First(&row).Error
	if err != nil {
		return 0, err
	}
	return row.IDConductor, nil
}
func (r *liquidacionRepository) SumarVentas(idProgramacion int64) (float64, int, error) {
	var row struct {
		Total    float64 `gorm:"column:total"`
		Cantidad int     `gorm:"column:cantidad"`
	}
	err := r.db.Raw(
		"SELECT COALESCE(SUM(TOTAL), 0) AS total, COUNT(ID_VENTA) AS cantidad FROM VENTA WHERE ID_PROGRAMACION = ?",
		idProgramacion,
	).Scan(&row).Error
	return row.Total, row.Cantidad, err
}
func (r *liquidacionRepository) SumarEncomiendas(idProgramacion int64) (float64, int, error) {
	var row struct {
		Total    float64 `gorm:"column:total"`
		Cantidad int     `gorm:"column:cantidad"`
	}
	err := r.db.Raw(
		"SELECT COALESCE(SUM(COSTO), 0) AS total, COUNT(ID_ENCOMIENDA) AS cantidad FROM ENCOMIENDA WHERE ID_PROGRAMACION = ?",
		idProgramacion,
	).Scan(&row).Error
	return row.Total, row.Cantidad, err
}
