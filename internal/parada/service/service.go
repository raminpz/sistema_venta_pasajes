package service

import (
	"errors"
	"sistema_venta_pasajes/internal/parada/domain"
	"sistema_venta_pasajes/internal/parada/input"
	"sistema_venta_pasajes/internal/parada/repository"
	"sistema_venta_pasajes/internal/parada/util"
	"sistema_venta_pasajes/pkg"

	"gorm.io/gorm"
)

type ParadaService interface {
	Create(input.CreateParadaInput) (*input.ParadaOutput, error)
	Update(input.UpdateParadaInput) (*input.ParadaOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*input.ParadaOutput, error)
	ListByRuta(idRuta int64) ([]input.ParadaOutput, error)
}

type paradaService struct {
	repo repository.ParadaRepository
}

func NewParadaService(repo repository.ParadaRepository) ParadaService {
	return &paradaService{repo: repo}
}

func (s *paradaService) Create(in input.CreateParadaInput) (*input.ParadaOutput, error) {
	if err := util.ValidarCamposCreate(in); err != nil {
		return nil, err
	}

	existsTerminal, err := s.repo.ExistsByRutaTerminal(in.IDRuta, in.IDTerminal)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	if existsTerminal {
		return nil, pkg.BadRequest("duplicate_terminal", util.ERR_DUPLICATE)
	}

	existsOrden, err := s.repo.ExistsByRutaOrden(in.IDRuta, in.Orden)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	if existsOrden {
		return nil, pkg.BadRequest("duplicate_orden", util.ERR_DUPLICATE_ORDEN)
	}

	parada := &domain.Parada{
		IDRuta:     in.IDRuta,
		IDTerminal: in.IDTerminal,
		Orden:      in.Orden,
	}
	if err := s.repo.Create(parada); err != nil {
		return nil, pkg.ParseDBError(err, "db_error", "Error al crear la parada",
			map[string]string{"FK_PARADA_RUTA": "La ruta indicada no existe", "FK_PARADA_TERMINAL": "El terminal indicado no existe"},
			map[string]string{"UQ_PARADA_RUTA_TERMINAL": util.ERR_DUPLICATE, "UQ_PARADA_RUTA_ORDEN": util.ERR_DUPLICATE_ORDEN},
		)
	}
	return mapOutput(parada), nil
}

func (s *paradaService) Update(in input.UpdateParadaInput) (*input.ParadaOutput, error) {
	if err := util.ValidarCamposUpdate(in); err != nil {
		return nil, err
	}

	parada, err := s.repo.GetByID(in.IDParada)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return nil, pkg.Internal(err.Error())
	}

	if in.IDTerminal != nil {
		existsTerminal, err := s.repo.ExistsByRutaTerminal(parada.IDRuta, *in.IDTerminal)
		if err != nil {
			return nil, pkg.Internal(err.Error())
		}
		if existsTerminal {
			return nil, pkg.BadRequest("duplicate_terminal", util.ERR_DUPLICATE)
		}
		parada.IDTerminal = *in.IDTerminal
	}
	if in.Orden != nil {
		existsOrden, err := s.repo.ExistsByRutaOrden(parada.IDRuta, *in.Orden)
		if err != nil {
			return nil, pkg.Internal(err.Error())
		}
		if existsOrden {
			return nil, pkg.BadRequest("duplicate_orden", util.ERR_DUPLICATE_ORDEN)
		}
		parada.Orden = *in.Orden
	}

	if err := s.repo.Update(parada); err != nil {
		return nil, pkg.ParseDBError(err, "db_error", "Error al actualizar la parada", nil, nil)
	}
	return mapOutput(parada), nil
}

func (s *paradaService) Delete(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return pkg.Internal(util.ERR_DELETE, err)
	}
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return pkg.ParseDBError(err, util.ERR_CODE_DELETE, util.ERR_DELETE,
			map[string]string{"FK_TRAMO_PARADA_ORIGEN": "No se puede eliminar: hay tramos que usan esta parada como origen", "FK_TRAMO_PARADA_DESTINO": "No se puede eliminar: hay tramos que usan esta parada como destino"},
			nil,
		)
	}
	return nil
}

func (s *paradaService) GetByID(id int64) (*input.ParadaOutput, error) {
	parada, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return nil, pkg.Internal(err.Error())
	}
	return mapOutput(parada), nil
}

func (s *paradaService) ListByRuta(idRuta int64) ([]input.ParadaOutput, error) {
	paradas, err := s.repo.ListByRuta(idRuta)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	out := make([]input.ParadaOutput, 0, len(paradas))
	for _, p := range paradas {
		out = append(out, *mapOutput(&p))
	}
	return out, nil
}

func mapOutput(p *domain.Parada) *input.ParadaOutput {
	return &input.ParadaOutput{
		IDParada:   p.IDParada,
		IDRuta:     p.IDRuta,
		IDTerminal: p.IDTerminal,
		Orden:      p.Orden,
	}
}

