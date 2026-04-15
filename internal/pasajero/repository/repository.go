package repository

import (
	"sistema_venta_pasajes/internal/pasajero/domain"

	"gorm.io/gorm"
)

type PasajeroRepository interface {
	Create(pasajero *domain.Pasajero) error
	GetByID(id int64) (*domain.Pasajero, error)
	Update(pasajero *domain.Pasajero) error
	Delete(id int64) error
	List(page, size int) ([]domain.Pasajero, int, error)
	Search(query string) ([]domain.Pasajero, int, error)
}

type pasajeroRepository struct {
	db *gorm.DB
}

func NewPasajeroRepository(db *gorm.DB) PasajeroRepository {
	return &pasajeroRepository{db: db}
}

func (r *pasajeroRepository) Create(pasajero *domain.Pasajero) error {
	return r.db.Create(pasajero).Error
}

func (r *pasajeroRepository) GetByID(id int64) (*domain.Pasajero, error) {
	var pasajero domain.Pasajero
	if err := r.db.First(&pasajero, "ID_PASAJERO = ?", id).Error; err != nil {
		return nil, err
	}
	return &pasajero, nil
}

func (r *pasajeroRepository) Update(pasajero *domain.Pasajero) error {
	return r.db.Save(pasajero).Error
}

func (r *pasajeroRepository) Delete(id int64) error {
	res := r.db.Delete(&domain.Pasajero{}, "ID_PASAJERO = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *pasajeroRepository) List(page, size int) ([]domain.Pasajero, int, error) {
	var pasajeros []domain.Pasajero
	var total int64
	db := r.db.Model(&domain.Pasajero{})
	db.Count(&total)
	db.Offset((page-1)*size).Limit(size).Find(&pasajeros)
	return pasajeros, int(total), db.Error
}

// Search permite buscar pasajeros por nombre, apellido o DNI
func (r *pasajeroRepository) Search(query string) ([]domain.Pasajero, int, error) {
	var pasajeros []domain.Pasajero
	var total int64
	db := r.db.Model(&domain.Pasajero{}).Where("NOMBRES LIKE ? OR APELLIDOS LIKE ? OR NRO_DOCUMENTO LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%")
	db.Count(&total)
	db.Order("ID_PASAJERO ASC").Find(&pasajeros)
	return pasajeros, int(total), db.Error
}


