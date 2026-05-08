package service

import (
	"errors"
	"sistema_venta_pasajes/internal/encomienda/domain"
	"sistema_venta_pasajes/internal/encomienda/input"
	"sistema_venta_pasajes/internal/encomienda/util"
	"testing"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type fakeEncomiendaRepo struct {
	createFn  func(*domain.Encomienda) error
	updateFn  func(*domain.Encomienda) error
	deleteFn  func(int64) error
	getByIDFn func(int64) (*domain.Encomienda, error)
	listFn    func(int, int) ([]domain.Encomienda, int, error)
}

func (f *fakeEncomiendaRepo) Create(e *domain.Encomienda) error {
	if f.createFn != nil {
		return f.createFn(e)
	}
	return nil
}

func (f *fakeEncomiendaRepo) Update(e *domain.Encomienda) error {
	if f.updateFn != nil {
		return f.updateFn(e)
	}
	return nil
}

func (f *fakeEncomiendaRepo) Delete(id int64) error {
	if f.deleteFn != nil {
		return f.deleteFn(id)
	}
	return nil
}

func (f *fakeEncomiendaRepo) GetByID(id int64) (*domain.Encomienda, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(id)
	}
	return nil, nil
}

func (f *fakeEncomiendaRepo) List(offset, limit int) ([]domain.Encomienda, int, error) {
	if f.listFn != nil {
		return f.listFn(offset, limit)
	}
	return nil, 0, nil
}

func TestServiceCreateOK(t *testing.T) {
	repo := &fakeEncomiendaRepo{createFn: func(e *domain.Encomienda) error {
		if e.Estado != util.STATUS_PENDIENTE {
			t.Fatalf("estado esperado PENDIENTE, obtuvo %s", e.Estado)
		}
		if e.RemitenteNombre != "Juan Perez" {
			t.Fatalf("capitalizacion esperada en remitente, obtuvo %s", e.RemitenteNombre)
		}
		return nil
	}}

	s := NewEncomiendaService(repo)
	out, err := s.Create(input.CreateEncomiendaInput{
		IDVenta:            1,
		IDProgramacion:     1,
		Costo:              50,
		RemitenteNombre:    " juan perez ",
		RemitenteDoc:       "12345678",
		DestinatarioNombre: " maria lopez ",
		DestinatarioTel:    "987654321",
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out == nil {
		t.Fatal("salida invalida")
	}
}

func TestServiceCreateInvalidEstado(t *testing.T) {
	s := NewEncomiendaService(&fakeEncomiendaRepo{})
	_, err := s.Create(input.CreateEncomiendaInput{
		IDVenta:            1,
		IDProgramacion:     1,
		Costo:              25,
		RemitenteNombre:    "A",
		RemitenteDoc:       "12345678",
		DestinatarioNombre: "B",
		DestinatarioTel:    "987654321",
		Estado:             "CANCELADO",
	})
	if err == nil {
		t.Fatal("se esperaba error por estado invalido")
	}
}

func TestServiceCreateEstadoEnCursoInvalido(t *testing.T) {
	s := NewEncomiendaService(&fakeEncomiendaRepo{})
	_, err := s.Create(input.CreateEncomiendaInput{
		IDVenta:            1,
		IDProgramacion:     1,
		Costo:              25,
		RemitenteNombre:    "Juan",
		RemitenteDoc:       "12345678",
		DestinatarioNombre: "Maria",
		DestinatarioTel:    "987654321",
		Estado:             "EN_CURSO",
	})
	if err == nil {
		t.Fatal("se esperaba error por estado EN_CURSO invalido")
	}
}

func TestServiceCreateDBForeignKeyError(t *testing.T) {
	repo := &fakeEncomiendaRepo{createFn: func(e *domain.Encomienda) error {
		return &mysqlDriver.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`SISTEMA_PASAJES`.`ENCOMIENDA`, CONSTRAINT `FK_ENCOMIENDA_VENTA` FOREIGN KEY (`ID_VENTA`) REFERENCES `VENTA` (`ID_VENTA`))",
		}
	}}
	s := NewEncomiendaService(repo)
	_, err := s.Create(input.CreateEncomiendaInput{
		IDVenta:            999,
		IDProgramacion:     1,
		Costo:              50,
		RemitenteNombre:    "Juan",
		RemitenteDoc:       "12345678",
		DestinatarioNombre: "Maria",
		DestinatarioTel:    "987654321",
	})
	if err == nil {
		t.Fatal("se esperaba error")
	}
	if err.Error() != util.MSG_ENCOMIENDA_VENTA_NOT_FOUND {
		t.Fatalf("mensaje esperado %q, obtuvo %q", util.MSG_ENCOMIENDA_VENTA_NOT_FOUND, err.Error())
	}
}

