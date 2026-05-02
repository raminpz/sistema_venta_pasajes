package util

import (
	"regexp"
	"sistema_venta_pasajes/internal/vehiculo/input"
	"sistema_venta_pasajes/pkg"
	"strings"
)

func ValidarPlaca(placa string) bool {
	reg := regexp.MustCompile(`^[A-Z0-9-]{6,10}$`)
	return reg.MatchString(placa)
}

func ValidarCamposCreate(in input.CreateVehiculoInput) error {
	if in.IDTipoVehiculo <= 0 {
		return pkg.BadRequest("required_tipo_vehiculo", ERR_REQUIRED_TIPO_VEHICULO)
	}
	if in.NroPlaca == "" {
		return pkg.BadRequest("required_placa", ERR_REQUIRED_PLACA)
	}
	if in.Marca == "" {
		return pkg.BadRequest("required_marca", ERR_REQUIRED_MARCA)
	}
	if in.Modelo == "" {
		return pkg.BadRequest("required_modelo", ERR_REQUIRED_MODELO)
	}
	if in.AnioFabricacion <= 0 {
		return pkg.BadRequest("required_anio", ERR_REQUIRED_ANIO)
	}
	if in.NumeroChasis == "" {
		return pkg.BadRequest("required_chasis", ERR_REQUIRED_CHASIS)
	}
	if in.Capacidad <= 0 {
		return pkg.BadRequest("required_capacidad", ERR_REQUIRED_CAPACIDAD)
	}
	if in.NroSoat == "" {
		return pkg.BadRequest("required_soat", ERR_REQUIRED_SOAT)
	}
	if in.FechaVencSoat == nil || in.FechaVencSoat.Time.IsZero() {
		return pkg.BadRequest("required_fecha_venc_soat", ERR_REQUIRED_FECHA_VENC_SOAT)
	}
	if in.NroRevisionTecnica == "" {
		return pkg.BadRequest("required_revision_tecnica", ERR_REQUIRED_REVISION_TECNICA)
	}
	if in.FechaVencRevisionTec == nil || in.FechaVencRevisionTec.Time.IsZero() {
		return pkg.BadRequest("required_fecha_venc_revision_tecnica", ERR_REQUIRED_FECHA_VENC_REV)
	}
	if in.Estado == "" {
		return pkg.BadRequest("required_estado", ERR_REQUIRED_ESTADO)
	}
	if in.Estado != "ACTIVO" && in.Estado != "INACTIVO" {
		return pkg.BadRequest("invalid_estado", ERR_INVALID_ESTADO)
	}
	return nil
}

func ValidarCamposUpdate(in input.UpdateVehiculoInput) error {
	if in.IDTipoVehiculo == nil && in.NroPlaca == nil && in.Marca == nil && in.Modelo == nil &&
		in.AnioFabricacion == nil && in.NumeroChasis == nil && in.Capacidad == nil && in.NroSoat == nil &&
		in.FechaVencSoat == nil && in.NroRevisionTecnica == nil && in.FechaVencRevisionTec == nil && in.Estado == nil {
		return pkg.BadRequest("empty_update", ERR_EMPTY_UPDATE)
	}

	if in.IDTipoVehiculo != nil && *in.IDTipoVehiculo <= 0 {
		return pkg.BadRequest("required_tipo_vehiculo", ERR_REQUIRED_TIPO_VEHICULO)
	}
	if in.NroPlaca != nil && strings.TrimSpace(*in.NroPlaca) == "" {
		return pkg.BadRequest("required_placa", ERR_REQUIRED_PLACA)
	}
	if in.Marca != nil && strings.TrimSpace(*in.Marca) == "" {
		return pkg.BadRequest("required_marca", ERR_REQUIRED_MARCA)
	}
	if in.Modelo != nil && strings.TrimSpace(*in.Modelo) == "" {
		return pkg.BadRequest("required_modelo", ERR_REQUIRED_MODELO)
	}
	if in.AnioFabricacion != nil && *in.AnioFabricacion <= 0 {
		return pkg.BadRequest("required_anio", ERR_REQUIRED_ANIO)
	}
	if in.NumeroChasis != nil && strings.TrimSpace(*in.NumeroChasis) == "" {
		return pkg.BadRequest("required_chasis", ERR_REQUIRED_CHASIS)
	}
	if in.Capacidad != nil && *in.Capacidad <= 0 {
		return pkg.BadRequest("required_capacidad", ERR_REQUIRED_CAPACIDAD)
	}
	if in.NroSoat != nil && strings.TrimSpace(*in.NroSoat) == "" {
		return pkg.BadRequest("required_soat", ERR_REQUIRED_SOAT)
	}
	if in.FechaVencSoat != nil && in.FechaVencSoat.Time.IsZero() {
		return pkg.BadRequest("required_fecha_venc_soat", ERR_REQUIRED_FECHA_VENC_SOAT)
	}
	if in.NroRevisionTecnica != nil && strings.TrimSpace(*in.NroRevisionTecnica) == "" {
		return pkg.BadRequest("required_revision_tecnica", ERR_REQUIRED_REVISION_TECNICA)
	}
	if in.FechaVencRevisionTec != nil && in.FechaVencRevisionTec.Time.IsZero() {
		return pkg.BadRequest("required_fecha_venc_revision_tecnica", ERR_REQUIRED_FECHA_VENC_REV)
	}
	if in.Estado != nil {
		estado := strings.ToUpper(strings.TrimSpace(*in.Estado))
		if estado == "" {
			return pkg.BadRequest("required_estado", ERR_REQUIRED_ESTADO)
		}
		if estado != "ACTIVO" && estado != "INACTIVO" {
			return pkg.BadRequest("invalid_estado", ERR_INVALID_ESTADO)
		}
	}

	return nil
}
