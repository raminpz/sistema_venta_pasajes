package service

import (
	"errors"
	"net/http"
	"sistema_venta_pasajes/internal/encomienda/domain"
	"sistema_venta_pasajes/internal/encomienda/input"
	"sistema_venta_pasajes/internal/encomienda/repository"
	"sistema_venta_pasajes/internal/encomienda/util"
	"sistema_venta_pasajes/pkg"

	"gorm.io/gorm"
)

type EncomiendaService interface {
	Create(in input.CreateEncomiendaInput) (*input.EncomiendaOutput, error)
	Update(id int64, in input.UpdateEncomiendaInput) (*input.EncomiendaOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*input.EncomiendaOutput, error)
	List(page, size int) ([]input.EncomiendaOutput, int, error)
}

type encomiendaService struct {
	repo repository.EncomiendaRepository
}

func NewEncomiendaService(repo repository.EncomiendaRepository) EncomiendaService {
	return &encomiendaService{repo: repo}
}

const dateTimeLayout = "2006-01-02 15:04:05"

func toOutput(e *domain.Encomienda) *input.EncomiendaOutput {
	var createdAt *string
	if e.CreatedAt != nil {
		v := e.CreatedAt.Format(dateTimeLayout)
		createdAt = &v
	}
	var updatedAt *string
	if e.UpdatedAt != nil {
		v := e.UpdatedAt.Format(dateTimeLayout)
		updatedAt = &v
	}
	return &input.EncomiendaOutput{
		IDEncomienda:       e.IDEncomienda,
		IDVenta:            e.IDVenta,
		IDProgramacion:     e.IDProgramacion,
		Descripcion:        e.Descripcion,
		PesoKg:             e.PesoKg,
		Costo:              e.Costo,
		RemitenteNombre:    e.RemitenteNombre,
		RemitenteDoc:       e.RemitenteDoc,
		DestinatarioNombre: e.DestinatarioNombre,
		DestinatarioDoc:    e.DestinatarioDoc,
		DestinatarioTel:    e.DestinatarioTel,
		Estado:             e.Estado,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}

func (s *encomiendaService) Create(in input.CreateEncomiendaInput) (*input.EncomiendaOutput, error) {
	pkg.TrimSpacesOnStruct(&in)
	in.RemitenteNombre = pkg.CapitalizeWords(in.RemitenteNombre)
	in.DestinatarioNombre = pkg.CapitalizeWords(in.DestinatarioNombre)

	if err := validateCreateEncomiendaInput(in); err != nil {
		return nil, err
	}

	if in.Estado == "" {
		in.Estado = util.STATUS_PENDIENTE
	}

	encomienda := &domain.Encomienda{
		IDVenta:            in.IDVenta,
		IDProgramacion:     in.IDProgramacion,
		Descripcion:        in.Descripcion,
		PesoKg:             in.PesoKg,
		Costo:              in.Costo,
		RemitenteNombre:    in.RemitenteNombre,
		RemitenteDoc:       in.RemitenteDoc,
		DestinatarioNombre: in.DestinatarioNombre,
		DestinatarioDoc:    in.DestinatarioDoc,
		DestinatarioTel:    in.DestinatarioTel,
		Estado:             in.Estado,
	}

	if err := s.repo.Create(encomienda); err != nil {
		return nil, util.ParseDBError(err, util.ERR_CODE_CREATE, util.MSG_ENCOMIENDA_CREATE_ERROR)
	}

	return toOutput(encomienda), nil
}

func (s *encomiendaService) Update(id int64, in input.UpdateEncomiendaInput) (*input.EncomiendaOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_ENCOMIENDA_INVALID_ID)
	}

	pkg.TrimSpacesOnStruct(&in)
	in.RemitenteNombre = pkg.CapitalizeWords(in.RemitenteNombre)
	in.DestinatarioNombre = pkg.CapitalizeWords(in.DestinatarioNombre)

	if err := validateUpdateEncomiendaInput(in); err != nil {
		return nil, err
	}

	encomienda, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_ENCOMIENDA_NOT_FOUND)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_UPDATE, util.MSG_ENCOMIENDA_UPDATE_ERROR).WithCause(err)
	}

	encomienda.IDVenta = in.IDVenta
	encomienda.IDProgramacion = in.IDProgramacion
	encomienda.Descripcion = in.Descripcion
	encomienda.PesoKg = in.PesoKg
	encomienda.Costo = in.Costo
	encomienda.RemitenteNombre = in.RemitenteNombre
	encomienda.RemitenteDoc = in.RemitenteDoc
	encomienda.DestinatarioNombre = in.DestinatarioNombre
	encomienda.DestinatarioDoc = in.DestinatarioDoc
	encomienda.DestinatarioTel = in.DestinatarioTel
	encomienda.Estado = in.Estado

	if err := s.repo.Update(encomienda); err != nil {
		return nil, util.ParseDBError(err, util.ERR_CODE_UPDATE, util.MSG_ENCOMIENDA_UPDATE_ERROR)
	}

	return toOutput(encomienda), nil
}

func (s *encomiendaService) Delete(id int64) error {
	if id <= 0 {
		return pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_ENCOMIENDA_INVALID_ID)
	}

	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_ENCOMIENDA_NOT_FOUND)
		}
		return util.ParseDBError(err, util.ERR_CODE_DELETE, util.MSG_ENCOMIENDA_DELETE_ERROR)
	}

	return nil
}

