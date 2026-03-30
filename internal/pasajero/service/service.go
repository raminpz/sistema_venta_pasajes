package service

import (
	"errors"
	"sistema_venta_pasajes/internal/pasajero/domain"
	"sistema_venta_pasajes/internal/pasajero/input"
	"sistema_venta_pasajes/internal/pasajero/repository"
	"sistema_venta_pasajes/internal/pasajero/util"
	"sistema_venta_pasajes/pkg"
	"strings"
	"time"
)

type PasajeroService interface {
	Create(in input.CreatePasajeroInput) (input.PasajeroOutput, error)
	Update(id int64, in input.UpdatePasajeroInput) (input.PasajeroOutput, error)
	Delete(id int64) error
	GetByID(id int64) (input.PasajeroOutput, error)
	List(page, size int) ([]input.PasajeroOutput, pkg.PaginationMeta, error)
	Search(query string) ([]input.PasajeroOutput, error)
}

type pasajeroService struct {
	repo repository.PasajeroRepository
}

func NewPasajeroService(repo repository.PasajeroRepository) PasajeroService {
	return &pasajeroService{repo: repo}
}

func toOutput(p *domain.Pasajero) input.PasajeroOutput {
	var fechaNac *string
	if p.FechaNacimiento != nil {
		f := p.FechaNacimiento.Format("2006-01-02")
		fechaNac = &f
	}
	return input.PasajeroOutput{
		IDPasajero:      int64(p.IDPasajero),
		TipoDocumento:   p.TipoDocumento,
		NroDocumento:    p.NroDocumento,
		Nombres:         p.Nombres,
		Apellidos:       p.Apellidos,
		Telefono:        p.Telefono,
		Email:           p.Email,
		FechaNacimiento: fechaNac,
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *pasajeroService) Create(in input.CreatePasajeroInput) (input.PasajeroOutput, error) {
	pkg.TrimSpacesOnStruct(&in)
	if err := in.Validate(); err != nil {
		return input.PasajeroOutput{}, err
	}
	var fechaNac *time.Time
	if in.FechaNacimiento != nil && *in.FechaNacimiento != "" {
		f, err := time.Parse("2006-01-02", *in.FechaNacimiento)
		if err != nil {
			return input.PasajeroOutput{}, errors.New(util.ERR_DATE_FORMAT)
		}
		fechaNac = &f
	}
	// Capitalizar nombres y apellidos usando pkg
	nombres := pkg.CapitalizeWords(in.Nombres)
	apellidos := pkg.CapitalizeWords(in.Apellidos)
	pasajero := &domain.Pasajero{
		TipoDocumento:   in.TipoDocumento,
		NroDocumento:    in.NroDocumento,
		Nombres:         nombres,
		Apellidos:       apellidos,
		Telefono:        in.Telefono,
		Email:           in.Email,
		FechaNacimiento: fechaNac,
	}
	err := s.repo.Create(pasajero)
	if err != nil {
		errApp := pkg.AsAppError(err)
		if errApp != nil && errApp.Code == "duplicate_resource" {
			if errApp.Err != nil {
				errStr := errApp.Err.Error()
				if strings.Contains(errStr, "nrodocumento") || strings.Contains(errStr, "NroDocumento") || strings.Contains(errStr, "dni") || strings.Contains(errStr, "DNI") {
					return input.PasajeroOutput{}, pkg.Conflict("duplicate_resource", util.MSG_DUPLICATE_DNI)
				}
			}
			return input.PasajeroOutput{}, pkg.Conflict("duplicate_resource", errApp.Message)
		}
		return input.PasajeroOutput{}, pkg.Internal(util.MSG_CREATE_ERROR, err)
	}
	return toOutput(pasajero), nil
}

func (s *pasajeroService) Update(id int64, in input.UpdatePasajeroInput) (input.PasajeroOutput, error) {
	pkg.TrimSpacesOnStruct(&in)
	if err := in.Validate(); err != nil {
		return input.PasajeroOutput{}, err
	}
	pasajero, err := s.repo.GetByID(id)
	if err != nil {
		return input.PasajeroOutput{}, err
	}
	if in.TipoDocumento != "" {
		pasajero.TipoDocumento = in.TipoDocumento
	}
	if in.NroDocumento != "" {
		pasajero.NroDocumento = in.NroDocumento
	}
	if in.Nombres != "" {
		pasajero.Nombres = pkg.CapitalizeWords(in.Nombres)
	}
	if in.Apellidos != "" {
		pasajero.Apellidos = pkg.CapitalizeWords(in.Apellidos)
	}
	if in.Telefono != "" {
		pasajero.Telefono = in.Telefono
	}
	if in.Email != nil {
		pasajero.Email = in.Email
	}
	if in.FechaNacimiento != nil && *in.FechaNacimiento != "" {
		f, err := time.Parse("2006-01-02", *in.FechaNacimiento)
		if err != nil {
			return input.PasajeroOutput{}, errors.New(util.ERR_DATE_FORMAT)
		}
		pasajero.FechaNacimiento = &f
	}
	err = s.repo.Update(pasajero)
	if err != nil {
		return input.PasajeroOutput{}, err
	}
	return toOutput(pasajero), nil
}

func (s *pasajeroService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *pasajeroService) GetByID(id int64) (input.PasajeroOutput, error) {
	pasajero, err := s.repo.GetByID(id)
	if err != nil {
		return input.PasajeroOutput{}, err
	}
	return toOutput(pasajero), nil
}

func (s *pasajeroService) List(page, size int) ([]input.PasajeroOutput, pkg.PaginationMeta, error) {
	pasajeros, total, err := s.repo.List(page, size)
	if err != nil {
		return nil, pkg.PaginationMeta{}, err
	}
	offset, limit, meta := pkg.Paginate(page, size, total)
	_ = offset
	_ = limit
	var out []input.PasajeroOutput
	for _, p := range pasajeros {
		pCopy := p
		out = append(out, toOutput(&pCopy))
	}
	return out, meta, nil
}

func (s *pasajeroService) Search(query string) ([]input.PasajeroOutput, error) {
	pasajeros, _, err := s.repo.Search(query)
	if err != nil {
		return nil, err
	}
	var out []input.PasajeroOutput
	for _, p := range pasajeros {
		pCopy := p
		out = append(out, toOutput(&pCopy))
	}
	return out, nil
}
