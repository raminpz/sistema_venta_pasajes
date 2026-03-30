package service

import (
	"os"
	"sistema_venta_pasajes/internal/empresa/domain"
	"sistema_venta_pasajes/internal/empresa/input"
	"sistema_venta_pasajes/internal/empresa/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:3306)/SISTEMA_PASAJES?parseTime=true"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Exec("DELETE FROM EMPRESA")
	return db
}

func TestEmpresaService_Create_Validations(t *testing.T) {
	repo := &repository.EmpresaRepositoryMock{}
	svc := NewEmpresaService(repo)

	// RUC inválido
	in := input.CreateEmpresaInput{
		RUC:           "123",
		RazonSocial:   "Empresa Test",
		Telefono:      "987654321",
		FechaCreacion: time.Now().Format("2006-01-02"),
	}
	_, err := svc.Create(in)
	assert.Error(t, err, "Debe fallar si el RUC es inválido")

	// Teléfono inválido
	in.RUC = "12345678901"
	in.Telefono = "123"
	_, err = svc.Create(in)
	assert.Error(t, err, "Debe fallar si el teléfono es inválido")

	// Creación exitosa
	in.Telefono = "987654321"
	repo.On("Create", mock.Anything).Return(nil)
	out, err := svc.Create(in)
	assert.NoError(t, err, "No debe fallar creación válida")
	assert.Equal(t, in.RUC, out.RUC)
}

func TestEmpresaService_CRUD(t *testing.T) {
	repo := &repository.EmpresaRepositoryMock{}
	svc := NewEmpresaService(repo)

	// Crear empresa válida
	in := input.CreateEmpresaInput{
		RUC:           "11111111111",
		RazonSocial:   "Empresa Uno",
		Telefono:      "987654321",
		FechaCreacion: time.Now().Format("2006-01-02"),
	}
	repo.On("Create", mock.Anything).Return(nil)
	out, err := svc.Create(in)
	assert.NoError(t, err, "Error al crear empresa")

	// GetByID existente
	domainEmpresa := &domain.Empresa{
		IDEmpresa:     int(out.IDEmpresa),
		RUC:           in.RUC,
		RazonSocial:   in.RazonSocial,
		Telefono:      in.Telefono,
		FechaCreacion: time.Now(),
	}
	repo.On("GetByID", out.IDEmpresa).Return(domainEmpresa, nil)
	got, err := svc.GetByID(out.IDEmpresa)
	assert.NoError(t, err, "Error al obtener empresa")
	assert.Equal(t, in.RUC, got.RUC)

	// Update exitoso
	repo.On("Update", mock.Anything).Return(nil)
	upd := input.UpdateEmpresaInput{
		RazonSocial: "Empresa Actualizada",
		Telefono:    "912345678",
	}
	updated, err := svc.Update(out.IDEmpresa, upd)
	assert.NoError(t, err, "Error al actualizar empresa")
	assert.Equal(t, "Empresa Actualizada", updated.RazonSocial)

	// Update con datos inválidos
	upd.Telefono = "123"
	_, err = svc.Update(out.IDEmpresa, upd)
	assert.Error(t, err, "Debe fallar si el teléfono es inválido en update")

	// Delete exitoso
	repo.On("Delete", out.IDEmpresa).Return(nil)
	err = svc.Delete(out.IDEmpresa)
	assert.NoError(t, err, "Error al eliminar empresa")

	// GetByID de empresa eliminada
	repo.ExpectedCalls = nil // Limpiar expectativas previas
	repo.On("GetByID", out.IDEmpresa).Return((*domain.Empresa)(nil), gorm.ErrRecordNotFound).Once()
	_, err = svc.GetByID(out.IDEmpresa)
	assert.Error(t, err, "No debe encontrar empresa eliminada")

	// Delete de empresa inexistente
	repo.On("Delete", int64(999999)).Return(gorm.ErrRecordNotFound)
	err = svc.Delete(999999)
	assert.Error(t, err, "Debe fallar al eliminar empresa inexistente")

	// List solo empresas no eliminadas
	in2 := input.CreateEmpresaInput{
		RUC:           "22222222222",
		RazonSocial:   "Empresa Dos",
		Telefono:      "987654322",
		FechaCreacion: time.Now().Format("2006-01-02"),
	}
	repo.On("Create", mock.Anything).Return(nil)
	_, _ = svc.Create(in2)
	repo.On("List").Return([]domain.Empresa{
		{IDEmpresa: 2, RUC: in2.RUC, RazonSocial: in2.RazonSocial, Telefono: in2.Telefono, FechaCreacion: time.Now()},
	}, nil)
	list, err := svc.List()
	assert.NoError(t, err, "Error al listar empresas")
	for _, e := range list {
		assert.NotEqual(t, out.IDEmpresa, e.IDEmpresa, "No debe listar empresas eliminadas")
	}
}
