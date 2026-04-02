package service

import (
	"errors"
	"fmt"
	"sistema_venta_pasajes/internal/venta/domain"
	"sistema_venta_pasajes/internal/venta/input"
	"sistema_venta_pasajes/internal/venta/repository"
	"sistema_venta_pasajes/internal/venta/util"
	"sistema_venta_pasajes/pkg"
)

type VentaService interface {
	Create(in input.VentaCreateInput) (*input.VentaOutput, error)
	Update(id int64, in input.VentaUpdateInput) (*input.VentaOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*input.VentaOutput, error)
	List(page, size int) ([]input.VentaOutput, int, error)
}

type VentaServiceImpl struct {
	repo repository.VentaRepository
}

func NewVentaService(repo repository.VentaRepository) VentaService {
	return &VentaServiceImpl{repo: repo}
}

const layoutDateTime = "2006-01-02 15:04:05"

func toVentaOutput(v *domain.Venta) *input.VentaOutput {
	return &input.VentaOutput{
		IDVenta:           v.IDVenta,
		IDUsuario:         v.IDUsuario,
		IDTipoComprobante: v.IDTipoComprobante,
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

func (s *VentaServiceImpl) Create(in input.VentaCreateInput) (*input.VentaOutput, error) {
	if !util.ValidarVentaInput(in.IDUsuario, in.IDTipoComprobante, in.Subtotal) {
		if in.IDUsuario <= 0 {
			return nil, errors.New(util.MSG_VENTA_USUARIO_REQUIRED)
		}
		if in.IDTipoComprobante <= 0 {
			return nil, errors.New(util.MSG_VENTA_COMPROBANTE_REQUIRED)
		}
		return nil, errors.New(util.MSG_VENTA_SUBTOTAL_REQUIRED)
	}

	serie, err := util.SerieFromTipoComprobante(in.IDTipoComprobante)
	if err != nil {
		return nil, err
	}

	correlativo, err := s.repo.NextCorrelativo(serie)
	if err != nil {
		return nil, errors.New(util.MSG_VENTA_CORRELATIVO_ERROR)
	}

	var igv, total float64
	switch in.IDTipoComprobante {
	case 2: // FACTURA → aplica 18% IGV
		igv = in.Subtotal * 0.18
		total = in.Subtotal + igv
	default: // BOLETA o TICKET → sin IGV
		igv = 0
		total = in.Subtotal
	}

	qrData := fmt.Sprintf("VENTA|%s-%06d|SUBTOTAL:%.2f|TOTAL:%.2f", serie, correlativo, in.Subtotal, total)
	qr, errQR := pkg.GenerateQRCode(qrData, 256)
	if errQR != nil {
		return nil, errors.New(util.MSG_VENTA_QR_ERROR)
	}

	venta := &domain.Venta{
		IDUsuario:         in.IDUsuario,
		IDTipoComprobante: in.IDTipoComprobante,
		Serie:             serie,
		Correlativo:       correlativo,
		Nota:              in.Nota,
		Observaciones:     in.Observaciones,
		Subtotal:          in.Subtotal,
		IGV:               igv,
		Total:             total,
		QRCode:            qr,
	}
	pkg.TrimSpacesOnStruct(venta)

	if err := s.repo.Create(venta); err != nil {
		return nil, err
	}
	return toVentaOutput(venta), nil
}

func (s *VentaServiceImpl) Update(id int64, in input.VentaUpdateInput) (*input.VentaOutput, error) {
	venta, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New(util.MSG_VENTA_NOT_FOUND)
	}
	venta.Nota = in.Nota
	venta.Observaciones = in.Observaciones
	pkg.TrimSpacesOnStruct(venta)

	if err := s.repo.Update(venta); err != nil {
		return nil, err
	}
	return toVentaOutput(venta), nil
}

func (s *VentaServiceImpl) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *VentaServiceImpl) GetByID(id int64) (*input.VentaOutput, error) {
	venta, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New(util.MSG_VENTA_NOT_FOUND)
	}
	return toVentaOutput(venta), nil
}

func (s *VentaServiceImpl) List(page, size int) ([]input.VentaOutput, int, error) {
	offset, limit, _ := pkg.Paginate(page, size, 0)
	ventas, total, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, 0, err
	}
	var outs []input.VentaOutput
	for _, v := range ventas {
		vCopy := v
		outs = append(outs, *toVentaOutput(&vCopy))
	}
	return outs, total, nil
}
