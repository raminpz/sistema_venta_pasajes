package service

import (
	"context"
	"sistema_venta_pasajes/internal/ruta/domain"
	"sistema_venta_pasajes/internal/ruta/input"
	"sistema_venta_pasajes/internal/ruta/repository"
	"sistema_venta_pasajes/pkg"
)

type Service interface {
	List(ctx context.Context) ([]domain.Ruta, error)
	GetByID(ctx context.Context, id int) (*domain.Ruta, error)
	Create(ctx context.Context, input input.CreateRutaInput) (*domain.Ruta, error)
	Update(ctx context.Context, id int, input input.UpdateRutaInput) (*domain.Ruta, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo repository.RutaRepository
}

func New(repo repository.RutaRepository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context) ([]domain.Ruta, error) {
	return s.repo.List()
}

func (s *service) GetByID(ctx context.Context, id int) (*domain.Ruta, error) {
	return s.repo.GetByID(id)
}

func (s *service) Create(ctx context.Context, in input.CreateRutaInput) (*domain.Ruta, error) {
	ruta := &domain.Ruta{
		IDOrigenTerminal:  in.IDOrigenTerminal,
		IDDestinoTerminal: in.IDDestinoTerminal,
		DuracionHoras:     in.DuracionHoras,
	}
	if err := s.repo.Create(ruta); err != nil {
		errApp := pkg.AsAppError(err)
		if errApp != nil && errApp.Code == "duplicate_resource" {
			return nil, pkg.Conflict("duplicate_resource", "Ya existe una ruta con los mismos terminales de origen y destino.")
		}
		return nil, pkg.Internal("Error al crear ruta", err)
	}
	return ruta, nil
}

func (s *service) Update(ctx context.Context, id int, in input.UpdateRutaInput) (*domain.Ruta, error) {
	ruta, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if in.IDOrigenTerminal != nil {
		ruta.IDOrigenTerminal = *in.IDOrigenTerminal
	}
	if in.IDDestinoTerminal != nil {
		ruta.IDDestinoTerminal = *in.IDDestinoTerminal
	}
	if in.DuracionHoras != nil {
		ruta.DuracionHoras = *in.DuracionHoras
	}
	if err := s.repo.Update(ruta); err != nil {
		return nil, err
	}
	return ruta, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