func (s *encomiendaService) GetByID(id int64) (*input.EncomiendaOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_ENCOMIENDA_INVALID_ID)
	}

	encomienda, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_ENCOMIENDA_NOT_FOUND)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_NOT_FOUND, util.MSG_ENCOMIENDA_NOT_FOUND).WithCause(err)
	}

	return toOutput(encomienda), nil
}

func (s *encomiendaService) List(page, size int) ([]input.EncomiendaOutput, int, error) {
	offset, limit, _ := pkg.Paginate(page, size, 0)
	encomiendas, total, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, 0, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_LIST, util.MSG_ENCOMIENDA_LIST_ERROR).WithCause(err)
	}

	outs := make([]input.EncomiendaOutput, 0, len(encomiendas))
	for i := range encomiendas {
		outs = append(outs, *toOutput(&encomiendas[i]))
	}
	return outs, total, nil
}

// Validación para Create con mensajes específicos por campo
func validateCreateEncomiendaInput(in input.CreateEncomiendaInput) error {
	details := map[string]string{}

	if in.IDVenta <= 0 {
		details["id_venta"] = "El ID de la venta es obligatorio y debe ser mayor a 0"
	}
	if in.IDProgramacion <= 0 {
		details["id_programacion"] = "El ID de la programación es obligatorio y debe ser mayor a 0"
	}
	if in.RemitenteNombre == "" {
		details["remitente_nombre"] = "El nombre del remitente es obligatorio"
	}
	if in.DestinatarioNombre == "" {
		details["destinatario_nombre"] = "El nombre del destinatario es obligatorio"
	}
	if in.Costo <= 0 {
		details["costo"] = "El costo es obligatorio y debe ser mayor a 0"
	}

	// Campos obligatorios adicionales
	if in.RemitenteDoc == "" {
		details["remitente_doc"] = "El documento del remitente es obligatorio"
	} else if len(in.RemitenteDoc) > 20 {
		details["remitente_doc"] = "El documento del remitente no puede exceder 20 caracteres"
	}
	if in.DestinatarioTel == "" {
		details["destinatario_tel"] = "El telefono del destinatario es obligatorio"
	} else if len(in.DestinatarioTel) != 9 {
		details["destinatario_tel"] = "El telefono del destinatario debe tener 9 digitos"
	} else if !isDigits(in.DestinatarioTel) {
		details["destinatario_tel"] = "El telefono del destinatario debe contener solo numeros"
	}

	// Validaciones de campos opcionales
	if in.PesoKg != nil && *in.PesoKg <= 0 {
		details["peso_kg"] = "El peso debe ser mayor a 0"
	}
	if in.DestinatarioDoc != nil && *in.DestinatarioDoc == "" {
		details["destinatario_doc"] = "El documento del destinatario no puede estar vacio"
	}
	if in.Descripcion != "" && len(in.Descripcion) > 300 {
		details["descripcion"] = "La descripción no puede exceder 300 caracteres"
	}
	if in.Estado != "" && !util.IsValidEstado(in.Estado) {
		details["estado"] = "Estado inválido. Use: PENDIENTE, ENTREGADO o DEVUELTO"
	}

	if len(details) > 0 {
		return pkg.Validation(util.MSG_ENCOMIENDA_VALIDATION, details)
	}
	return nil
}

// Validación para Update con mensajes específicos por campo
func validateUpdateEncomiendaInput(in input.UpdateEncomiendaInput) error {
	details := map[string]string{}

	if in.IDVenta <= 0 {
		details["id_venta"] = "El ID de la venta es obligatorio y debe ser mayor a 0"
	}
	if in.IDProgramacion <= 0 {
		details["id_programacion"] = "El ID de la programación es obligatorio y debe ser mayor a 0"
	}
	if in.RemitenteNombre == "" {
		details["remitente_nombre"] = "El nombre del remitente es obligatorio"
	}
	if in.DestinatarioNombre == "" {
		details["destinatario_nombre"] = "El nombre del destinatario es obligatorio"
	}
	if in.Costo <= 0 {
		details["costo"] = "El costo es obligatorio y debe ser mayor a 0"
	}

	// Campos obligatorios adicionales
	if in.RemitenteDoc == "" {
		details["remitente_doc"] = "El documento del remitente es obligatorio"
	} else if len(in.RemitenteDoc) > 20 {
		details["remitente_doc"] = "El documento del remitente no puede exceder 20 caracteres"
	}
	if in.DestinatarioTel == "" {
		details["destinatario_tel"] = "El telefono del destinatario es obligatorio"
	} else if len(in.DestinatarioTel) != 9 {
		details["destinatario_tel"] = "El telefono del destinatario debe tener 9 digitos"
	} else if !isDigits(in.DestinatarioTel) {
		details["destinatario_tel"] = "El telefono del destinatario debe contener solo numeros"
	}

	// Validaciones de campos opcionales
	if in.PesoKg != nil && *in.PesoKg <= 0 {
		details["peso_kg"] = "El peso debe ser mayor a 0"
	}
	if in.DestinatarioDoc != nil && *in.DestinatarioDoc == "" {
		details["destinatario_doc"] = "El documento del destinatario no puede estar vacio"
	}
	if in.Descripcion != "" && len(in.Descripcion) > 300 {
		details["descripcion"] = "La descripción no puede exceder 300 caracteres"
	}
	if in.Estado != "" && !util.IsValidEstado(in.Estado) {
		details["estado"] = "Estado inválido. Use: PENDIENTE, ENTREGADO o DEVUELTO"
	}

	if len(details) > 0 {
		return pkg.Validation(util.MSG_ENCOMIENDA_VALIDATION, details)
	}
	return nil
}

func isDigits(v string) bool {
	if v == "" {
		return false
	}
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
