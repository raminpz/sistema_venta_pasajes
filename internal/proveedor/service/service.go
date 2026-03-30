package service

import (
	"context"
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"sistema_venta_pasajes/pkg"
	"strings"

	"sistema_venta_pasajes/internal/proveedor/domain"
	providerinput "sistema_venta_pasajes/internal/proveedor/input"
	"sistema_venta_pasajes/internal/proveedor/repository"
	"sistema_venta_pasajes/internal/proveedor/util"
)

var rucRegex = regexp.MustCompile(`^\d{11}$`)

type (
	Service interface {
		List(ctx context.Context) ([]domain.ProveedorSistema, error)
		GetByID(ctx context.Context, id int64) (*domain.ProveedorSistema, error)
		Create(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error)
		Update(ctx context.Context, id int64, input providerinput.UpdateInput) (*domain.ProveedorSistema, error)
		Delete(ctx context.Context, id int64) error
	}

	service struct {
		repo repository.Repository
	}
)

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context) ([]domain.ProveedorSistema, error) {
	proveedores, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	if proveedores == nil {
		return make([]domain.ProveedorSistema, 0), nil
	}
	return proveedores, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*domain.ProveedorSistema, error) {
	if id <= 0 {
		return nil, pkg.BadRequest("invalid_provider_id", util.MSG_PROVIDER_ID_GT_ZERO)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) Create(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error) {
	pkg.TrimSpacesOnStruct(&input)
	input.Email = strings.ToLower(input.Email)
	if err := validateInput(input.RUC, input.RazonSocial, input.NombreComercial, input.Direccion, input.Telefono, input.Email, input.Web); err != nil {
		return nil, err
	}

	proveedor, err := s.repo.Create(ctx, input)
	if err != nil {
		errApp := pkg.AsAppError(err)
		if errApp != nil && errApp.Code == "duplicate_resource" {
			if errApp.Err != nil {
				errStr := errApp.Err.Error()
				if strings.Contains(errStr, "ruc") || strings.Contains(errStr, "RUC") {
					return nil, pkg.Conflict("duplicate_resource", util.ERR_RUC_DUPLICADO)
				}
				if strings.Contains(errStr, "email") || strings.Contains(errStr, "EMAIL") {
					return nil, pkg.Conflict("duplicate_resource", util.ERR_EMAIL_DUPLICADO)
				}
			}
			return nil, pkg.Conflict("duplicate_resource", errApp.Message)
		}
		return nil, pkg.Internal("Error al crear proveedor del sistema", err)
	}
	return proveedor, nil
}

func (s *service) Update(ctx context.Context, id int64, input providerinput.UpdateInput) (*domain.ProveedorSistema, error) {
	if id <= 0 {
		return nil, pkg.BadRequest("invalid_provider_id", util.MSG_PROVIDER_ID_GT_ZERO)
	}

	pkg.TrimSpacesOnStruct(&input)
	input.Email = strings.ToLower(input.Email)
	if err := validateInput(input.RUC, input.RazonSocial, input.NombreComercial, input.Direccion, input.Telefono, input.Email, input.Web); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, input)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return pkg.BadRequest("invalid_provider_id", util.MSG_PROVIDER_ID_GT_ZERO)
	}

	return s.repo.Delete(ctx, id)
}

func validateInput(ruc, razonSocial, nombreComercial, direccion, telefono, email, web string) error {
	details := map[string]string{}

	if !rucRegex.MatchString(ruc) {
		details["ruc"] = util.MSG_RUC_FORMAT
	}
	if razonSocial == "" {
		details["razon_social"] = util.MSG_RAZON_SOCIAL_REQUIRED
	}
	if len(razonSocial) > 150 {
		details["razon_social"] = util.MSG_RAZON_SOCIAL_LENGTH
	}
	if len(nombreComercial) > 150 {
		details["nombre_comercial"] = util.MSG_NOMBRE_COMERCIAL_LENGTH
	}
	if len(direccion) > 200 {
		details["direccion"] = util.MSG_DIRECCION_LENGTH
	}
	if len(telefono) > 20 {
		details["telefono"] = util.MSG_TELEFONO_LENGTH
	}
	if email != "" {
		if len(email) > 100 {
			details["email"] = util.MSG_EMAIL_LENGTH
		} else if _, err := mail.ParseAddress(email); err != nil {
			details["email"] = util.MSG_EMAIL_FORMAT
		}
	}
	if web != "" {
		if len(web) > 150 {
			details["web"] = util.MSG_WEB_LENGTH
		} else if err := validateWebURL(web); err != nil {
			details["web"] = err.Error()
		}
	}

	if len(details) > 0 {
		return pkg.Validation(util.MSG_VALIDATION, details)
	}

	return nil
}

func validateWebURL(rawURL string) error {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf(util.MSG_WEB_URL)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf(util.MSG_WEB_SCHEMA_HOST)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf(util.MSG_WEB_HTTP_HTTPS)
	}
	return nil
}
