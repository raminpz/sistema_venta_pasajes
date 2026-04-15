package service

import (
	"context"
	"errors"
	"sistema_venta_pasajes/internal/usuario/domain"
	"sistema_venta_pasajes/internal/usuario/input"
	"sistema_venta_pasajes/internal/usuario/repository"
	"sistema_venta_pasajes/internal/usuario/util"
	"sistema_venta_pasajes/pkg"

	"golang.org/x/crypto/bcrypt"
)

type UsuarioService interface {
	Create(input.UsuarioCreateInput) (*input.UsuarioOutput, error)
	Update(id int, in input.UsuarioUpdateInput) (*input.UsuarioOutput, error)
	Delete(id int) error
	GetByID(id int) (*input.UsuarioOutput, error)
	List(page, size int) ([]input.UsuarioOutput, int, error)
}

type usuarioService struct {
	repo repository.UsuarioRepository
}

func NewUsuarioService(repo repository.UsuarioRepository) UsuarioService {
	return &usuarioService{repo: repo}
}

func (s *usuarioService) Create(in input.UsuarioCreateInput) (*input.UsuarioOutput, error) {
	ctx := context.Background()
	if !util.ValidarCamposObligatorios(in.Nombre, in.Apellidos, in.Email, in.Password, in.Telefono) {
		return nil, errors.New(util.MSG_INVALID_DATA)
	}
	if !util.ValidarEmail(in.Email) {
		return nil, errors.New(util.MSG_INVALID_DATA)
	}
	if !util.ValidarDNI(in.DNI) {
		return nil, errors.New(util.MSG_INVALID_DATA)
	}
	// Validar unicidad de EMAIL
	existingEmail, _ := s.repo.GetByEmail(ctx, in.Email)
	if existingEmail != nil {
		return nil, errors.New(util.MSG_EMAIL_DUPLICATE)
	}
	// Validar unicidad de DNI
	existingDNI, _ := s.repo.GetByDNI(ctx, in.DNI)
	if existingDNI != nil {
		return nil, errors.New(util.MSG_DNI_DUPLICATE)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// Capitalizar nombre y apellidos
	in.Nombre = pkg.CapitalizeWords(in.Nombre)
	in.Apellidos = pkg.CapitalizeWords(in.Apellidos)
	usuario := &domain.Usuario{
		IDRol:     in.IDRol,
		Nombre:    in.Nombre,
		Apellidos: in.Apellidos,
		DNI:       in.DNI,
		Email:     in.Email,
		Password:  string(hashed),
		Telefono:  in.Telefono,
		Estado:    "ACTIVO",
	}
	err = s.repo.Create(ctx, usuario)
	if err != nil {
		return nil, err
	}
	return toUsuarioOutput(usuario), nil
}

func (s *usuarioService) Update(id int, in input.UsuarioUpdateInput) (*input.UsuarioOutput, error) {
	ctx := context.Background()
	usuario, err := s.repo.GetByID(ctx, id)
	if err != nil || usuario == nil {
		return nil, errors.New(util.MSG_USER_NOT_FOUND)
	}
	if in.Email != usuario.Email {
		existingEmail, _ := s.repo.GetByEmail(ctx, in.Email)
		if existingEmail != nil && existingEmail.IDUsuario != usuario.IDUsuario {
			return nil, errors.New(util.MSG_EMAIL_DUPLICATE)
		}
	}
	usuario.Nombre = in.Nombre
	usuario.Apellidos = in.Apellidos
	usuario.Email = in.Email
	if in.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		usuario.Password = string(hashed)
	}
	usuario.Telefono = in.Telefono
	usuario.Estado = in.Estado
	if in.Nombre != "" {
		usuario.Nombre = pkg.CapitalizeWords(in.Nombre)
	}
	if in.Apellidos != "" {
		usuario.Apellidos = pkg.CapitalizeWords(in.Apellidos)
	}
	err = s.repo.Update(ctx, id, usuario)
	if err != nil {
		return nil, err
	}
	return toUsuarioOutput(usuario), nil
}

func (s *usuarioService) Delete(id int) error {
	ctx := context.Background()
	usuario, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New(util.MSG_USER_NOT_FOUND)
	}
	if usuario == nil {
		return errors.New(util.MSG_USER_NOT_FOUND)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return pkg.Internal(util.MSG_DELETE_ERROR, err)
	}
	return nil
}

func (s *usuarioService) GetByID(id int) (*input.UsuarioOutput, error) {
	ctx := context.Background()
	usuario, err := s.repo.GetByID(ctx, id)
	if err != nil || usuario == nil {
		return nil, errors.New(util.MSG_USER_NOT_FOUND)
	}
	return toUsuarioOutput(usuario), nil
}

func (s *usuarioService) List(page, size int) ([]input.UsuarioOutput, int, error) {
	ctx := context.Background()
	offset, limit, _ := pkg.Paginate(page, size, 0)
	usuarios, total, err := s.repo.List(ctx, map[string]interface{}{}, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	var outs []input.UsuarioOutput
	for _, u := range usuarios {
		outs = append(outs, *toUsuarioOutput(u))
	}
	return outs, total, nil
}

func toUsuarioOutput(u *domain.Usuario) *input.UsuarioOutput {
	return &input.UsuarioOutput{
		IDUsuario: u.IDUsuario,
		IDRol:     u.IDRol,
		Nombre:    u.Nombre,
		Apellidos: u.Apellidos,
		DNI:       u.DNI,
		Email:     u.Email,
		Telefono:  u.Telefono,
		Estado:    u.Estado,
	}
}
