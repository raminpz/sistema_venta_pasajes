package service

import (
	"errors"
	"net/http"
	"sistema_venta_pasajes/internal/asiento_tramo/domain"
	"sistema_venta_pasajes/internal/asiento_tramo/input"
	"sistema_venta_pasajes/internal/asiento_tramo/repository"
	"sistema_venta_pasajes/internal/asiento_tramo/util"
	"sistema_venta_pasajes/pkg"

	"gorm.io/gorm"
)

type AsientoTramoService interface {
	Create(in input.CreateAsientoTramoInput) (*input.AsientoTramoOutput, error)
	GetByID(id int64) (*input.AsientoTramoOutput, error)
	GetDisponiblesEnTramo(idTramo int64) ([]input.AsientoTramoOutput, error)
	MarkAsOccupied(idAsiento, idTramo int64, idVenta *int64) error
	MarkAsAvailable(idAsiento, idTramo int64) error
	DeleteByVenta(idVenta int64) error
	IsAsientoDisponible(idAsiento, idTramo int64) (bool, error)
}

type asientoTramoService struct {
	repo repository.AsientoTramoRepository
}

func NewAsientoTramoService(repo repository.AsientoTramoRepository) AsientoTramoService {
	return &asientoTramoService{repo: repo}
}

func toAsientoTramoOutput(at *domain.AsientoTramo) *input.AsientoTramoOutput {
	return &input.AsientoTramoOutput{
		IDAsientoTramo: at.IDAsientoTramo,
		IDVenta:        at.IDVenta,
		IDAsiento:      at.IDAsiento,
		IDTramo:        at.IDTramo,
		Estado:         at.Estado,
	}
}

func (s *asientoTramoService) Create(in input.CreateAsientoTramoInput) (*input.AsientoTramoOutput, error) {
	if in.IDAsiento <= 0 || in.IDTramo <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_ASIENTO_TRAMO_DUPLICATE, util.MSG_ASIENTO_TRAMO_DUPLICATE)
	}

	at := &domain.AsientoTramo{
		IDVenta:   in.IDVenta,
		IDAsiento: in.IDAsiento,
		IDTramo:   in.IDTramo,
		Estado:    in.Estado,
	}

	if err := s.repo.Create(at); err != nil {
		return nil, util.ParseDBError(err, util.ERR_CODE_ASIENTO_TRAMO_CREATE, util.MSG_ASIENTO_TRAMO_CREATE_ERROR)
	}
	return toAsientoTramoOutput(at), nil
}

func (s *asientoTramoService) GetByID(id int64) (*input.AsientoTramoOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_ASIENTO_TRAMO_NOT_FOUND, util.MSG_ASIENTO_TRAMO_NOT_FOUND)
	}
	at, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_ASIENTO_TRAMO_NOT_FOUND, util.MSG_ASIENTO_TRAMO_NOT_FOUND)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_ASIENTO_TRAMO_NOT_FOUND, util.MSG_ASIENTO_TRAMO_NOT_FOUND).WithCause(err)
	}
	return toAsientoTramoOutput(at), nil
}

func (s *asientoTramoService) GetDisponiblesEnTramo(idTramo int64) ([]input.AsientoTramoOutput, error) {
	if idTramo <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_ASIENTO_TRAMO_NOT_FOUND, util.MSG_ASIENTO_TRAMO_NOT_FOUND)
	}
	asientos, err := s.repo.GetDisponiblesEnTramo(idTramo)
	if err != nil {
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_ASIENTO_TRAMO_NOT_FOUND, util.MSG_ASIENTO_TRAMO_NOT_FOUND).WithCause(err)
	}
	outs := make([]input.AsientoTramoOutput, 0, len(asientos))
	for _, at := range asientos {
		outs = append(outs, *toAsientoTramoOutput(&at))
	}
	return outs, nil
}

func (s *asientoTramoService) MarkAsOccupied(idAsiento, idTramo int64, idVenta *int64) error {
	if idAsiento <= 0 || idTramo <= 0 {
		return pkg.BadRequest(util.ERR_CODE_ASIENTO_TRAMO_DUPLICATE, util.MSG_ASIENTO_TRAMO_DUPLICATE)
	}
	return s.repo.MarkAsOccupied(idAsiento, idTramo, idVenta)
}

func (s *asientoTramoService) MarkAsAvailable(idAsiento, idTramo int64) error {
	if idAsiento <= 0 || idTramo <= 0 {
		return pkg.BadRequest(util.ERR_CODE_ASIENTO_TRAMO_DUPLICATE, util.MSG_ASIENTO_TRAMO_DUPLICATE)
	}
	return s.repo.MarkAsAvailable(idAsiento, idTramo)
}

func (s *asientoTramoService) DeleteByVenta(idVenta int64) error {
	if idVenta <= 0 {
		return pkg.BadRequest(util.ERR_CODE_ASIENTO_TRAMO_DELETE, util.MSG_ASIENTO_TRAMO_DELETE_ERROR)
	}
	return s.repo.DeleteByVenta(idVenta)
}

func (s *asientoTramoService) IsAsientoDisponible(idAsiento, idTramo int64) (bool, error) {
	if idAsiento <= 0 || idTramo <= 0 {
		return false, pkg.BadRequest(util.ERR_CODE_ASIENTO_TRAMO_DUPLICATE, util.MSG_ASIENTO_TRAMO_DUPLICATE)
	}
	at, err := s.repo.GetByAsientoTramo(idAsiento, idTramo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // Si no existe el registro, está disponible (pero necesita inicializarse)
		}
		return false, err
	}
	return at.Estado == util.ESTADO_DISPONIBLE, nil
}
