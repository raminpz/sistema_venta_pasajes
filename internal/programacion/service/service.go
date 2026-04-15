package service

import (
	"errors"
	"net/http"
	"sistema_venta_pasajes/internal/programacion/domain"
	"sistema_venta_pasajes/internal/programacion/input"
	"sistema_venta_pasajes/internal/programacion/repository"
	"sistema_venta_pasajes/internal/programacion/util"
	"sistema_venta_pasajes/pkg"
	"time"

	"gorm.io/gorm"
)

type ProgramacionService interface {
	Create(in input.CreateProgramacionInput) (*input.ProgramacionOutput, error)
	Update(id int64, in input.UpdateProgramacionInput) (*input.ProgramacionOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*input.ProgramacionOutput, error)
	List(page, size int) ([]input.ProgramacionOutput, int, error)
}

type programacionService struct {
	repo repository.ProgramacionRepository
}

func NewProgramacionService(repo repository.ProgramacionRepository) ProgramacionService {
	return &programacionService{repo: repo}
}

const dateTimeLayout = "2006-01-02 15:04:05"

func toOutput(p *domain.Programacion) *input.ProgramacionOutput {
	var fechaLlegada *string
	if p.FechaLlegada != nil {
		v := p.FechaLlegada.Format(dateTimeLayout)
		fechaLlegada = &v
	}

	var createdAt *string
	if p.CreatedAt != nil {
		v := p.CreatedAt.Format(dateTimeLayout)
		createdAt = &v
	}

	var updatedAt *string
	if p.UpdatedAt != nil {
		v := p.UpdatedAt.Format(dateTimeLayout)
		updatedAt = &v
	}

	return &input.ProgramacionOutput{
		IDProgramacion: p.IDProgramacion,
		IDRuta:         p.IDRuta,
		IDVehiculo:     p.IDVehiculo,
		IDConductor:    p.IDConductor,
		FechaSalida:    p.FechaSalida.Format(dateTimeLayout),
		FechaLlegada:   fechaLlegada,
		Estado:         p.Estado,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func (s *programacionService) Create(in input.CreateProgramacionInput) (*input.ProgramacionOutput, error) {
	pkg.TrimSpacesOnStruct(&in)

	if in.IDRuta <= 0 || in.IDVehiculo <= 0 || in.IDConductor <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_CREATE, util.MSG_PROGRAMACION_REQUIRED_IDS)
	}

	fechaSalida, err := util.ParseDateTime(in.FechaSalida)
	if err != nil {
		return nil, pkg.BadRequest(util.ERR_CODE_CREATE, util.MSG_PROGRAMACION_INVALID_DATETIME)
	}

	var fechaLlegada *timeRef
	if in.FechaLlegada != nil && *in.FechaLlegada != "" {
		llegada, err := util.ParseDateTime(*in.FechaLlegada)
		if err != nil {
			return nil, pkg.BadRequest(util.ERR_CODE_CREATE, util.MSG_PROGRAMACION_INVALID_DATETIME)
		}
		fechaLlegada = &timeRef{v: llegada}
	}

	estado := in.Estado
	if estado == "" {
		estado = util.STATUS_PROGRAMADO
	}
	if !util.IsValidEstado(estado) {
		return nil, pkg.BadRequest(util.ERR_CODE_CREATE, util.MSG_PROGRAMACION_INVALID_STATUS)
	}

	if fechaLlegada != nil && fechaLlegada.v.Before(fechaSalida) {
		return nil, pkg.BadRequest(util.ERR_CODE_CREATE, util.MSG_PROGRAMACION_INVALID_DATES)
	}

	programacion := &domain.Programacion{
		IDRuta:      in.IDRuta,
		IDVehiculo:  in.IDVehiculo,
		IDConductor: in.IDConductor,
		FechaSalida: fechaSalida,
		Estado:      estado,
	}
	if fechaLlegada != nil {
		programacion.FechaLlegada = &fechaLlegada.v
	}

	if err := s.repo.Create(programacion); err != nil {
		appErr := pkg.AsAppError(err)
		if appErr != nil && appErr.Code == "foreign_key_conflict" {
			return nil, pkg.NewAppError(http.StatusConflict, util.ERR_CODE_CREATE, util.MSG_PROGRAMACION_FK_CONFLICT)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_CREATE, util.MSG_PROGRAMACION_CREATE_ERROR).WithCause(err)
	}

	return toOutput(programacion), nil
}

