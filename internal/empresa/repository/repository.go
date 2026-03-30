package repository

import (
	"sistema_venta_pasajes/internal/empresa/domain"
	"gorm.io/gorm"
)

type EmpresaRepository interface {
	Create(empresa *domain.Empresa) error
	GetByID(id int64) (*domain.Empresa, error)
	Update(empresa *domain.Empresa) error
	Delete(id int64) error
	List() ([]domain.Empresa, error)
}

type empresaRepository struct {
	db *gorm.DB
}

func NewEmpresaRepository(db *gorm.DB) EmpresaRepository {
	return &empresaRepository{db: db}
}

func (r *empresaRepository) Create(empresa *domain.Empresa) error {
	return r.db.Create(empresa).Error
}

func (r *empresaRepository) GetByID(id int64) (*domain.Empresa, error) {
	   var empresa domain.Empresa
	   if err := r.db.First(&empresa, "ID_EMPRESA = ?", id).Error; err != nil {
			   return nil, err
	   }
	   return &empresa, nil
}

func (r *empresaRepository) Update(empresa *domain.Empresa) error {
	return r.db.Save(empresa).Error
}

func (r *empresaRepository) Delete(id int64) error {
	   res := r.db.Delete(&domain.Empresa{}, "ID_EMPRESA = ?", id)
	   if res.Error != nil {
			   return res.Error
	   }
	   if res.RowsAffected == 0 {
			   return gorm.ErrRecordNotFound
	   }
	   return nil
}

func (r *empresaRepository) List() ([]domain.Empresa, error) {
	   var empresas []domain.Empresa
	   err := r.db.Model(&domain.Empresa{}).Order("ID_EMPRESA ASC").Find(&empresas).Error
	   return empresas, err
}
