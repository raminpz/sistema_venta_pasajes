package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sistema_venta_pasajes/internal/control_acceso/domain"
	"sistema_venta_pasajes/internal/control_acceso/input"
	"sistema_venta_pasajes/internal/control_acceso/repository"
	"sistema_venta_pasajes/internal/control_acceso/util"
	"sistema_venta_pasajes/pkg"
)

const dateLayout = "2006-01-02"

type Service interface {
	GetStatus() (*input.ControlAccesoStatusOutput, error)
	GetLatest() (*input.ControlAccesoOutput, error)
	Create(in input.ActivarControlAccesoInput) (*input.ControlAccesoOutput, error)
	Activar(id int64) error
	Bloquear(id int64) error
	Renovar(id int64, in input.RenovarControlAccesoInput) (*input.ControlAccesoOutput, error)
}

type service struct {
	repo repository.ControlAccesoRepository
}

func New(repo repository.ControlAccesoRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetStatus() (*input.ControlAccesoStatusOutput, error) {
	a, err := s.repo.GetLatest()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &input.ControlAccesoStatusOutput{
				EstadoEfectivo: "BLOQUEADO",
				Mensaje:        fmt.Sprintf("No hay un control de acceso activo. Comuníquese con el proveedor al %s.", util.PROVEEDOR_TELEFONO),
				DiasParaVencer: 0,
				EnAlerta:       false,
				EnGracia:       false,
			}, nil
		}
		return nil, pkg.Internal("Error al consultar el estado del sistema.")
	}
	info := computeState(a)
	return &input.ControlAccesoStatusOutput{
		EstadoEfectivo: info.estado,
		Mensaje:        info.mensaje,
		DiasParaVencer: info.diasParaVencer,
		EnAlerta:       info.enAlerta,
		EnGracia:       info.enGracia,
	}, nil
}

func (s *service) GetLatest() (*input.ControlAccesoOutput, error) {
	a, err := s.repo.GetLatest()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.CODE_ACCESO_NO_ENCONTRADO, util.ERR_ACCESO_NO_ENCONTRADO)
		}
		return nil, pkg.Internal("Error al obtener el control de acceso.")
	}
	return toOutput(a), nil
}

func (s *service) Create(in input.ActivarControlAccesoInput) (*input.ControlAccesoOutput, error) {
	if in.FechaActivacion == "" || in.FechaExpiracion == "" {
		return nil, pkg.BadRequest(util.CODE_ACCESO_FECHA_INVALIDA, util.ERR_ACCESO_FECHA_REQUERIDA)
	}
	activacion, err := time.Parse(dateLayout, in.FechaActivacion)
	if err != nil {
		return nil, pkg.BadRequest(util.CODE_ACCESO_FECHA_INVALIDA, util.ERR_ACCESO_FECHA_FORMATO)
	}
	expiracion, err := time.Parse(dateLayout, in.FechaExpiracion)
	if err != nil {
		return nil, pkg.BadRequest(util.CODE_ACCESO_FECHA_INVALIDA, util.ERR_ACCESO_FECHA_FORMATO)
	}
	if expiracion.Before(activacion) {
		return nil, pkg.BadRequest(util.CODE_ACCESO_FECHA_VALID, util.ERR_ACCESO_FECHA_EXPIRACION)
	}
	a := &domain.ControlAcceso{
		FechaActivacion: activacion,
		FechaExpiracion: expiracion,
		Estado:          "OPERATIVO",
	}
	if err := s.repo.Create(a); err != nil {
		return nil, pkg.Internal("Error al crear el control de acceso.")
	}
	return toOutput(a), nil
}

func (s *service) Activar(id int64) error {
	if err := s.repo.SetEstado(id, "OPERATIVO"); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.CODE_ACCESO_NO_ENCONTRADO, util.ERR_ACCESO_NO_ENCONTRADO)
		}
		return pkg.Internal("Error al activar el control de acceso.")
	}
	return nil
}

func (s *service) Bloquear(id int64) error {
	if err := s.repo.SetEstado(id, "BLOQUEADO"); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.CODE_ACCESO_NO_ENCONTRADO, util.ERR_ACCESO_NO_ENCONTRADO)
		}
		return pkg.Internal("Error al bloquear el control de acceso.")
	}
	return nil
}

