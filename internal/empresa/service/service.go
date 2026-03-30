package service

import (
	"sistema_venta_pasajes/internal/empresa/domain"
	"sistema_venta_pasajes/internal/empresa/input"
	"sistema_venta_pasajes/internal/empresa/repository"
	"sistema_venta_pasajes/internal/empresa/util"
	"sistema_venta_pasajes/pkg"
	"strings"
	"time"
)

type EmpresaService interface {
	Create(in input.CreateEmpresaInput) (input.EmpresaOutput, error)
	Update(id int64, in input.UpdateEmpresaInput) (input.EmpresaOutput, error)
	Delete(id int64) error
	GetByID(id int64) (input.EmpresaOutput, error)
	List() ([]input.EmpresaOutput, error)
}

type empresaService struct {
	repo repository.EmpresaRepository
}

func NewEmpresaService(repo repository.EmpresaRepository) EmpresaService {
	return &empresaService{repo: repo}
}

func toOutput(e *domain.Empresa) input.EmpresaOutput {
	fechaCreacion := e.FechaCreacion.Format("2006-01-02")
	return input.EmpresaOutput{
		IDEmpresa:       int64(e.IDEmpresa),
		RUC:             e.RUC,
		RazonSocial:     e.RazonSocial,
		NombreComercial: e.NombreComercial,
		Direccion:       e.Direccion,
		Telefono:        e.Telefono,
		Email:           e.Email,
		Logo:            e.Logo,
		FechaCreacion:   fechaCreacion,
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       e.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *empresaService) Create(in input.CreateEmpresaInput) (input.EmpresaOutput, error) {
	   pkg.TrimSpacesOnStruct(&in)
	   if err := in.Validate(); err != nil {
			   return input.EmpresaOutput{}, pkg.Validation(err.Error(), nil)
	   }
	   in.RazonSocial = pkg.CapitalizeWords(in.RazonSocial)
	   var fechaCreacion time.Time
	   if t, err := time.Parse("2006-01-02", in.FechaCreacion); err == nil {
		   fechaCreacion = t
	   } else {
		   return input.EmpresaOutput{}, pkg.BadRequest("invalid_date", util.ERR_DATE_FORMAT)
	   }
	   emp := &domain.Empresa{
		   RUC:             in.RUC,
		   RazonSocial:     in.RazonSocial,
		   NombreComercial: in.NombreComercial,
		   Direccion:       in.Direccion,
		   Telefono:        in.Telefono,
		   Email:           in.Email,
		   Logo:            in.Logo,
		   FechaCreacion:   fechaCreacion,
	   }
	   err := s.repo.Create(emp)
	   if err != nil {
		   errApp := pkg.AsAppError(err)
		   if errApp != nil && errApp.Code == "duplicate_resource" {
			   if errApp.Err != nil {
				   errStr := errApp.Err.Error()
				   if strings.Contains(errStr, "ruc") || strings.Contains(errStr, "RUC") {
					   return input.EmpresaOutput{}, pkg.Conflict("duplicate_resource", util.MSG_DUPLICATE_RUC)
				   }
			   }
			   return input.EmpresaOutput{}, pkg.Conflict("duplicate_resource", errApp.Message)
		   }
		   return input.EmpresaOutput{}, pkg.Internal("Error al crear empresa", err)
	   }
	   return toOutput(emp), nil
}

func (s *empresaService) Update(id int64, in input.UpdateEmpresaInput) (input.EmpresaOutput, error) {
	   pkg.TrimSpacesOnStruct(&in)
	   if err := in.Validate(); err != nil {
			   return input.EmpresaOutput{}, err
	   }
	emp, err := s.repo.GetByID(id)
	if err != nil {
		return input.EmpresaOutput{}, err
	}
	if in.RUC != nil && *in.RUC != "" {
		emp.RUC = *in.RUC
	}
	if in.RazonSocial != "" {
		emp.RazonSocial = pkg.CapitalizeWords(in.RazonSocial)
	}
	if in.NombreComercial != nil {
		emp.NombreComercial = in.NombreComercial
	}
	if in.Direccion != nil {
		emp.Direccion = in.Direccion
	}
	if in.Telefono != "" {
		emp.Telefono = in.Telefono
	}
	if in.Email != nil {
		emp.Email = in.Email
	}
	if in.Logo != nil {
		emp.Logo = in.Logo
	}
	if in.FechaCreacion != nil && *in.FechaCreacion != "" {
		if t, err := time.Parse("2006-01-02", *in.FechaCreacion); err == nil {
			emp.FechaCreacion = t
		} else {
			return input.EmpresaOutput{}, pkg.BadRequest("invalid_date", util.ERR_DATE_FORMAT)
		}
	}
	err = s.repo.Update(emp)
	if err != nil {
		return input.EmpresaOutput{}, err
	}
	return toOutput(emp), nil
}

func (s *empresaService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *empresaService) GetByID(id int64) (input.EmpresaOutput, error) {
	emp, err := s.repo.GetByID(id)
	if err != nil {
		return input.EmpresaOutput{}, err
	}
	return toOutput(emp), nil
}

func (s *empresaService) List() ([]input.EmpresaOutput, error) {
	empresas, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	outputs := make([]input.EmpresaOutput, len(empresas))
	for i, e := range empresas {
		outputs[i] = toOutput(&e)
	}
	return outputs, nil
}