type timeRef struct {
	v time.Time
}

func (s *programacionService) Update(id int64, in input.UpdateProgramacionInput) (*input.ProgramacionOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PROGRAMACION_INVALID_ID)
	}

	pkg.TrimSpacesOnStruct(&in)

	current, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_PROGRAMACION_NOT_FOUND)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_UPDATE_ERROR).WithCause(err)
	}

	if in.IDRuta != nil {
		if *in.IDRuta <= 0 {
			return nil, pkg.BadRequest(util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_REQUIRED_IDS)
		}
		current.IDRuta = *in.IDRuta
	}
	if in.IDVehiculo != nil {
		if *in.IDVehiculo <= 0 {
			return nil, pkg.BadRequest(util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_REQUIRED_IDS)
		}
		current.IDVehiculo = *in.IDVehiculo
	}
	if in.IDConductor != nil {
		if *in.IDConductor <= 0 {
			return nil, pkg.BadRequest(util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_REQUIRED_IDS)
		}
		current.IDConductor = *in.IDConductor
	}
	if in.FechaSalida != nil {
		t, err := util.ParseDateTime(*in.FechaSalida)
		if err != nil {
			return nil, pkg.BadRequest(util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_INVALID_DATETIME)
		}
		current.FechaSalida = t
	}
	if in.FechaLlegada != nil {
		if *in.FechaLlegada == "" {
			current.FechaLlegada = nil
		} else {
			t, err := util.ParseDateTime(*in.FechaLlegada)
			if err != nil {
				return nil, pkg.BadRequest(util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_INVALID_DATETIME)
			}
			current.FechaLlegada = &t
		}
	}
	if in.Estado != nil {
		if !util.IsValidEstado(*in.Estado) {
			return nil, pkg.BadRequest(util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_INVALID_STATUS)
		}
		current.Estado = *in.Estado
	}

	if current.FechaLlegada != nil && current.FechaLlegada.Before(current.FechaSalida) {
		return nil, pkg.BadRequest(util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_INVALID_DATES)
	}

	if err := s.repo.Update(current); err != nil {
		appErr := pkg.AsAppError(err)
		if appErr != nil && appErr.Code == "foreign_key_conflict" {
			return nil, pkg.NewAppError(http.StatusConflict, util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_FK_CONFLICT)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_UPDATE, util.MSG_PROGRAMACION_UPDATE_ERROR).WithCause(err)
	}

	return toOutput(current), nil
}

func (s *programacionService) Delete(id int64) error {
	if id <= 0 {
		return pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PROGRAMACION_INVALID_ID)
	}
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_PROGRAMACION_NOT_FOUND)
		}
		return pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_DELETE, util.MSG_PROGRAMACION_DELETE_ERROR).WithCause(err)
	}
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_PROGRAMACION_NOT_FOUND)
		}
		return pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_DELETE, util.MSG_PROGRAMACION_DELETE_ERROR).WithCause(err)
	}
	return nil
}

func (s *programacionService) GetByID(id int64) (*input.ProgramacionOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PROGRAMACION_INVALID_ID)
	}
	programacion, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_PROGRAMACION_NOT_FOUND)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_NOT_FOUND, util.MSG_PROGRAMACION_NOT_FOUND).WithCause(err)
	}
	return toOutput(programacion), nil
}

func (s *programacionService) List(page, size int) ([]input.ProgramacionOutput, int, error) {
	offset, limit, _ := pkg.Paginate(page, size, 0)
	programaciones, total, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, 0, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_LIST, util.MSG_PROGRAMACION_LIST_ERROR).WithCause(err)
	}
	outs := make([]input.ProgramacionOutput, 0, len(programaciones))
	for i := range programaciones {
		outs = append(outs, *toOutput(&programaciones[i]))
	}
	return outs, total, nil
}
