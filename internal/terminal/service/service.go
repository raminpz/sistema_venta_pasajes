package service

import (
	"errors"
	"sistema_venta_pasajes/internal/terminal/domain"
	input "sistema_venta_pasajes/internal/terminal/input"
	"sistema_venta_pasajes/internal/terminal/repository"
	"sistema_venta_pasajes/internal/terminal/util"
	"sistema_venta_pasajes/pkg"

	"gorm.io/gorm"
)

type TerminalService interface {
	Create(input input.CreateTerminalInput) (*domain.Terminal, error)
	GetByID(id int64) (*domain.Terminal, error)
	Update(id int64, input input.UpdateTerminalInput) (*domain.Terminal, error)
	Delete(id int64) error
	List() ([]domain.Terminal, error)
}

type terminalService struct {
	repo repository.TerminalRepository
}

func NewTerminalService(repo repository.TerminalRepository) TerminalService {
	return &terminalService{repo: repo}
}

func (s *terminalService) Create(in input.CreateTerminalInput) (*domain.Terminal, error) {
	pkg.TrimSpacesOnStruct(&in)
	if err := validateTerminalInput(in); err != nil {
		return nil, err
	}
	terminal := &domain.Terminal{
		NOMBRE:       in.Nombre,
		CIUDAD:       in.Ciudad,
		DEPARTAMENTO: in.Departamento,
		DIRECCION:    in.Direccion,
		ESTADO:       in.Estado,
	}
	err := s.repo.Create(terminal)
	if err != nil {
		errApp := pkg.AsAppError(err)
		if errApp != nil && errApp.Code == "duplicate_resource" {
			return nil, pkg.Conflict(util.ERR_CODE_DUPLICATE_RESOURCE, util.MSG_TERMINAL_DUPLICATE)
		}
		return nil, pkg.Internal(util.MSG_TERMINAL_CREATE_ERROR, err)
	}
	return terminal, nil
}

func (s *terminalService) GetByID(id int64) (*domain.Terminal, error) {
	return s.repo.GetByID(id)
}

func (s *terminalService) Update(id int64, in input.UpdateTerminalInput) (*domain.Terminal, error) {
	terminal, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	pkg.TrimSpacesOnStruct(&in)
	if err := validateTerminalUpdateInput(in); err != nil {
		return nil, err
	}
	terminal.NOMBRE = in.Nombre
	terminal.CIUDAD = in.Ciudad
	terminal.DEPARTAMENTO = in.Departamento
	terminal.DIRECCION = in.Direccion
	terminal.ESTADO = in.Estado
	err = s.repo.Update(terminal)
	return terminal, err
}

// Validación de campos obligatorios y longitudes para Create
func validateTerminalInput(in input.CreateTerminalInput) error {
	details := map[string]string{}
	if in.Nombre == "" {
		details["nombre"] = "El nombre es obligatorio"
	}
	if in.Ciudad == "" {
		details["ciudad"] = "La ciudad es obligatoria"
	}
	if in.Departamento == "" {
		details["departamento"] = "El departamento es obligatorio"
	}
	if in.Direccion == "" {
		details["direccion"] = "La dirección es obligatoria"
	}
	if in.Estado == "" {
		details["estado"] = "El estado es obligatorio"
	}
	if len(details) > 0 {
		return pkg.Validation(util.MSG_TERMINAL_VALIDATION, details)
	}
	return nil
}

// Validación para Update (puedes ajustar si quieres permitir campos vacíos)
func validateTerminalUpdateInput(in input.UpdateTerminalInput) error {
	details := map[string]string{}
	if in.Nombre != "" && len(in.Nombre) > 100 {
		details["nombre"] = "El nombre no puede exceder 100 caracteres"
	}
	if in.Ciudad != "" && len(in.Ciudad) > 100 {
		details["ciudad"] = "La ciudad no puede exceder 100 caracteres"
	}
	if in.Departamento != "" && len(in.Departamento) > 100 {
		details["departamento"] = "El departamento no puede exceder 100 caracteres"
	}
	if in.Direccion != "" && len(in.Direccion) > 200 {
		details["direccion"] = "La dirección no puede exceder 200 caracteres"
	}
	if len(details) > 0 {
		return pkg.Validation(util.MSG_TERMINAL_VALIDATION, details)
	}
	return nil
}

func (s *terminalService) Delete(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_TERMINAL_NOT_FOUND)
		}
		return pkg.Internal(util.MSG_TERMINAL_DELETE_ERROR, err)
	}
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_TERMINAL_NOT_FOUND)
		}
		return pkg.Internal(util.MSG_TERMINAL_DELETE_ERROR, err)
	}
	return nil
}

func (s *terminalService) List() ([]domain.Terminal, error) {
	return s.repo.List()
}
