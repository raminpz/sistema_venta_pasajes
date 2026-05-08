package service

import (
	"context"
	"errors"
	"sistema_venta_pasajes/pkg"
	"testing"

	domain "sistema_venta_pasajes/internal/proveedor/domain"
	providerinput "sistema_venta_pasajes/internal/proveedor/input"
)

type fakeRepository struct {
	listFn    func(ctx context.Context) ([]domain.ProveedorSistema, error)
	getByIDFn func(ctx context.Context, id int64) (*domain.ProveedorSistema, error)
	createFn  func(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error)
	updateFn  func(ctx context.Context, id int64, input providerinput.UpdateInput) (*domain.ProveedorSistema, error)
	deleteFn  func(ctx context.Context, id int64) error
}

func (f fakeRepository) List(ctx context.Context) ([]domain.ProveedorSistema, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, nil
}

func (f fakeRepository) GetByID(ctx context.Context, id int64) (*domain.ProveedorSistema, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (f fakeRepository) Create(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error) {
	if f.createFn != nil {
		return f.createFn(ctx, input)
	}
	return nil, nil
}

func (f fakeRepository) Update(ctx context.Context, id int64, input providerinput.UpdateInput) (*domain.ProveedorSistema, error) {
	if f.updateFn != nil {
		return f.updateFn(ctx, id, input)
	}
	return nil, nil
}

func (f fakeRepository) Delete(ctx context.Context, id int64) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

func TestServiceCreateNormalizesAndDelegates(t *testing.T) {
	var received providerinput.CreateInput
	svc := NewService(fakeRepository{
		createFn: func(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error) {
			received = input
			return &domain.ProveedorSistema{IDProveedor: 1, RUC: input.RUC, RazonSocial: input.RazonSocial, Email: input.Email, Web: input.Web}, nil
		},
	})

	proveedor, err := svc.Create(context.Background(), providerinput.CreateInput{
		RUC:         " 20123456789 ",
		RazonSocial: "  Transportes del Norte SAC ",
		Email:       "  VENTAS@EMPRESA.COM ",
		Web:         "https://empresa.com",
	})
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo %v", err)
	}

	if received.RUC != "20123456789" {
		t.Fatalf("se esperaba el RUC normalizado, se obtuvo %q", received.RUC)
	}
	if received.RazonSocial != "Transportes del Norte SAC" {
		t.Fatalf("se esperaba la razón social sin espacios extra, se obtuvo %q", received.RazonSocial)
	}
	if received.Email != "ventas@empresa.com" {
		t.Fatalf("se esperaba el email en minúsculas, se obtuvo %q", received.Email)
	}
	if proveedor.IDProveedor != 1 {
		t.Fatalf("se esperaba el ID de proveedor 1, se obtuvo %d", proveedor.IDProveedor)
	}
}

func TestServiceCreateReturnsValidationError(t *testing.T) {
	repositoryCalled := false
	svc := NewService(fakeRepository{
		createFn: func(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error) {
			repositoryCalled = true
			return nil, errors.New("no debería ejecutarse")
		},
	})

	_, err := svc.Create(context.Background(), providerinput.CreateInput{
		RUC:         "123",
		RazonSocial: "",
		Email:       "correo-invalido",
		Web:         "empresa.com",
	})
	if err == nil {
		t.Fatal("se esperaba un error de validación, se obtuvo nil")
	}
	if repositoryCalled {
		t.Fatal("no se esperaba que el repositorio se ejecutara cuando falla la validación")
	}

	appErr := pkg.AsAppError(err)
	if appErr.Code != "validation_error" {
		t.Fatalf("se esperaba validation_error, se obtuvo %q", appErr.Code)
	}

	details, ok := appErr.Details.(map[string]string)
	if !ok {
		t.Fatalf("se esperaba details como map[string]string, se obtuvo %#v", appErr.Details)
	}
	if details["ruc"] == "" || details["razon_social"] == "" || details["email"] == "" || details["web"] == "" {
		t.Fatalf("se esperaban detalles de validación para ruc, razon_social, email y web, se obtuvo %#v", details)
	}
}

