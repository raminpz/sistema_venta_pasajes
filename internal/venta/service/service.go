package service

import (
	"errors"
	"fmt"
	"net/http"
	asientoTramoService "sistema_venta_pasajes/internal/asiento_tramo/service"
	"sistema_venta_pasajes/internal/venta/domain"
	ventaInput "sistema_venta_pasajes/internal/venta/input"
	"sistema_venta_pasajes/internal/venta/repository"
	"sistema_venta_pasajes/internal/venta/util"
	"sistema_venta_pasajes/pkg"

	"gorm.io/gorm"
)

type VentaService interface {
	Create(in ventaInput.VentaCreateInput) (*ventaInput.VentaOutput, error)
	Update(id int64, in ventaInput.VentaUpdateInput) (*ventaInput.VentaOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*ventaInput.VentaOutput, error)
	List(page, size int) ([]ventaInput.VentaOutput, int, error)
}

type VentaServiceImpl struct {
	repo            repository.VentaRepository
	asientoTramoSvc asientoTramoService.AsientoTramoService
}

func NewVentaService(repo repository.VentaRepository, asientoTramoSvc asientoTramoService.AsientoTramoService) VentaService {
	return &VentaServiceImpl{repo: repo, asientoTramoSvc: asientoTramoSvc}
}

const layoutDateTime = "2006-01-02 15:04:05"

func toVentaOutput(v *domain.Venta) *ventaInput.VentaOutput {
	return &ventaInput.VentaOutput{
		IDVenta:           v.IDVenta,
		IDUsuario:         v.IDUsuario,
		IDTipoComprobante: v.IDTipoComprobante,
		IDProgramacion:    v.IDProgramacion,
		IDPasajero:        v.IDPasajero,
		IDAsiento:         v.IDAsiento,
		IDTramo:           v.IDTramo,
		Precio:            v.Precio,
		Descuento:         v.Descuento,
		Serie:             v.Serie,
		Correlativo:       v.Correlativo,
		NumeroComprobante: fmt.Sprintf("%s-%06d", v.Serie, v.Correlativo),
		Nota:              v.Nota,
		Observaciones:     v.Observaciones,
		Subtotal:          v.Subtotal,
		IGV:               v.IGV,
		Total:             v.Total,
		QRCode:            v.QRCode,
		CreatedAt:         v.CreatedAt.Format(layoutDateTime),
		UpdatedAt:         v.UpdatedAt.Format(layoutDateTime),
	}
}