func (s *service) Renovar(id int64, in input.RenovarControlAccesoInput) (*input.ControlAccesoOutput, error) {
	if in.FechaExpiracion == "" {
		return nil, pkg.BadRequest(util.CODE_ACCESO_FECHA_INVALIDA, util.ERR_ACCESO_EXP_REQUERIDA)
	}
	nuevaFecha, err := time.Parse(dateLayout, in.FechaExpiracion)
	if err != nil {
		return nil, pkg.BadRequest(util.CODE_ACCESO_FECHA_INVALIDA, util.ERR_ACCESO_FECHA_FORMATO)
	}
	if err := s.repo.Renovar(id, nuevaFecha); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.CODE_ACCESO_NO_ENCONTRADO, util.ERR_ACCESO_NO_ENCONTRADO)
		}
		return nil, pkg.Internal("Error al renovar el control de acceso.")
	}
	a, err := s.repo.GetByID(id)
	if err != nil {
		return nil, pkg.Internal("Error al obtener el control de acceso renovado.")
	}
	return toOutput(a), nil
}

// --- helpers internos ---

type stateInfo struct {
	estado         string
	mensaje        string
	diasParaVencer int
	enAlerta       bool
	enGracia       bool
}

// computeState calcula el estado efectivo del control de acceso.
//
// Flujo de estados:
//  1. Proveedor bloqueó manualmente → BLOQUEADO (prioridad máxima)
//  2. Fecha aún vigente, más de 30 días → OPERATIVO normal
//  3. Fecha aún vigente, 30 días o menos → OPERATIVO + en_alerta
//  4. Vencida, dentro de 30 días de gracia → SOLO_LECTURA + en_gracia
//  5. Vencida, gracia agotada → BLOQUEADO automático
func computeState(a *domain.ControlAcceso) stateInfo {
	if a.Estado == "BLOQUEADO" {
		return stateInfo{
			estado:  "BLOQUEADO",
			mensaje: fmt.Sprintf("El sistema está bloqueado. Comuníquese con el proveedor al %s para reactivarlo.", util.PROVEEDOR_TELEFONO),
		}
	}

	now := time.Now()
	exp := a.FechaExpiracion

	if !now.After(exp) {
		dias := int(exp.Sub(now).Hours() / 24)
		if dias <= util.DIAS_ALERTA {
			return stateInfo{
				estado: "OPERATIVO",
				mensaje: fmt.Sprintf(
					"Su fecha de renovación del sistema termina el %s. Por favor comuníquese con el proveedor al celular %s para renovar y no se perjudique. Muchas gracias por su atención.",
					exp.Format("02/01/2006"),
					util.PROVEEDOR_TELEFONO,
				),
				diasParaVencer: dias,
				enAlerta:       true,
			}
		}
		return stateInfo{
			estado:         "OPERATIVO",
			mensaje:        "Sistema operativo.",
			diasParaVencer: dias,
		}
	}

	graceEnd := exp.AddDate(0, 0, util.DIAS_GRACIA)
	if !now.After(graceEnd) {
		diasGracia := int(graceEnd.Sub(now).Hours() / 24)
		return stateInfo{
			estado: "SOLO_LECTURA",
			mensaje: fmt.Sprintf(
				"Su suscripción venció el %s. El sistema está en modo solo lectura por %d días de gracia. Comuníquese con el proveedor al %s para renovar.",
				exp.Format("02/01/2006"),
				diasGracia,
				util.PROVEEDOR_TELEFONO,
			),
			diasParaVencer: diasGracia,
			enGracia:       true,
		}
	}

	return stateInfo{
		estado: "BLOQUEADO",
		mensaje: fmt.Sprintf(
			"El sistema está bloqueado. Su suscripción y el período de gracia han vencido. Comuníquese con el proveedor al %s para renovar.",
			util.PROVEEDOR_TELEFONO,
		),
	}
}

func toOutput(a *domain.ControlAcceso) *input.ControlAccesoOutput {
	info := computeState(a)
	return &input.ControlAccesoOutput{
		IDAcceso:        a.IDAcceso,
		FechaActivacion: a.FechaActivacion.Format(dateLayout),
		FechaExpiracion: a.FechaExpiracion.Format(dateLayout),
		EstadoDB:        a.Estado,
		EstadoEfectivo:  info.estado,
		DiasParaVencer:  info.diasParaVencer,
		EnAlerta:        info.enAlerta,
		EnGracia:        info.enGracia,
	}
}
