package service

import (
	"context"
	"errors"
	"sistema_venta_pasajes/internal/conductor/domain"
	"sistema_venta_pasajes/internal/conductor/input"
	"sistema_venta_pasajes/internal/conductor/repository"
	"sistema_venta_pasajes/internal/conductor/util"
	"sistema_venta_pasajes/pkg"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	List(ctx context.Context) ([]domain.Conductor, error)
	GetByID(ctx context.Context, id int64) (*domain.Conductor, error)
	Create(ctx context.Context, input input.CreateConductorInput) (*domain.Conductor, error)
	Update(ctx context.Context, id int64, input input.UpdateConductorInput) (*domain.Conductor, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo repository.ConductorRepository
}

func New(repo repository.ConductorRepository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context) ([]domain.Conductor, error) {
	return s.repo.List()
}

func (s *service) GetByID(ctx context.Context, id int64) (*domain.Conductor, error) {
	if id <= 0 {
		return nil, pkg.BadRequest("invalid_conductor_id", util.MSG_INVALID_ID)
	}
	conductor, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.NotFound("conductor_not_found", util.MSG_NOT_FOUND)
		}
		return nil, pkg.Internal(util.MSG_LIST_ERROR, err)
	}
	return conductor, nil
}

func (s *service) Create(ctx context.Context, input input.CreateConductorInput) (*domain.Conductor, error) {
	pkg.TrimSpacesOnStruct(&input)
	if err := input.Validate(); err != nil {
		return nil, err
	}
	fechaVenc, err := time.Parse("2006-01-02", input.FechaVencLicencia)
	if err != nil {
		return nil, errors.New("El formato de fecha_venc_licencia debe ser YYYY-MM-DD")
	}
	conductor := &domain.Conductor{
		Nombres:           pkg.CapitalizeWords(input.Nombres),
		Apellidos:         pkg.CapitalizeWords(input.Apellidos),
		DNI:               input.DNI,
		NumeroLicencia:    input.NumeroLicencia,
		Telefono:          input.Telefono,
		Direccion:         input.Direccion,
		FechaVencLicencia: fechaVenc,
	}
	if err := s.repo.Create(conductor); err != nil {
		errApp := pkg.AsAppError(err)
		if errApp != nil && errApp.Code == "duplicate_resource" {
			// Analizar el error original para saber si es por DNI o control_acceso
			if errApp.Err != nil {
				errStr := errApp.Err.Error()
				if strings.Contains(errStr, "dni") || strings.Contains(errStr, "DNI") {
					return nil, pkg.Conflict("duplicate_resource", util.ERR_DNI_DUPLICADO)
				}
				if strings.Contains(errStr, "control_acceso") || strings.Contains(errStr, "LICENCIA") {
					return nil, pkg.Conflict("duplicate_resource", util.ERR_LICENCIA_DUPLICADA)
				}
			}
			// Si no se puede determinar el campo, devolver mensaje genérico
			return nil, pkg.Conflict("duplicate_resource", errApp.Message)
		}
		return nil, pkg.Internal(util.MSG_CREATE_ERROR, err)
	}
	return conductor, nil
}

func (s *service) Update(ctx context.Context, id int64, input input.UpdateConductorInput) (*domain.Conductor, error) {
	if id <= 0 {
		return nil, pkg.BadRequest("invalid_conductor_id", util.MSG_INVALID_ID)
	}
	pkg.TrimSpacesOnStruct(&input)
	if err := input.Validate(); err != nil {
		return nil, err
	}
	conductor, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkg.NotFound("conductor_not_found", util.MSG_NOT_FOUND)
		}
		return nil, pkg.Internal(util.MSG_UPDATE_ERROR, err)
	}
	if conductor == nil {
		return nil, pkg.NotFound("conductor_not_found", util.MSG_NOT_FOUND)
	}
	if input.Nombres != nil {
		conductor.Nombres = pkg.CapitalizeWords(*input.Nombres)
	}
	if input.Apellidos != nil {
		conductor.Apellidos = pkg.CapitalizeWords(*input.Apellidos)
	}
	if input.DNI != nil {
		conductor.DNI = *input.DNI
	}
	if input.NumeroLicencia != nil {
		conductor.NumeroLicencia = *input.NumeroLicencia
	}
	if input.Telefono != nil {
		conductor.Telefono = *input.Telefono
	}
	if input.Direccion != nil {
		conductor.Direccion = input.Direccion
	}
	if err := s.repo.Update(conductor); err != nil {
		return nil, pkg.Internal(util.MSG_UPDATE_ERROR, err)
	}
	return conductor, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return pkg.BadRequest("invalid_conductor_id", util.MSG_INVALID_ID)
	}
	err := s.repo.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return pkg.NotFound("conductor_not_found", util.MSG_NOT_FOUND)
		}
		return pkg.Internal(util.MSG_DELETE_ERROR, err)
	}
	return nil
}