func TestServiceGetByIDRejectsInvalidID(t *testing.T) {
	svc := NewService(fakeRepository{})

	_, err := svc.GetByID(context.Background(), 0)
	if err == nil {
		t.Fatal("se esperaba un error, se obtuvo nil")
	}

	appErr := pkg.AsAppError(err)
	if appErr.Code != "invalid_provider_id" {
		t.Fatalf("se esperaba invalid_provider_id, se obtuvo %q", appErr.Code)
	}
}

func TestServiceDeleteRejectsInvalidID(t *testing.T) {
	svc := NewService(fakeRepository{})

	err := svc.Delete(context.Background(), 0)
	if err == nil {
		t.Fatal("se esperaba un error, se obtuvo nil")
	}

	appErr := pkg.AsAppError(err)
	if appErr.Code != "invalid_provider_id" {
		t.Fatalf("se esperaba invalid_provider_id, se obtuvo %q", appErr.Code)
	}
}

func TestServiceListAndUpdateDeleteBranches(t *testing.T) {
	t.Run("list nil retorna slice vacio", func(t *testing.T) {
		svc := NewService(fakeRepository{listFn: func(ctx context.Context) ([]domain.ProveedorSistema, error) {
			return nil, nil
		}})
		list, err := svc.List(context.Background())
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if list == nil || len(list) != 0 {
			t.Fatalf("se esperaba slice vacio, obtuvo %#v", list)
		}
	})

	t.Run("list error", func(t *testing.T) {
		svc := NewService(fakeRepository{listFn: func(ctx context.Context) ([]domain.ProveedorSistema, error) {
			return nil, errors.New("db")
		}})
		if _, err := svc.List(context.Background()); err == nil {
			t.Fatal("se esperaba error")
		}
	})

	t.Run("update invalid id", func(t *testing.T) {
		svc := NewService(fakeRepository{})
		_, err := svc.Update(context.Background(), 0, providerinput.UpdateInput{})
		if err == nil {
			t.Fatal("se esperaba error")
		}
	})

	t.Run("update normaliza email", func(t *testing.T) {
		var received providerinput.UpdateInput
		svc := NewService(fakeRepository{updateFn: func(ctx context.Context, id int64, in providerinput.UpdateInput) (*domain.ProveedorSistema, error) {
			received = in
			return &domain.ProveedorSistema{IDProveedor: id, Email: in.Email}, nil
		}})
		_, err := svc.Update(context.Background(), 1, providerinput.UpdateInput{RUC: "20123456789", RazonSocial: "Empresa", Email: "INFO@MAIL.COM"})
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if received.Email != "info@mail.com" {
			t.Fatalf("se esperaba email normalizado, obtuvo %q", received.Email)
		}
	})

	t.Run("delete delega al repo", func(t *testing.T) {
		called := false
		svc := NewService(fakeRepository{deleteFn: func(ctx context.Context, id int64) error {
			called = true
			return nil
		}})
		if err := svc.Delete(context.Background(), 1); err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if !called {
			t.Fatal("se esperaba invocacion al repositorio")
		}
	})
}

func TestServiceCreateDuplicateMapping(t *testing.T) {
	t.Run("duplicado ruc", func(t *testing.T) {
		svc := NewService(fakeRepository{createFn: func(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error) {
			return nil, pkg.NewAppError(409, "duplicate_resource", "dup").WithCause(errors.New("Duplicate entry RUC"))
		}})
		_, err := svc.Create(context.Background(), providerinput.CreateInput{RUC: "20123456789", RazonSocial: "Empresa"})
		if err == nil {
			t.Fatal("se esperaba error")
		}
	})

	t.Run("duplicado email", func(t *testing.T) {
		svc := NewService(fakeRepository{createFn: func(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error) {
			return nil, pkg.NewAppError(409, "duplicate_resource", "dup").WithCause(errors.New("Duplicate entry EMAIL"))
		}})
		_, err := svc.Create(context.Background(), providerinput.CreateInput{RUC: "20123456789", RazonSocial: "Empresa", Email: "a@mail.com"})
		if err == nil {
			t.Fatal("se esperaba error")
		}
	})
}

func TestValidateWebURL(t *testing.T) {
	if err := validateWebURL("empresa.com"); err == nil {
		t.Fatal("se esperaba error por URL invalida")
	}
	if err := validateWebURL("ftp://empresa.com"); err == nil {
		t.Fatal("se esperaba error por esquema no permitido")
	}
	if err := validateWebURL("https://empresa.com"); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}
