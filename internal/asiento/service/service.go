package service

import (
	"sistema_venta_pasajes/internal/asiento/domain"
	"sistema_venta_pasajes/internal/asiento/input"
	"sistema_venta_pasajes/internal/asiento/repository"
	"sistema_venta_pasajes/internal/asiento/util"
)

type AsientoService interface {
	Create(in input.CreateAsientoInput) (*domain.Asiento, error)
	GetByID(id int64) (*domain.Asiento, error)
	ListByVehiculo(idVehiculo int64) ([]*domain.Asiento, error)
	Update(id int64, in input.UpdateAsientoInput) error
	Delete(id int64) error
	CambiarEstado(id int64, estado string) error
}

type asientoService struct {
	repo repository.AsientoRepository
}

func NewAsientoService(repo repository.AsientoRepository) AsientoService {
	return &asientoService{repo: repo}
}

func (s *asientoService) Create(in input.CreateAsientoInput) (*domain.Asiento, error) {
	if err := util.ValidateEstadoAsiento(in.Estado); err != nil {
		return nil, err
	}
	asiento := &domain.Asiento{
		IDVehiculo:    in.IDVehiculo,
		NumeroAsiento: in.NumeroAsiento,
		Estado:        in.Estado,
	}
	err := s.repo.Create(asiento)
	if err != nil {
		return nil, err
	}
	return asiento, nil
}

func (s *asientoService) GetByID(id int64) (*domain.Asiento, error) {
	return s.repo.GetByID(id)
}

func (s *asientoService) ListByVehiculo(idVehiculo int64) ([]*domain.Asiento, error) {
	return s.repo.ListByVehiculo(idVehiculo)
}

func (s *asientoService) Update(id int64, in input.UpdateAsientoInput) error {
	if err := util.ValidateEstadoAsiento(in.Estado); err != nil {
		return err
	}
	asiento, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	asiento.NumeroAsiento = in.NumeroAsiento
	asiento.Estado = in.Estado
	return s.repo.Update(asiento)
}

func (s *asientoService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *asientoService) CambiarEstado(id int64, estado string) error {
	return s.repo.CambiarEstado(id, estado)
}