func TestServiceUpdateNotFound(t *testing.T) {
	repo := &fakeEncomiendaRepo{getByIDFn: func(int64) (*domain.Encomienda, error) {
		return nil, gorm.ErrRecordNotFound
	}}
	s := NewEncomiendaService(repo)
	_, err := s.Update(1, input.UpdateEncomiendaInput{})
	if err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestServiceListOK(t *testing.T) {
	repo := &fakeEncomiendaRepo{listFn: func(offset, limit int) ([]domain.Encomienda, int, error) {
		if offset != 15 || limit != 15 {
			t.Fatalf("offset/limit esperados 15/15, obtuvo %d/%d", offset, limit)
		}
		return []domain.Encomienda{{IDEncomienda: 1, Estado: util.STATUS_PENDIENTE}}, 20, nil
	}}
	s := NewEncomiendaService(repo)
	items, total, err := s.List(2, 15)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if total != 20 || len(items) != 1 {
		t.Fatalf("resultado inesperado total=%d len=%d", total, len(items))
	}
}

func TestServiceDeleteInvalidID(t *testing.T) {
	s := NewEncomiendaService(&fakeEncomiendaRepo{})
	if err := s.Delete(0); err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestServiceDeleteRepositoryError(t *testing.T) {
	s := NewEncomiendaService(&fakeEncomiendaRepo{deleteFn: func(int64) error {
		return errors.New("db")
	}})
	if err := s.Delete(1); err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestServiceCreateRequiredDocAndTel(t *testing.T) {
	s := NewEncomiendaService(&fakeEncomiendaRepo{})
	_, err := s.Create(input.CreateEncomiendaInput{
		IDVenta:            1,
		IDProgramacion:     1,
		Costo:              50,
		RemitenteNombre:    "Juan",
		DestinatarioNombre: "Maria",
	})
	if err == nil {
		t.Fatal("se esperaba error por remitente_doc y destinatario_tel obligatorios")
	}
}

func TestServiceUpdateGetAndHelpersBranches(t *testing.T) {
	now := time.Now()
	repo := &fakeEncomiendaRepo{
		getByIDFn: func(id int64) (*domain.Encomienda, error) {
			if id == 404 {
				return nil, gorm.ErrRecordNotFound
			}
			return &domain.Encomienda{IDEncomienda: id, IDVenta: 1, IDProgramacion: 1, Costo: 10, RemitenteNombre: "Juan", RemitenteDoc: "123", DestinatarioNombre: "Maria", DestinatarioTel: "987654321", Estado: util.STATUS_PENDIENTE, CreatedAt: &now, UpdatedAt: &now}, nil
		},
		updateFn: func(e *domain.Encomienda) error {
			if e.IDEncomienda == 500 {
				return errors.New("db")
			}
			return nil
		},
		listFn: func(offset, limit int) ([]domain.Encomienda, int, error) {
			if offset == 999 {
				return nil, 0, errors.New("db")
			}
			return []domain.Encomienda{{IDEncomienda: 1, IDVenta: 1, IDProgramacion: 1, Costo: 10, RemitenteNombre: "Juan", RemitenteDoc: "123", DestinatarioNombre: "Maria", DestinatarioTel: "987654321", Estado: util.STATUS_PENDIENTE, CreatedAt: &now, UpdatedAt: &now}}, 1, nil
		},
	}

	s := NewEncomiendaService(repo)

	t.Run("update id invalido", func(t *testing.T) {
		_, err := s.Update(0, input.UpdateEncomiendaInput{})
		if err == nil {
			t.Fatal("se esperaba error")
		}
	})

	t.Run("update validacion", func(t *testing.T) {
		_, err := s.Update(1, input.UpdateEncomiendaInput{})
		if err == nil {
			t.Fatal("se esperaba error de validacion")
		}
	})

	t.Run("update not found", func(t *testing.T) {
		_, err := s.Update(404, input.UpdateEncomiendaInput{IDVenta: 1, IDProgramacion: 1, Costo: 10, RemitenteNombre: "A", RemitenteDoc: "123", DestinatarioNombre: "B", DestinatarioTel: "987654321", Estado: util.STATUS_PENDIENTE})
		if err == nil {
			t.Fatal("se esperaba not found")
		}
	})

	t.Run("update ok", func(t *testing.T) {
		out, err := s.Update(1, input.UpdateEncomiendaInput{IDVenta: 1, IDProgramacion: 1, Costo: 10, RemitenteNombre: "juan perez", RemitenteDoc: "123", DestinatarioNombre: "maria", DestinatarioTel: "987654321", Estado: util.STATUS_PENDIENTE})
		if err != nil || out == nil {
			t.Fatalf("resultado inesperado out=%v err=%v", out, err)
		}
	})

	t.Run("get id invalido", func(t *testing.T) {
		_, err := s.GetByID(0)
		if err == nil {
			t.Fatal("se esperaba error")
		}
	})

	t.Run("get not found", func(t *testing.T) {
		_, err := s.GetByID(404)
		if err == nil {
			t.Fatal("se esperaba not found")
		}
	})

	t.Run("get ok", func(t *testing.T) {
		out, err := s.GetByID(1)
		if err != nil || out.IDEncomienda != 1 {
			t.Fatalf("resultado inesperado out=%v err=%v", out, err)
		}
	})

	t.Run("list ok", func(t *testing.T) {
		items, total, err := s.List(1, 15)
		if err != nil || total != 1 || len(items) != 1 {
			t.Fatalf("resultado inesperado total=%d len=%d err=%v", total, len(items), err)
		}
	})
}

func TestValidateUpdateAndHelpers(t *testing.T) {
	doc := ""
	if err := validateUpdateEncomiendaInput(input.UpdateEncomiendaInput{
		IDVenta:            1,
		IDProgramacion:     1,
		Costo:              20,
		RemitenteNombre:    "Juan",
		RemitenteDoc:       "12345678",
		DestinatarioNombre: "Maria",
		DestinatarioTel:    "987654321",
		DestinatarioDoc:    &doc,
		Estado:             util.STATUS_PENDIENTE,
	}); err == nil {
		t.Fatal("se esperaba error por destinatario_doc vacio")
	}

	if isDigits("123456789") != true {
		t.Fatal("se esperaba true para solo digitos")
	}
	if isDigits("12A") != false {
		t.Fatal("se esperaba false para no digitos")
	}
}
