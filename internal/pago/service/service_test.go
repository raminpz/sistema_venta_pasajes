package service

import (
	"errors"
	"sistema_venta_pasajes/internal/pago/domain"
	"sistema_venta_pasajes/internal/pago/input"
	"sistema_venta_pasajes/internal/pago/util"
	"testing"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

type fakePagoRepo struct {
	createFn  func(*domain.Pago) error
	updateFn  func(*domain.Pago) error
	deleteFn  func(int64) error
	getByIDFn func(int64) (*domain.Pago, error)
	listFn    func(int, int, *int64) ([]domain.Pago, int, error)
}

func (f *fakePagoRepo) Create(p *domain.Pago) error {
	if f.createFn != nil {
		return f.createFn(p)
	}
	return nil
}
func (f *fakePagoRepo) Update(p *domain.Pago) error {
	if f.updateFn != nil {
		return f.updateFn(p)
	}
	return nil
}
func (f *fakePagoRepo) Delete(id int64) error {
	if f.deleteFn != nil {
		return f.deleteFn(id)
	}
	return nil
}
func (f *fakePagoRepo) GetByID(id int64) (*domain.Pago, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(id)
	}
	return nil, nil
}
func (f *fakePagoRepo) List(offset, limit int, idVenta *int64) ([]domain.Pago, int, error) {
	if f.listFn != nil {
		return f.listFn(offset, limit, idVenta)
	}
	return nil, 0, nil
}

// ---------- Create ----------

func TestServiceCreateEstadoDefectoRegistrada(t *testing.T) {
	repo := &fakePagoRepo{createFn: func(p *domain.Pago) error {
		if p.Estado != util.STATUS_REGISTRADA {
			t.Fatalf("estado esperado REGISTRADA, obtuvo %s", p.Estado)
		}
		return nil
	}}
	s := NewPagoService(repo)
	out, err := s.Create(input.CreatePagoInput{IDVenta: 1, IDMetodo: 1, Monto: 0})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out == nil || out.IDVenta != 1 {
		t.Fatal("salida invalida")
	}
}

func TestServiceCreateEstadoParcial(t *testing.T) {
	repo := &fakePagoRepo{createFn: func(p *domain.Pago) error {
		if p.Estado != util.STATUS_PARCIAL {
			t.Fatalf("estado esperado PARCIAL, obtuvo %s", p.Estado)
		}
		return nil
	}}
	s := NewPagoService(repo)
	out, err := s.Create(input.CreatePagoInput{IDVenta: 1, IDMetodo: 1, Monto: 50, Estado: util.STATUS_PARCIAL})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out == nil {
		t.Fatal("salida invalida")
	}
}

func TestServiceCreateEstadoPagada(t *testing.T) {
	repo := &fakePagoRepo{}
	s := NewPagoService(repo)
	out, err := s.Create(input.CreatePagoInput{IDVenta: 1, IDMetodo: 2, Monto: 100, Estado: util.STATUS_PAGADA})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out == nil || out.Estado != util.STATUS_PAGADA {
		t.Fatal("salida invalida o estado incorrecto")
	}
}

func TestServiceCreateEstadoAnulada(t *testing.T) {
	repo := &fakePagoRepo{}
	s := NewPagoService(repo)
	out, err := s.Create(input.CreatePagoInput{IDVenta: 1, IDMetodo: 1, Monto: 0, Estado: util.STATUS_ANULADA})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out == nil || out.Estado != util.STATUS_ANULADA {
		t.Fatal("estado incorrecto")
	}
}

func TestServiceCreateEstadoInvalido(t *testing.T) {
	s := NewPagoService(&fakePagoRepo{})
	_, err := s.Create(input.CreatePagoInput{IDVenta: 1, IDMetodo: 1, Monto: 0, Estado: "PENDIENTE"})
	if err == nil {
		t.Fatal("se esperaba error por estado invalido")
	}
}

