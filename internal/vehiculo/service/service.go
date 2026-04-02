package service

import (
	"sistema_venta_pasajes/internal/vehiculo/domain"
	"sistema_venta_pasajes/internal/vehiculo/input"
	"sistema_venta_pasajes/internal/vehiculo/repository"
	"sistema_venta_pasajes/internal/vehiculo/util"
	"sistema_venta_pasajes/pkg"
	"strings"
	"time"
)

type VehiculoService interface {
	Create(input.CreateVehiculoInput) (*input.VehiculoOutput, error)
	Update(input.UpdateVehiculoInput) (*input.VehiculoOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*input.VehiculoOutput, error)
	List(page, size int) ([]input.VehiculoOutput, int, error)
}

type vehiculoService struct {
	repo repository.VehiculoRepository
}

func NewVehiculoService(repo repository.VehiculoRepository) VehiculoService {
	return &vehiculoService{repo: repo}
}

func (s *vehiculoService) Create(in input.CreateVehiculoInput) (*input.VehiculoOutput, error) {
	if !util.ValidarPlaca(in.NroPlaca) {
		return nil, pkg.BadRequest("invalid_plate", util.ERR_INVALID_ID)
	}
	in.NroPlaca = strings.ToUpper(in.NroPlaca)
	in.Marca = pkg.CapitalizeWords(in.Marca)
	in.Modelo = pkg.CapitalizeWords(in.Modelo)
	pkg.TrimSpacesOnStruct(&in)
	var fechaVencSoat *time.Time
	if in.FechaVencSoat != nil && !in.FechaVencSoat.Time.IsZero() {
		t := in.FechaVencSoat.Time
		fechaVencSoat = &t
	}
	var fechaVencRevisionTec *time.Time
	if in.FechaVencRevisionTec != nil && !in.FechaVencRevisionTec.Time.IsZero() {
		t := in.FechaVencRevisionTec.Time
		fechaVencRevisionTec = &t
	}
	vehiculo := &domain.Vehiculo{
		IDTipoVehiculo:       in.IDTipoVehiculo,
		NroPlaca:             in.NroPlaca,
		Marca:                in.Marca,
		Modelo:               in.Modelo,
		AnioFabricacion:      in.AnioFabricacion,
		NumeroChasis:         in.NumeroChasis,
		Capacidad:            in.Capacidad,
		NroSoat:              in.NroSoat,
		FechaVencSoat:        fechaVencSoat,
		NroRevisionTecnica:   in.NroRevisionTecnica,
		FechaVencRevisionTec: fechaVencRevisionTec,
		Estado:               in.Estado,
	}
	// Validar duplicidad específica antes de crear
	existsPlaca, err := s.repo.ExistsByPlaca(in.NroPlaca)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	if existsPlaca {
		return nil, pkg.BadRequest("duplicate_placa", util.ERR_DUPLICATE_PLATE)
	}

	existsChasis, err := s.repo.ExistsByChasis(in.NumeroChasis)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	if existsChasis {
		return nil, pkg.BadRequest("duplicate_chasis", util.ERR_DUPLICATE_CHASSIS)
	}

	existsSoat, err := s.repo.ExistsBySoat(in.NroSoat)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	if existsSoat {
		return nil, pkg.BadRequest("duplicate_soat", util.ERR_DUPLICATE_SOAT)
	}

	err = s.repo.Create(vehiculo)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	return mapVehiculoOutput(vehiculo), nil
}

func (s *vehiculoService) Update(in input.UpdateVehiculoInput) (*input.VehiculoOutput, error) {
	vehiculo, err := s.repo.GetByID(in.IDVehiculo)
	if err != nil {
		return nil, pkg.NotFound("vehiculo_not_found", util.ERR_NOT_FOUND)
	}
	in.NroPlaca = strings.ToUpper(in.NroPlaca)
	in.Marca = pkg.CapitalizeWords(in.Marca)
	in.Modelo = pkg.CapitalizeWords(in.Modelo)
	pkg.TrimSpacesOnStruct(&in)
	vehiculo.IDTipoVehiculo = in.IDTipoVehiculo
	vehiculo.NroPlaca = in.NroPlaca
	vehiculo.Marca = in.Marca
	vehiculo.Modelo = in.Modelo
	vehiculo.AnioFabricacion = in.AnioFabricacion
	vehiculo.NumeroChasis = in.NumeroChasis
	vehiculo.Capacidad = in.Capacidad
	vehiculo.NroSoat = in.NroSoat
	if in.FechaVencSoat != nil {
		t := in.FechaVencSoat.Time
		vehiculo.FechaVencSoat = &t
	} else {
		vehiculo.FechaVencSoat = nil
	}
	vehiculo.NroRevisionTecnica = in.NroRevisionTecnica
	if in.FechaVencRevisionTec != nil {
		t := in.FechaVencRevisionTec.Time
		vehiculo.FechaVencRevisionTec = &t
	} else {
		vehiculo.FechaVencRevisionTec = nil
	}
	vehiculo.Estado = in.Estado
	err = s.repo.Update(vehiculo)
	if err != nil {
		return nil, pkg.Internal(err.Error())
	}
	return mapVehiculoOutput(vehiculo), nil
}

func (s *vehiculoService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *vehiculoService) GetByID(id int64) (*input.VehiculoOutput, error) {
	vehiculo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, pkg.NotFound("vehiculo_not_found", util.ERR_NOT_FOUND)
	}
	return mapVehiculoOutput(vehiculo), nil
}

func (s *vehiculoService) List(page, size int) ([]input.VehiculoOutput, int, error) {
	offset, limit, _ := pkg.Paginate(page, size, 0)
	vehiculos, total, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, 0, pkg.Internal(err.Error())
	}
	var out []input.VehiculoOutput
	for _, v := range vehiculos {
		out = append(out, *mapVehiculoOutput(&v))
	}
	return out, total, nil
}

func mapVehiculoOutput(v *domain.Vehiculo) *input.VehiculoOutput {
	var fechaVencSoat *input.DateOnly
	if v.FechaVencSoat != nil {
		fechaVencSoat = &input.DateOnly{Time: *v.FechaVencSoat}
	}
	var fechaVencRevisionTec *input.DateOnly
	if v.FechaVencRevisionTec != nil {
		fechaVencRevisionTec = &input.DateOnly{Time: *v.FechaVencRevisionTec}
	}
	return &input.VehiculoOutput{
		IDVehiculo:           v.IDVehiculo,
		IDTipoVehiculo:       v.IDTipoVehiculo,
		NroPlaca:             v.NroPlaca,
		Marca:                v.Marca,
		Modelo:               v.Modelo,
		AnioFabricacion:      v.AnioFabricacion,
		NumeroChasis:         v.NumeroChasis,
		Capacidad:            v.Capacidad,
		NroSoat:              v.NroSoat,
		FechaVencSoat:        fechaVencSoat,
		NroRevisionTecnica:   v.NroRevisionTecnica,
		FechaVencRevisionTec: fechaVencRevisionTec,
		Estado:               v.Estado,
	}
}
