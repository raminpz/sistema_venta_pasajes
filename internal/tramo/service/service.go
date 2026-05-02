package service

import (
	"errors"
	"sistema_venta_pasajes/internal/tramo/domain"
	"sistema_venta_pasajes/internal/tramo/input"
	"sistema_venta_pasajes/internal/tramo/repository"
	"sistema_venta_pasajes/internal/tramo/util"
	"sistema_venta_pasajes/pkg"

	"gorm.io/gorm"
)

type TramoService interface {
	Create(input.CreateTramoInput) (*input.TramoOutput, error)
	Update(input.UpdateTramoInput) (*input.TramoOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*input.TramoOutput, error)
	List(page, size int) ([]input.TramoOutput, int, error)
	ListByRuta(idRuta int64) ([]input.TramoOutput, error)
}

type tramoService struct {
	repo repository.TramoRepository
}

func NewTramoService(repo repository.TramoRepository) TramoService {
	return &tramoService{repo: repo}
}

func (s *tramoService) Create(in input.CreateTramoInput) (*input.TramoOutput, error) {
	pkg.TrimSpacesOnStruct(&in)
	if err := util.ValidarCamposCreate(in); err != nil {
		return nil, err
	}
	exists, err := s.repo.ExistsByRutaParadas(in.IDRuta, in.IDParadaOrigen, in.IDParadaDestino)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	if exists {
		return nil, pkg.BadRequest("duplicate_tramo", util.ERR_DUPLICATE)
	}
	tramo := &domain.Tramo{
		IDRuta:          in.IDRuta,
		IDParadaOrigen:  in.IDParadaOrigen,
		IDParadaDestino: in.IDParadaDestino,
	}
	if err := s.repo.Create(tramo); err != nil {
		return nil, pkg.ParseDBError(err, "db_error", "Error al guardar el tramo", nil,
			map[string]string{"UQ_TRAMO_RUTA_PARADAS": util.ERR_DUPLICATE})
	}
	return mapOutput(tramo), nil
}

func (s *tramoService) Update(in input.UpdateTramoInput) (*input.TramoOutput, error) {
	if err := util.ValidarCamposUpdate(in); err != nil {
		return nil, err
	}
	tramo, err := s.repo.GetByID(in.IDTramo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return nil, pkg.Internal(err.Error())
	}
	if in.IDRuta != nil {
		tramo.IDRuta = *in.IDRuta
	}
	if in.IDParadaOrigen != nil {
		tramo.IDParadaOrigen = *in.IDParadaOrigen
	}
	if in.IDParadaDestino != nil {
		tramo.IDParadaDestino = *in.IDParadaDestino
	}
	if tramo.IDParadaOrigen == tramo.IDParadaDestino {
		return nil, pkg.BadRequest("paradas_iguales", util.ERR_PARADAS_IGUALES)
	}
	if err := s.repo.Update(tramo); err != nil {
		return nil, pkg.ParseDBError(err, "db_error", "Error al actualizar el tramo", nil,
			map[string]string{"UQ_TRAMO_RUTA_PARADAS": util.ERR_DUPLICATE})
	}
	return mapOutput(tramo), nil
}

func (s *tramoService) Delete(id int64) error {
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
			map[string]string{"FK_VENTA_TRAMO": "No se puede eliminar: hay ventas asociadas a este tramo"}, nil)
	}
	return nil
}

func (s *tramoService) GetByID(id int64) (*input.TramoOutput, error) {
	tramo, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return nil, pkg.Internal(err.Error())
	}
	return mapOutput(tramo), nil
}

func (s *tramoService) List(page, size int) ([]input.TramoOutput, int, error) {
	offset, limit, _ := pkg.Paginate(page, size, 0)
	tramos, total, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, 0, pkg.Internal(err.Error())
	}
	out := make([]input.TramoOutput, 0, len(tramos))
	for _, t := range tramos {
		out = append(out, *mapOutput(&t))
	}
	return out, total, nil
}

func (s *tramoService) ListByRuta(idRuta int64) ([]input.TramoOutput, error) {
	tramos, err := s.repo.ListByRuta(idRuta)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	out := make([]input.TramoOutput, 0, len(tramos))
	for _, t := range tramos {
		out = append(out, *mapOutput(&t))
	}
	return out, nil
}

func mapOutput(t *domain.Tramo) *input.TramoOutput {
	return &input.TramoOutput{
		IDTramo:         t.IDTramo,
		IDRuta:          t.IDRuta,
		IDParadaOrigen:  t.IDParadaOrigen,
		IDParadaDestino: t.IDParadaDestino,
	}
}
