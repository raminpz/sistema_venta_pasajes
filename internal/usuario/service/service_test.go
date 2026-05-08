package service

import (
	"context"
	"errors"
	"sistema_venta_pasajes/internal/usuario/domain"
	"sistema_venta_pasajes/internal/usuario/input"
	"sistema_venta_pasajes/internal/usuario/util"
	"testing"
)

type fakeUsuarioRepo struct {
	createFn     func(context.Context, *domain.Usuario) error
	updateFn     func(context.Context, int, *domain.Usuario) error
	deleteFn     func(context.Context, int) error
	getByIDFn    func(context.Context, int) (*domain.Usuario, error)
	getByEmailFn func(context.Context, string) (*domain.Usuario, error)
	getByDNIFn   func(context.Context, string) (*domain.Usuario, error)
	listFn       func(context.Context, map[string]interface{}, int, int) ([]*domain.Usuario, int, error)
}

func (f *fakeUsuarioRepo) Create(ctx context.Context, usuario *domain.Usuario) error {
	if f.createFn != nil {
		return f.createFn(ctx, usuario)
	}
	return nil
}

func (f *fakeUsuarioRepo) Update(ctx context.Context, id int, usuario *domain.Usuario) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, id, usuario)
	}
	return nil
}

func (f *fakeUsuarioRepo) Delete(ctx context.Context, id int) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

func (f *fakeUsuarioRepo) GetByID(ctx context.Context, id int) (*domain.Usuario, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeUsuarioRepo) GetByEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	if f.getByEmailFn != nil {
		return f.getByEmailFn(ctx, email)
	}
	return nil, nil
}

func (f *fakeUsuarioRepo) GetByDNI(ctx context.Context, dni string) (*domain.Usuario, error) {
	if f.getByDNIFn != nil {
		return f.getByDNIFn(ctx, dni)
	}
	return nil, nil
}

func (f *fakeUsuarioRepo) List(ctx context.Context, filtro map[string]interface{}, offset, limit int) ([]*domain.Usuario, int, error) {
	if f.listFn != nil {
		return f.listFn(ctx, filtro, offset, limit)
	}
	return nil, 0, nil
}

func validCreateInput() input.UsuarioCreateInput {
	return input.UsuarioCreateInput{
		IDRol:     1,
		Nombre:    "juan",
		Apellidos: "perez",
		DNI:       "12345678",
		Email:     "juan@mail.com",
		Password:  "1234",
		Telefono:  "999999999",
	}
}

func TestUsuarioService_Create(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		repo := &fakeUsuarioRepo{}
		serv := NewUsuarioService(repo)
		usuario, err := serv.Create(validCreateInput())
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if usuario.Nombre != "Juan" || usuario.Apellidos != "Perez" {
			t.Fatalf("se esperaba capitalizacion, obtenido %s %s", usuario.Nombre, usuario.Apellidos)
		}
	})

	t.Run("datos invalidos", func(t *testing.T) {
		repo := &fakeUsuarioRepo{}
		serv := NewUsuarioService(repo)
		in := validCreateInput()
		in.Email = "correo-invalido"
		_, err := serv.Create(in)
		if err == nil || err.Error() != util.MSG_INVALID_DATA {
			t.Fatalf("se esperaba error de datos invalidos, obtenido %v", err)
		}
	})

	t.Run("email duplicado", func(t *testing.T) {
		repo := &fakeUsuarioRepo{getByEmailFn: func(context.Context, string) (*domain.Usuario, error) {
			return &domain.Usuario{IDUsuario: 9}, nil
		}}
		serv := NewUsuarioService(repo)
		_, err := serv.Create(validCreateInput())
		if err == nil || err.Error() != util.MSG_EMAIL_DUPLICATE {
			t.Fatalf("se esperaba duplicado email, obtenido %v", err)
		}
	})

	t.Run("dni duplicado", func(t *testing.T) {
		repo := &fakeUsuarioRepo{getByDNIFn: func(context.Context, string) (*domain.Usuario, error) {
			return &domain.Usuario{IDUsuario: 10}, nil
		}}
		serv := NewUsuarioService(repo)
		_, err := serv.Create(validCreateInput())
		if err == nil || err.Error() != util.MSG_DNI_DUPLICATE {
			t.Fatalf("se esperaba duplicado dni, obtenido %v", err)
		}
	})
}

func TestUsuarioService_UpdateDeleteGetList(t *testing.T) {
	repo := &fakeUsuarioRepo{
		getByIDFn: func(_ context.Context, id int) (*domain.Usuario, error) {
			if id == 404 {
				return nil, nil
			}
			return &domain.Usuario{IDUsuario: id, Email: "x@mail.com", Nombre: "juan", Apellidos: "diaz", Estado: "ACTIVO"}, nil
		},
		listFn: func(_ context.Context, _ map[string]interface{}, _, _ int) ([]*domain.Usuario, int, error) {
			return []*domain.Usuario{{IDUsuario: 1, Nombre: "A", Apellidos: "B", Email: "a@mail.com", Estado: "ACTIVO"}}, 1, nil
		},
	}
	serv := NewUsuarioService(repo)

	t.Run("update usuario inexistente", func(t *testing.T) {
		_, err := serv.Update(404, input.UsuarioUpdateInput{})
		if err == nil || err.Error() != util.MSG_USER_NOT_FOUND {
			t.Fatalf("se esperaba no encontrado, obtenido %v", err)
		}
	})

	t.Run("update ok", func(t *testing.T) {
		out, err := serv.Update(1, input.UsuarioUpdateInput{Nombre: "luis", Apellidos: "perez", Email: "x@mail.com", Telefono: "999999999", Estado: "ACTIVO"})
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if out.Nombre != "Luis" {
			t.Fatalf("se esperaba nombre capitalizado, obtenido %s", out.Nombre)
		}
	})

	t.Run("delete not found", func(t *testing.T) {
		err := serv.Delete(404)
		if err == nil || err.Error() != util.MSG_USER_NOT_FOUND {
			t.Fatalf("se esperaba not found, obtenido %v", err)
		}
	})

	t.Run("delete error repo", func(t *testing.T) {
		repo.deleteFn = func(context.Context, int) error { return errors.New("db") }
		err := serv.Delete(1)
		if err == nil {
			t.Fatal("se esperaba error")
		}
		repo.deleteFn = nil
	})

	t.Run("get by id ok", func(t *testing.T) {
		out, err := serv.GetByID(1)
		if err != nil || out.IDUsuario != 1 {
			t.Fatalf("resultado inesperado: %+v err=%v", out, err)
		}
	})

	t.Run("list ok", func(t *testing.T) {
		out, total, err := serv.List(1, 10)
		if err != nil || total != 1 || len(out) != 1 {
			t.Fatalf("resultado inesperado: total=%d len=%d err=%v", total, len(out), err)
		}
	})
}
