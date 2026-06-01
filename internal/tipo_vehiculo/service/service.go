package service

import (
	"errors"
	"sistema_venta_pasajes/internal/tipo_vehiculo/domain"
	"sistema_venta_pasajes/internal/tipo_vehiculo/input"
	"sistema_venta_pasajes/internal/tipo_vehiculo/repository"
	"sistema_venta_pasajes/internal/tipo_vehiculo/util"
	"sistema_venta_pasajes/pkg"

	"gorm.io/gorm"
)

type TipoVehiculoService interface {
	Create(in input.CreateTipoVehiculoInput) (*input.TipoVehiculoOutput, error)
	Update(id int64, in input.UpdateTipoVehiculoInput) (*input.TipoVehiculoOutput, error)
	Delete(id int64) error
	GetByID(id int64) (*input.TipoVehiculoOutput, error)
	List() ([]input.TipoVehiculoOutput, error)
}

type tipoVehiculoService struct {
	repo repository.TipoVehiculoRepository
}

func NewTipoVehiculoService(repo repository.TipoVehiculoRepository) TipoVehiculoService {
	return &tipoVehiculoService{repo: repo}
}

func mapOutput(tv *domain.TipoVehiculo) *input.TipoVehiculoOutput {
	return &input.TipoVehiculoOutput{
		IDTipoVehiculo: tv.IDTipoVehiculo,
		Nombre:         tv.Nombre,
		Descripcion:    tv.Descripcion,
	}
}

func (s *tipoVehiculoService) Create(in input.CreateTipoVehiculoInput) (*input.TipoVehiculoOutput, error) {
	pkg.TrimSpacesOnStruct(&in)
	if err := util.ValidateCreate(in); err != nil {
		return nil, err
	}

	tv := &domain.TipoVehiculo{
		Nombre:      in.Nombre,
		Descripcion: in.Descripcion,
	}

	if err := s.repo.Create(tv); err != nil {
		return nil, pkg.ParseDBError(err, util.ERR_CODE_CREATE, util.ERR_CREATE, nil,
			map[string]string{"NOMBRE": util.ERR_DUPLICATE_NAME})
	}

	return mapOutput(tv), nil
}

func (s *tipoVehiculoService) Update(id int64, in input.UpdateTipoVehiculoInput) (*input.TipoVehiculoOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_INVALIDID, util.ERR_INVALID_ID)
	}

	pkg.TrimSpacesOnStruct(&in)
	if err := util.ValidateUpdate(in); err != nil {
		return nil, err
	}

	tv, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return nil, pkg.Internal(util.ERR_UPDATE, err)
	}

	if in.Nombre != nil {
		tv.Nombre = *in.Nombre
	}
	if in.Descripcion != nil {
		tv.Descripcion = *in.Descripcion
	}

	if err := s.repo.Update(tv); err != nil {
		return nil, pkg.ParseDBError(err, util.ERR_CODE_UPDATE, util.ERR_UPDATE, nil,
			map[string]string{"NOMBRE": util.ERR_DUPLICATE_NAME})
	}
	return mapOutput(tv), nil
}

func (s *tipoVehiculoService) Delete(id int64) error {
	if id <= 0 {
		return pkg.BadRequest(util.ERR_CODE_INVALIDID, util.ERR_INVALID_ID)
	}

	if _, err := s.repo.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return pkg.Internal(util.ERR_DELETE, err)
	}

	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return pkg.ParseDBError(err, util.ERR_CODE_DELETE, util.ERR_DELETE,
			map[string]string{"*": util.ERR_DELETE_CONFLICT}, nil)
	}
	return nil
}

func (s *tipoVehiculoService) GetByID(id int64) (*input.TipoVehiculoOutput, error) {
	if id <= 0 {
		return nil, pkg.BadRequest(util.ERR_CODE_INVALIDID, util.ERR_INVALID_ID)
	}

	tv, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.ERR_NOT_FOUND)
		}
		return nil, pkg.Internal(util.ERR_NOT_FOUND, err)
	}
	return mapOutput(tv), nil
}

func (s *tipoVehiculoService) List() ([]input.TipoVehiculoOutput, error) {
	items, err := s.repo.List()
	if err != nil {
		return nil, pkg.Internal(util.ERR_LIST, err)
	}

	outputs := make([]input.TipoVehiculoOutput, 0, len(items))
	for _, item := range items {
		itemCopy := item
		outputs = append(outputs, *mapOutput(&itemCopy))
	}
	return outputs, nil
}