func TestServiceCreateMontoNegativo(t *testing.T) {
	s := NewPagoService(&fakePagoRepo{})
	_, err := s.Create(input.CreatePagoInput{IDVenta: 1, IDMetodo: 1, Monto: -1})
	if err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestServiceCreateSinIDVenta(t *testing.T) {
	s := NewPagoService(&fakePagoRepo{})
	_, err := s.Create(input.CreatePagoInput{IDVenta: 0, IDMetodo: 1, Monto: 0})
	if err == nil {
		t.Fatal("se esperaba error por id_venta invalido")
	}
}

func TestServiceCreateSinIDMetodo(t *testing.T) {
	s := NewPagoService(&fakePagoRepo{})
	_, err := s.Create(input.CreatePagoInput{IDVenta: 1, IDMetodo: 0, Monto: 0})
	if err == nil {
		t.Fatal("se esperaba error por id_metodo invalido")
	}
}

// ---------- Update ----------

func TestServiceUpdateNotFound(t *testing.T) {
	repo := &fakePagoRepo{getByIDFn: func(int64) (*domain.Pago, error) {
		return nil, errors.New("not found")
	}}
	s := NewPagoService(repo)
	_, err := s.Update(1, input.UpdatePagoInput{})
	if err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestServiceUpdateEstadoValido(t *testing.T) {
	estadoPagada := util.STATUS_PAGADA
	repo := &fakePagoRepo{
		getByIDFn: func(id int64) (*domain.Pago, error) {
			return &domain.Pago{IDPago: 1, IDVenta: 1, IDMetodo: 1, Monto: 100, Estado: util.STATUS_REGISTRADA}, nil
		},
	}
	s := NewPagoService(repo)
	out, err := s.Update(1, input.UpdatePagoInput{Estado: &estadoPagada})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out == nil || out.Estado != util.STATUS_PAGADA {
		t.Fatal("estado no actualizado correctamente")
	}
}

func TestServiceUpdateEstadoInvalido(t *testing.T) {
	estadoInvalido := "PAGADO"
	repo := &fakePagoRepo{
		getByIDFn: func(id int64) (*domain.Pago, error) {
			return &domain.Pago{IDPago: 1, Estado: util.STATUS_REGISTRADA}, nil
		},
	}
	s := NewPagoService(repo)
	_, err := s.Update(1, input.UpdatePagoInput{Estado: &estadoInvalido})
	if err == nil {
		t.Fatal("se esperaba error por estado invalido")
	}
}

// ---------- List ----------

func TestServiceListOK(t *testing.T) {
	repo := &fakePagoRepo{listFn: func(offset, limit int, idVenta *int64) ([]domain.Pago, int, error) {
		if offset != 0 || limit != 15 {
			t.Fatalf("offset/limit inesperados %d/%d", offset, limit)
		}
		return []domain.Pago{{IDPago: 1, IDVenta: 1, IDMetodo: 1, Monto: 10, Estado: util.STATUS_PAGADA}}, 1, nil
	}}
	s := NewPagoService(repo)
	outs, total, err := s.List(1, 15, nil)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if total != 1 || len(outs) != 1 {
		t.Fatal("resultado inesperado")
	}
}

func TestServiceListFiltradoPorVenta(t *testing.T) {
	idVenta := int64(5)
	repo := &fakePagoRepo{listFn: func(offset, limit int, v *int64) ([]domain.Pago, int, error) {
		if v == nil || *v != 5 {
			t.Fatal("id_venta no se pasó correctamente")
		}
		return []domain.Pago{{IDPago: 2, IDVenta: 5, Estado: util.STATUS_PARCIAL}}, 1, nil
	}}
	s := NewPagoService(repo)
	outs, total, err := s.List(1, 15, &idVenta)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if total != 1 || len(outs) != 1 {
		t.Fatal("resultado inesperado")
	}
}

// ---------- Delete ----------

func TestServiceDeleteOK(t *testing.T) {
	s := NewPagoService(&fakePagoRepo{deleteFn: func(int64) error { return nil }})
	if err := s.Delete(1); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
}

func TestServiceDeleteError(t *testing.T) {
	s := NewPagoService(&fakePagoRepo{deleteFn: func(int64) error { return errors.New("db") }})
	if err := s.Delete(1); err == nil {
		t.Fatal("se esperaba error")
	}
}

func TestServiceCreateIDVentaNoExiste(t *testing.T) {
	repo := &fakePagoRepo{createFn: func(p *domain.Pago) error {
		return &mysqlDriver.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`SISTEMA_PASAJES`.`PAGO`, CONSTRAINT `FK_PAGO_VENTA` FOREIGN KEY (`ID_VENTA`) REFERENCES `VENTA` (`ID_VENTA`))",
		}
	}}
	s := NewPagoService(repo)
	_, err := s.Create(input.CreatePagoInput{IDVenta: 999, IDMetodo: 1, Monto: 0})
	if err == nil {
		t.Fatal("se esperaba error")
	}
	if err.Error() != util.MSG_PAGO_VENTA_NOT_FOUND {
		t.Fatalf("mensaje esperado: %q, obtuvo: %q", util.MSG_PAGO_VENTA_NOT_FOUND, err.Error())
	}
}

