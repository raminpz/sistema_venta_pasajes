package service

import (
	"errors"
	"net/http"
	"sistema_venta_pasajes/internal/pago/domain"
	"sistema_venta_pasajes/internal/pago/input"
	"sistema_venta_pasajes/internal/pago/repository"
	"sistema_venta_pasajes/internal/pago/util"
	"sistema_venta_pasajes/pkg"

	"gorm.io/gorm"
)

type PagoService interface {
	Create(in input.CreatePagoInput) (*input.PagoOutput, error)
	Update(id int64, in input.UpdatePagoInput) (*input.PagoOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*input.PagoOutput, error)
	List(page, size int, idVenta *int64) ([]input.PagoOutput, int, error)
}

type pagoService struct {
	repo repository.PagoRepository
}

func NewPagoService(repo repository.PagoRepository) PagoService {
	return &pagoService{repo: repo}
}

const dateTimeLayout = "2006-01-02 15:04:05"

func toOutput(p *domain.Pago) *input.PagoOutput {
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
	return &input.PagoOutput{
		IDPago:    p.IDPago,
		IDVenta:   p.IDVenta,
		IDMetodo:  p.IDMetodo,
		Monto:     p.Monto,
		Estado:    p.Estado,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (s *pagoService) Create(in input.CreatePagoInput) (*input.PagoOutput, error) {
	pkg.TrimSpacesOnStruct(&in)

	// Validar con mensajes específicos por campo
	if err := validateCreatePagoInput(in); err != nil {
		return nil, err
	}

	estado := in.Estado
	if estado == "" {
		estado = util.STATUS_REGISTRADA
	}

	pago := &domain.Pago{
		IDVenta:  in.IDVenta,
		IDMetodo: in.IDMetodo,
		Monto:    in.Monto,
		Estado:   estado,
	}
	if err := s.repo.Create(pago); err != nil {
		return nil, util.ParseDBError(err, util.ERR_CODE_CREATE, util.MSG_PAGO_CREATE_ERROR)
	}

	return toOutput(pago), nil
}

func (s *pagoService) Update(id int64, in input.UpdatePagoInput) (*input.PagoOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PAGO_INVALID_ID)
	}
	pkg.TrimSpacesOnStruct(&in)

	pago, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_PAGO_NOT_FOUND)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_UPDATE, util.MSG_PAGO_UPDATE_ERROR).WithCause(err)
	}

	// Validar cambios con mensajes específicos
	if err := validateUpdatePagoInput(in); err != nil {
		return nil, err
	}

	if in.IDMetodo != nil {
		pago.IDMetodo = *in.IDMetodo
	}
	if in.Monto != nil {
		pago.Monto = *in.Monto
	}
	if in.Estado != nil {
		pago.Estado = *in.Estado
	}

	if err := s.repo.Update(pago); err != nil {
		return nil, util.ParseDBError(err, util.ERR_CODE_UPDATE, util.MSG_PAGO_UPDATE_ERROR)
	}
	return toOutput(pago), nil
}

func (s *pagoService) Delete(id int64) error {
	if id <= 0 {
		return pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PAGO_INVALID_ID)
	}
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_PAGO_NOT_FOUND)
		}
		return pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_DELETE, util.MSG_PAGO_DELETE_ERROR).WithCause(err)
	}
	return nil
}

func (s *pagoService) GetByID(id int64) (*input.PagoOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PAGO_INVALID_ID)
	}
	pago, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_PAGO_NOT_FOUND)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_NOT_FOUND, util.MSG_PAGO_NOT_FOUND).WithCause(err)
	}
	return toOutput(pago), nil
}

func (s *pagoService) List(page, size int, idVenta *int64) ([]input.PagoOutput, int, error) {
	offset, limit, _ := pkg.Paginate(page, size, 0)
	pagos, total, err := s.repo.List(offset, limit, idVenta)
	if err != nil {
		return nil, 0, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_LIST, util.MSG_PAGO_LIST_ERROR).WithCause(err)
	}
	outs := make([]input.PagoOutput, 0, len(pagos))
	for i := range pagos {
		outs = append(outs, *toOutput(&pagos[i]))
	}
	return outs, total, nil
}

// Validación para Create con mensajes específicos por campo
func validateCreatePagoInput(in input.CreatePagoInput) error {
	details := map[string]string{}

	if in.IDVenta <= 0 {
		details["id_venta"] = "El ID de la venta es obligatorio y debe ser mayor a 0"
	}
	if in.IDMetodo <= 0 {
		details["id_metodo"] = "El ID del metodo de pago es obligatorio y debe ser mayor a 0"
	}
	if in.Monto < 0 {
		details["monto"] = "El monto no puede ser negativo"
	}
	if in.Estado != "" && !util.IsValidEstado(in.Estado) {
		details["estado"] = "Estado inválido. Use: REGISTRADA, PARCIAL, PAGADA o ANULADA"
	}

	if len(details) > 0 {
		return pkg.Validation(util.MSG_PAGO_VALIDATION, details)
	}
	return nil
}

// Validación para Update con mensajes específicos por campo
func validateUpdatePagoInput(in input.UpdatePagoInput) error {
	details := map[string]string{}

	if in.IDMetodo != nil && *in.IDMetodo <= 0 {
		details["id_metodo"] = "El ID del metodo de pago debe ser mayor a 0"
	}
	if in.Monto != nil && *in.Monto < 0 {
		details["monto"] = "El monto no puede ser negativo"
	}
	if in.Estado != nil && !util.IsValidEstado(*in.Estado) {
		details["estado"] = "Estado inválido. Use: REGISTRADA, PARCIAL, PAGADA o ANULADA"
	}

	if len(details) > 0 {
		return pkg.Validation(util.MSG_PAGO_VALIDATION, details)
	}
	return nil
}
