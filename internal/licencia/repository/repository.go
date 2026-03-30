package repository

import (
	"context"

	"gorm.io/gorm"

	domain "sistema_venta_pasajes/internal/licencia/domain"
	input "sistema_venta_pasajes/internal/licencia/input"
)

type Repository interface {
	List(ctx context.Context) ([]domain.LicenciaSistema, error)
	Create(ctx context.Context, input input.CreateInput) (*domain.LicenciaSistema, error)
	Delete(ctx context.Context, id int64) error
	Activar(ctx context.Context, id int64) error
	Bloquear(ctx context.Context, id int64) error
	Renovar(ctx context.Context, id int64, nuevaFechaExpiracion string) error
	GenerarClave(ctx context.Context) (string, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context) ([]domain.LicenciaSistema, error) {
	// Implementación real aquí
	return nil, nil
}

func (r *repository) Create(ctx context.Context, input input.CreateInput) (*domain.LicenciaSistema, error) {
	// Implementación real aquí
	return nil, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	// Implementación real aquí
	return nil
}

func (r *repository) Activar(ctx context.Context, id int64) error {
	// Cambiar estado a OPERATIVO
	return nil
}

func (r *repository) Bloquear(ctx context.Context, id int64) error {
	// Cambiar estado a BLOQUEADO
	return nil
}

func (r *repository) Renovar(ctx context.Context, id int64, nuevaFechaExpiracion string) error {
	// Actualizar fechas y estado
	return nil
}

func (r *repository) GenerarClave(ctx context.Context) (string, error) {
	// Lógica para generar clave única
	return "", nil
}