func TestServiceCreateIDMetodoNoExiste(t *testing.T) {
	repo := &fakePagoRepo{createFn: func(p *domain.Pago) error {
		return &mysqlDriver.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`SISTEMA_PASAJES`.`PAGO`, CONSTRAINT `FK_PAGO_METODO` FOREIGN KEY (`ID_METODO`) REFERENCES `METODO_PAGO` (`ID_METODO`))",
		}
	}}
	s := NewPagoService(repo)
	_, err := s.Create(input.CreatePagoInput{IDVenta: 1, IDMetodo: 999, Monto: 0})
	if err == nil {
		t.Fatal("se esperaba error")
	}
	if err.Error() != util.MSG_PAGO_METODO_NOT_FOUND {
		t.Fatalf("mensaje esperado: %q, obtuvo: %q", util.MSG_PAGO_METODO_NOT_FOUND, err.Error())
	}
}

func TestServiceCreateEnumDBNoActualizado(t *testing.T) {
	repo := &fakePagoRepo{createFn: func(p *domain.Pago) error {
		return &mysqlDriver.MySQLError{
			Number:  1265,
			Message: "Data truncated for column 'ESTADO' at row 1",
		}
	}}
	s := NewPagoService(repo)
	_, err := s.Create(input.CreatePagoInput{IDVenta: 1, IDMetodo: 1, Monto: 0, Estado: util.STATUS_PARCIAL})
	if err == nil {
		t.Fatal("se esperaba error")
	}
	if err.Error() != util.MSG_PAGO_ENUM_DB_ERROR {
		t.Fatalf("mensaje esperado: %q, obtuvo: %q", util.MSG_PAGO_ENUM_DB_ERROR, err.Error())
	}
}

func TestServiceUpdateIDMetodoNoExiste(t *testing.T) {
	nuevoMetodo := int64(999)
	repo := &fakePagoRepo{
		getByIDFn: func(id int64) (*domain.Pago, error) {
			return &domain.Pago{IDPago: 1, IDVenta: 1, IDMetodo: 1, Estado: util.STATUS_REGISTRADA}, nil
		},
		updateFn: func(p *domain.Pago) error {
			return &mysqlDriver.MySQLError{
				Number:  1452,
				Message: "CONSTRAINT `FK_PAGO_METODO` FOREIGN KEY (`ID_METODO`) REFERENCES `METODO_PAGO`",
			}
		},
	}
	s := NewPagoService(repo)
	_, err := s.Update(1, input.UpdatePagoInput{IDMetodo: &nuevoMetodo})
	if err == nil {
		t.Fatal("se esperaba error")
	}
	if err.Error() != util.MSG_PAGO_METODO_NOT_FOUND {
		t.Fatalf("mensaje esperado: %q, obtuvo: %q", util.MSG_PAGO_METODO_NOT_FOUND, err.Error())
	}
}
