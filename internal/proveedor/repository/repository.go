package repository

import (
	"context"
	"fmt"
	"sistema_venta_pasajes/internal/proveedor/util"
	"sistema_venta_pasajes/pkg"

	"sistema_venta_pasajes/internal/proveedor/domain"
	providerinput "sistema_venta_pasajes/internal/proveedor/input"

	"gorm.io/gorm"
)

const queryByIDProveedor = "ID_PROVEEDOR = ?"

type Repository interface {
	List(ctx context.Context) ([]domain.ProveedorSistema, error)
	GetByID(ctx context.Context, id int64) (*domain.ProveedorSistema, error)
	Create(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error)
	Update(ctx context.Context, id int64, input providerinput.UpdateInput) (*domain.ProveedorSistema, error)
	Delete(ctx context.Context, id int64) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) List(ctx context.Context) ([]domain.ProveedorSistema, error) {
	proveedores := make([]domain.ProveedorSistema, 0)
	if err := r.db.WithContext(ctx).
		Model(&domain.ProveedorSistema{}).
		Order("ID_PROVEEDOR ASC").
		Find(&proveedores).Error; err != nil {
		return nil, fmt.Errorf("iterar proveedores del sistema: %w", err)
	}

	return proveedores, nil
}

func (r *repo) GetByID(ctx context.Context, id int64) (*domain.ProveedorSistema, error) {
	var proveedor domain.ProveedorSistema
	if err := r.db.WithContext(ctx).
		Model(&domain.ProveedorSistema{}).
		First(&proveedor, queryByIDProveedor, id).Error; err != nil {
		return nil, fmt.Errorf("obtener proveedor del sistema por id: %w", err)
	}

	return &proveedor, nil
}

func (r *repo) Create(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error) {
	record := proveedorSistemaRecord{
		RUC:             input.RUC,
		RazonSocial:     input.RazonSocial,
		NombreComercial: stringPtrIfNotEmpty(input.NombreComercial),
		Direccion:       stringPtrIfNotEmpty(input.Direccion),
		Telefono:        stringPtrIfNotEmpty(input.Telefono),
		Email:           stringPtrIfNotEmpty(input.Email),
		Web:             stringPtrIfNotEmpty(input.Web),
	}

	   if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
			   return nil, err // No envolver, para que el service lo analice
	   }

	return r.GetByID(ctx, record.IDProveedor)
}

func (r *repo) Update(ctx context.Context, id int64, input providerinput.UpdateInput) (*domain.ProveedorSistema, error) {
	result := r.db.WithContext(ctx).
		Model(&domain.ProveedorSistema{}).
		Where(queryByIDProveedor, id).
		Updates(map[string]any{
			"RUC":              input.RUC,
			"RAZON_SOCIAL":     input.RazonSocial,
			"NOMBRE_COMERCIAL": nilIfEmpty(input.NombreComercial),
			"DIRECCION":        nilIfEmpty(input.Direccion),
			"TELEFONO":         nilIfEmpty(input.Telefono),
			"EMAIL":            nilIfEmpty(input.Email),
			"WEB":              nilIfEmpty(input.Web),
		})
	if result.Error != nil {
		return nil, fmt.Errorf(util.MSG_UPDATE_PROVIDER_ERROR, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, pkg.NotFound("provider_not_found", util.MSG_PROVIDER_NOT_FOUND)
	}

	return r.GetByID(ctx, id)
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).
		Delete(&domain.ProveedorSistema{}, queryByIDProveedor, id)
	if result.Error != nil {
		return fmt.Errorf(util.MSG_DELETE_PROVIDER_ERROR, result.Error)
	}

	if result.RowsAffected == 0 {
		return pkg.NotFound("provider_not_found", util.MSG_PROVIDER_NOT_FOUND)
	}

	return nil
}

type proveedorSistemaRecord struct {
	IDProveedor     int64   `gorm:"column:ID_PROVEEDOR;primaryKey;autoIncrement"`
	RUC             string  `gorm:"column:RUC"`
	RazonSocial     string  `gorm:"column:RAZON_SOCIAL"`
	NombreComercial *string `gorm:"column:NOMBRE_COMERCIAL"`
	Direccion       *string `gorm:"column:DIRECCION"`
	Telefono        *string `gorm:"column:TELEFONO"`
	Email           *string `gorm:"column:EMAIL"`
	Web             *string `gorm:"column:WEB"`
}

func (proveedorSistemaRecord) TableName() string {
	return "PROVEEDOR_SISTEMA"
}

func nilIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func stringPtrIfNotEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
