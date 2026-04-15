package repository

import (
	"context"
	"errors"
	"sistema_venta_pasajes/internal/usuario/domain"
)

// UsuarioRepositoryMock es un mock para pruebas.
type UsuarioRepositoryMock struct {
	Usuarios map[int]*domain.Usuario
}

func NewUsuarioRepositoryMock() *UsuarioRepositoryMock {
	return &UsuarioRepositoryMock{
		Usuarios: make(map[int]*domain.Usuario),
	}
}

func (m *UsuarioRepositoryMock) Create(ctx context.Context, usuario *domain.Usuario) error {
	m.Usuarios[usuario.IDUsuario] = usuario
	return nil
}

func (m *UsuarioRepositoryMock) Update(ctx context.Context, id int, usuario *domain.Usuario) error {
	if _, ok := m.Usuarios[id]; !ok {
		return ErrNotFound
	}
	m.Usuarios[id] = usuario
	return nil
}

func (m *UsuarioRepositoryMock) Delete(ctx context.Context, id int) error {
	if _, ok := m.Usuarios[id]; !ok {
		return ErrNotFound
	}
	delete(m.Usuarios, id)
	return nil
}

func (m *UsuarioRepositoryMock) GetByID(ctx context.Context, id int) (*domain.Usuario, error) {
	usuario, ok := m.Usuarios[id]
	if !ok {
		return nil, ErrNotFound
	}
	return usuario, nil
}

func (m *UsuarioRepositoryMock) GetByEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	for _, u := range m.Usuarios {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (m *UsuarioRepositoryMock) GetByDNI(ctx context.Context, dni string) (*domain.Usuario, error) {
	for _, u := range m.Usuarios {
		if u.DNI == dni {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (m *UsuarioRepositoryMock) List(ctx context.Context, filtro map[string]interface{}, offset, limit int) ([]*domain.Usuario, int, error) {
	usuarios := []*domain.Usuario{}
	for _, u := range m.Usuarios {
		usuarios = append(usuarios, u)
	}
	total := len(usuarios)
	end := offset + limit
	if offset > total {
		return []*domain.Usuario{}, total, nil
	}
	if end > total {
		end = total
	}
	return usuarios[offset:end], total, nil
}


var ErrNotFound = errors.New("usuario no encontrado")
