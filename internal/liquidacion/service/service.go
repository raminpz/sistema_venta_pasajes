package service

import (
	"sistema_venta_pasajes/internal/liquidacion/domain"
	"sistema_venta_pasajes/internal/liquidacion/input"
	"sistema_venta_pasajes/internal/liquidacion/repository"
	"sistema_venta_pasajes/internal/liquidacion/util"
	"sistema_venta_pasajes/pkg"
	"time"
)

// LiquidacionService define los casos de uso del módulo.
type LiquidacionService interface {
	Generar(in input.GenerarLiquidacionInput) (*input.LiquidacionOutput, error)
	ActualizarEstado(id int64, in input.ActualizarEstadoInput) (*input.LiquidacionOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*input.LiquidacionOutput, error)
	List(page, size int) ([]input.LiquidacionOutput, int, error)
	ObtenerResumenCaja(idProgramacion int64) (*input.ResumenCajaOutput, error)
}
type liquidacionService struct {
	repo repository.LiquidacionRepository
}

// NewLiquidacionService crea el servicio de liquidación.
func NewLiquidacionService(repo repository.LiquidacionRepository) LiquidacionService {
	return &liquidacionService{repo: repo}
}

// Generar calcula los totales de la programación y persiste la liquidación.
func (s *liquidacionService) Generar(in input.GenerarLiquidacionInput) (*input.LiquidacionOutput, error) {
	if in.IDProgramacion <= 0 {
		return nil, pkg.BadRequest("programacion_requerida", util.ERR_PROGRAMACION_REQUERIDA)
	}
	exists, err := s.repo.ExistsByProgramacion(in.IDProgramacion)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	if exists {
		return nil, pkg.Conflict("liquidacion_existente", util.ERR_LIQUIDACION_EXISTENTE)
	}
	idConductor, err := s.repo.GetConductorByProgramacion(in.IDProgramacion)
	if err != nil {
		return nil, pkg.NotFound("programacion_not_found", util.ERR_PROGRAMACION_NOT_FOUND)
	}
	totalPasajes, cantPasajes, err := s.repo.SumarVentas(in.IDProgramacion)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	totalEncomiendas, cantEncomiendas, err := s.repo.SumarEncomiendas(in.IDProgramacion)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	pkg.TrimSpacesOnStruct(&in)
	liq := &domain.LiquidacionViaje{
		IDProgramacion:   in.IDProgramacion,
		IDConductor:      idConductor,
		TotalPasajes:     totalPasajes,
		TotalEncomiendas: totalEncomiendas,
		TotalCaja:        totalPasajes + totalEncomiendas,
		Estado:           "PENDIENTE",
		Observaciones:    in.Observaciones,
	}
	if err := s.repo.Create(liq); err != nil {
		return nil, pkg.Internal(err.Error())
	}
	out := mapLiquidacionOutput(liq)
	out.CantidadPasajes = cantPasajes
	out.CantidadEncomiendas = cantEncomiendas
	return out, nil
}

// ActualizarEstado cambia el estado de la liquidación (PENDIENTE → ENTREGADO).
func (s *liquidacionService) ActualizarEstado(id int64, in input.ActualizarEstadoInput) (*input.LiquidacionOutput, error) {
	if !util.EsEstadoValido(in.Estado) {
		return nil, pkg.BadRequest("estado_invalido", util.ERR_ESTADO_INVALIDO)
	}
	liq, err := s.repo.GetByID(id)
	if err != nil {
		return nil, pkg.NotFound("liquidacion_not_found", util.ERR_NOT_FOUND)
	}
	liq.Estado = in.Estado
	pkg.TrimSpacesOnStruct(&in)
	if in.Observaciones != "" {
		liq.Observaciones = in.Observaciones
	}
	if in.Estado == "ENTREGADO" {
		now := time.Now()
		liq.FechaLiquidacion = &now
	}
	if err := s.repo.Update(liq); err != nil {
		return nil, pkg.Internal(err.Error())
	}
	return mapLiquidacionOutput(liq), nil
}

// Delete elimina la liquidación por ID.
func (s *liquidacionService) Delete(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return pkg.NotFound("liquidacion_not_found", util.ERR_NOT_FOUND)
	}
	return s.repo.Delete(id)
}

// GetByID obtiene una liquidación por su ID.
func (s *liquidacionService) GetByID(id int64) (*input.LiquidacionOutput, error) {
	liq, err := s.repo.GetByID(id)
	if err != nil {
		return nil, pkg.NotFound("liquidacion_not_found", util.ERR_NOT_FOUND)
	}
	return mapLiquidacionOutput(liq), nil
}

// List retorna liquidaciones paginadas.
func (s *liquidacionService) List(page, size int) ([]input.LiquidacionOutput, int, error) {
	offset, limit, _ := pkg.Paginate(page, size, 0)
	liqs, total, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, 0, pkg.Internal(err.Error())
	}
	var out []input.LiquidacionOutput
	for _, l := range liqs {
		out = append(out, *mapLiquidacionOutput(&l))
	}
	if out == nil {
		out = []input.LiquidacionOutput{}
	}
	return out, total, nil
}

// ObtenerResumenCaja calcula el total de caja sin persistir el resultado.
func (s *liquidacionService) ObtenerResumenCaja(idProgramacion int64) (*input.ResumenCajaOutput, error) {
	if idProgramacion <= 0 {
		return nil, pkg.BadRequest("id_invalido", util.ERR_INVALID_ID)
	}
	_, err := s.repo.GetConductorByProgramacion(idProgramacion)
	if err != nil {
		return nil, pkg.NotFound("programacion_not_found", util.ERR_PROGRAMACION_NOT_FOUND)
	}
	totalPasajes, cantPasajes, err := s.repo.SumarVentas(idProgramacion)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	totalEncomiendas, cantEncomiendas, err := s.repo.SumarEncomiendas(idProgramacion)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	return &input.ResumenCajaOutput{
		IDProgramacion:      idProgramacion,
		TotalPasajes:        totalPasajes,
		TotalEncomiendas:    totalEncomiendas,
		TotalCaja:           totalPasajes + totalEncomiendas,
		CantidadPasajes:     cantPasajes,
		CantidadEncomiendas: cantEncomiendas,
	}, nil
}
func mapLiquidacionOutput(l *domain.LiquidacionViaje) *input.LiquidacionOutput {
	return &input.LiquidacionOutput{
		IDLiquidacion:    l.IDLiquidacion,
		IDProgramacion:   l.IDProgramacion,
		IDConductor:      l.IDConductor,
		TotalPasajes:     l.TotalPasajes,
		TotalEncomiendas: l.TotalEncomiendas,
		TotalCaja:        l.TotalCaja,
		Estado:           l.Estado,
		FechaLiquidacion: l.FechaLiquidacion,
		Observaciones:    l.Observaciones,
	}
}