func (s *VentaServiceImpl) Create(in ventaInput.VentaCreateInput) (*ventaInput.VentaOutput, error) {
	pkg.TrimSpacesOnStruct(&in)

	if err := util.ValidarCreateInput(
		in.IDUsuario, in.IDTipoComprobante,
		in.IDProgramacion, in.IDPasajero, in.IDAsiento, in.IDTramo,
		in.Precio, in.Descuento,
	); err != nil {
		return nil, err
	}

	// Validar disponibilidad del asiento en el tramo
	disponible, err := s.asientoTramoSvc.IsAsientoDisponible(in.IDAsiento, in.IDTramo)
	if err != nil {
		return nil, pkg.NewAppError(500, util.ERR_CODE_CREATE, util.MSG_VENTA_CREATE_ERROR).WithCause(err)
	}
	if !disponible {
		return nil, pkg.BadRequest(util.ERR_CODE_DUPLICATE, util.MSG_VENTA_ASIENTO_OCUPADO)
	}

	serie, err := util.SerieFromTipoComprobante(in.IDTipoComprobante)
	if err != nil {
		return nil, err
	}

	correlativo, err := s.repo.NextCorrelativo(serie)
	if err != nil {
		return nil, errors.New(util.MSG_VENTA_CORRELATIVO_ERROR)
	}

	// Calcular subtotal automáticamente desde precio y descuento
	descuento := 0.0
	if in.Descuento != nil {
		descuento = *in.Descuento
	}
	subtotal := in.Precio - descuento

	var igv, total float64
	switch in.IDTipoComprobante {
	case 2:
		igv = subtotal * 0.18
		total = subtotal + igv
	default:
		igv = 0
		total = subtotal
	}

	qrData := fmt.Sprintf("VENTA|%s-%06d|PRECIO:%.2f|TOTAL:%.2f", serie, correlativo, in.Precio, total)
	qr, errQR := pkg.GenerateQRCode(qrData, 256)
	if errQR != nil {
		return nil, errors.New(util.MSG_VENTA_QR_ERROR)
	}

	venta := &domain.Venta{
		IDUsuario:         in.IDUsuario,
		IDTipoComprobante: in.IDTipoComprobante,
		IDProgramacion:    in.IDProgramacion,
		IDPasajero:        in.IDPasajero,
		IDAsiento:         in.IDAsiento,
		IDTramo:           in.IDTramo,
		Precio:            in.Precio,
		Descuento:         in.Descuento,
		Serie:             serie,
		Correlativo:       correlativo,
		Nota:              in.Nota,
		Observaciones:     in.Observaciones,
		Subtotal:          subtotal,
		IGV:               igv,
		Total:             total,
		QRCode:            qr,
	}

	if err := s.repo.Create(venta); err != nil {
		return nil, util.ParseDBError(err, util.ERR_CODE_CREATE, util.MSG_VENTA_CREATE_ERROR)
	}

	// Marcar asiento como ocupado en el tramo
	if err := s.asientoTramoSvc.MarkAsOccupied(in.IDAsiento, in.IDTramo, &venta.IDVenta); err != nil {
		// Si falla, eliminar la venta creada
		_ = s.repo.Delete(venta.IDVenta)
		return nil, pkg.NewAppError(500, util.ERR_CODE_CREATE, util.MSG_VENTA_CREATE_ERROR).WithCause(err)
	}

	return toVentaOutput(venta), nil
}

func (s *VentaServiceImpl) Update(id int64, in ventaInput.VentaUpdateInput) (*ventaInput.VentaOutput, error) {
	venta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_VENTA_NOT_FOUND)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_UPDATE, util.MSG_VENTA_UPDATE_ERROR).WithCause(err)
	}
	venta.Nota = in.Nota
	venta.Observaciones = in.Observaciones
	pkg.TrimSpacesOnStruct(venta)

	if err := s.repo.Update(venta); err != nil {
		return nil, util.ParseDBError(err, util.ERR_CODE_UPDATE, util.MSG_VENTA_UPDATE_ERROR)
	}
	return toVentaOutput(venta), nil
}

func (s *VentaServiceImpl) Delete(id int64) error {
	if id <= 0 {
		return pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_VENTA_NOT_FOUND)
	}
	venta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_VENTA_NOT_FOUND)
		}
		return pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_DELETE, util.MSG_VENTA_DELETE_ERROR).WithCause(err)
	}

	// Liberar asiento en el tramo
	if err := s.asientoTramoSvc.MarkAsAvailable(venta.IDAsiento, venta.IDTramo); err != nil {
		return pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_DELETE, util.MSG_VENTA_DELETE_ERROR).WithCause(err)
	}

	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_VENTA_NOT_FOUND)
		}
		return util.ParseDBError(err, util.ERR_CODE_DELETE, util.MSG_VENTA_DELETE_ERROR)
	}
	return nil
}

func (s *VentaServiceImpl) GetByID(id int64) (*ventaInput.VentaOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_VENTA_NOT_FOUND)
	}
	venta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_VENTA_NOT_FOUND)
		}
		return nil, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_NOT_FOUND, util.MSG_VENTA_NOT_FOUND).WithCause(err)
	}
	return toVentaOutput(venta), nil
}

func (s *VentaServiceImpl) List(page, size int) ([]ventaInput.VentaOutput, int, error) {
	offset, limit, _ := pkg.Paginate(page, size, 0)
	ventas, total, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, 0, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_LIST, util.MSG_VENTA_LIST_ERROR).WithCause(err)
	}
	outs := make([]ventaInput.VentaOutput, 0, len(ventas))
	for _, v := range ventas {
		vCopy := v
		outs = append(outs, *toVentaOutput(&vCopy))
	}
	return outs, total, nil
}
