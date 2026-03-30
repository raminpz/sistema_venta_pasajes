package service

import (
	"context"
	domain "sistema_venta_pasajes/internal/licencia/domain"
	input "sistema_venta_pasajes/internal/licencia/input"
)

type Service interface {
	List(ctx context.Context) ([]domain.LicenciaSistema, error)
	GetByID(ctx context.Context, id int64) (*domain.LicenciaSistema, error)
	Create(ctx context.Context, input input.CreateInput) (*domain.LicenciaSistema, error)
	Update(ctx context.Context, id int64, input input.UpdateInput) (*domain.LicenciaSistema, error)
	Delete(ctx context.Context, id int64) error
}
