package service

import (
	"errors"
	"sistema_venta_pasajes/internal/venta/domain"
	"sistema_venta_pasajes/internal/venta/input"
	"sistema_venta_pasajes/internal/venta/util"
	"testing"
)

// ---------------------------------------------------------------------------
// fakeRepo mock
// ---------------------------------------------------------------------------

type fakeRepo struct {
	createErr          error
	updateErr          error
	deleteErr          error
	getByIDErr         error
	listErr            error
	nextCorrelativoErr error
	nextCorrelativo    uint
	venta              domain.Venta
	ventas             []domain.Venta
}

func (f *fakeRepo) Create(_ *domain.Venta) error { return f.createErr }
func (f *fakeRepo) Update(_ *domain.Venta) error { return f.updateErr }
func (f *fakeRepo) Delete(_ int64) error         { return f.deleteErr }
func (f *fakeRepo) GetByID(_ int64) (*domain.Venta, error) {
	if f.getByIDErr != nil {
		return nil, f.getByIDErr
	}
	return &f.venta, nil
}
func (f *fakeRepo) List(_ int, _ int) ([]domain.Venta, int, error) {
	if f.listErr != nil {
		return nil, 0, f.listErr
	}
	return f.ventas, len(f.ventas), nil
}
func (f *fakeRepo) NextCorrelativo(_ string) (uint, error) {
	if f.nextCorrelativoErr != nil {
		return 0, f.nextCorrelativoErr
	}
	if f.nextCorrelativo == 0 {
		return 1, nil
	}
	return f.nextCorrelativo, nil
}

// ---------------------------------------------------------------------------
// Tests Create
// ---------------------------------------------------------------------------

func TestVentaService_Create_Factura_AutoSerie(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{nextCorrelativo: 1}}
	in := input.VentaCreateInput{
		IDUsuario:         1,
		IDTipoComprobante: 2,
		Subtotal:          100,
	}
	out, err := s.Create(in)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out.Serie != "F001" {
		t.Errorf("serie esperada F001, obtenida %s", out.Serie)
	}
	if out.Correlativo != 1 {
		t.Errorf("correlativo esperado 1, obtenido %d", out.Correlativo)
	}
	if out.NumeroComprobante != "F001-000001" {
		t.Errorf("numero_comprobante esperado F001-000001, obtenido %s", out.NumeroComprobante)
	}
	if out.IGV != 18 || out.Total != 118 {
		t.Errorf("IGV o Total incorrectos: IGV=%v, Total=%v", out.IGV, out.Total)
	}
	if out.QRCode == "" {
		t.Error("QRCode no generado")
	}
}

func TestVentaService_Create_Boleta_AutoSerie(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{nextCorrelativo: 5}}
	in := input.VentaCreateInput{
		IDUsuario:         1,
		IDTipoComprobante: 1, // BOLETA
		Subtotal:          200,
	}
	out, err := s.Create(in)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out.Serie != "B001" {
		t.Errorf("serie esperada B001, obtenida %s", out.Serie)
	}
	if out.Correlativo != 5 {
		t.Errorf("correlativo esperado 5, obtenido %d", out.Correlativo)
	}
	if out.NumeroComprobante != "B001-000005" {
		t.Errorf("numero_comprobante esperado B001-000005, obtenido %s", out.NumeroComprobante)
	}
	if out.IGV != 0 || out.Total != 200 {
		t.Errorf("IGV o Total incorrectos para boleta: IGV=%v, Total=%v", out.IGV, out.Total)
	}
}

func TestVentaService_Create_Ticket_AutoSerie(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{nextCorrelativo: 3}}
	in := input.VentaCreateInput{
		IDUsuario:         1,
		IDTipoComprobante: 3, // TICKET
		Subtotal:          150,
	}
	out, err := s.Create(in)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out.Serie != "T001" {
		t.Errorf("serie esperada T001, obtenida %s", out.Serie)
	}
	if out.NumeroComprobante != "T001-000003" {
		t.Errorf("numero_comprobante esperado T001-000003, obtenido %s", out.NumeroComprobante)
	}
	if out.IGV != 0 || out.Total != 150 {
		t.Errorf("IGV o Total incorrectos para ticket: IGV=%v, Total=%v", out.IGV, out.Total)
	}
}

func TestVentaService_Create_TipoComprobanteInvalido(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{}}
	in := input.VentaCreateInput{
		IDUsuario:         1,
		IDTipoComprobante: 99, // invalido
		Subtotal:          100,
	}
	_, err := s.Create(in)
	if err == nil {
		t.Error("debe fallar por tipo de comprobante invalido")
	}
}

func TestVentaService_Create_SubtotalCero(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{}}
	_, err := s.Create(input.VentaCreateInput{IDUsuario: 1, IDTipoComprobante: 1, Subtotal: 0})
	if err == nil {
		t.Error("debe fallar por subtotal invalido")
	}
}

func TestVentaService_Create_UsuarioRequerido(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{}}
	_, err := s.Create(input.VentaCreateInput{IDUsuario: 0, IDTipoComprobante: 1, Subtotal: 100})
	if err == nil {
		t.Error("debe fallar por usuario requerido")
	}
}

func TestVentaService_Create_ErrorCorrelativo(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{nextCorrelativoErr: errors.New("db error")}}
	_, err := s.Create(input.VentaCreateInput{IDUsuario: 1, IDTipoComprobante: 1, Subtotal: 100})
	if err == nil {
		t.Error("debe fallar si no se puede obtener el correlativo")
	}
}

func TestVentaService_Create_ErrorRepo(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{createErr: errors.New("fail")}}
	_, err := s.Create(input.VentaCreateInput{IDUsuario: 1, IDTipoComprobante: 1, Subtotal: 100})
	if err == nil {
		t.Error("debe fallar si el repo falla al crear")
	}
}

// ---------------------------------------------------------------------------
// Tests Update
// ---------------------------------------------------------------------------

func TestVentaService_Update_OK(t *testing.T) {
	venta := domain.Venta{IDVenta: 1, Nota: "old"}
	s := &VentaServiceImpl{repo: &fakeRepo{venta: venta}}
	out, err := s.Update(1, input.VentaUpdateInput{Nota: "nueva", Observaciones: "obs"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out.Nota != "nueva" || out.Observaciones != "obs" {
		t.Error("no se actualizaron los campos correctamente")
	}
}

func TestVentaService_Update_NoEncontrado(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{getByIDErr: errors.New("not found")}}
	_, err := s.Update(99, input.VentaUpdateInput{Nota: "x"})
	if err == nil {
		t.Error("debe fallar si la venta no existe")
	}
}

// ---------------------------------------------------------------------------
// Tests Delete
// ---------------------------------------------------------------------------

func TestVentaService_Delete_OK(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{}}
	if err := s.Delete(1); err != nil {
		t.Errorf("no se esperaba error al eliminar: %v", err)
	}
}

func TestVentaService_Delete_Error(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{deleteErr: errors.New("fail")}}
	if err := s.Delete(1); err == nil {
		t.Error("debe retornar error si falla repo.Delete")
	}
}

// ---------------------------------------------------------------------------
// Tests GetByID
// ---------------------------------------------------------------------------

func TestVentaService_GetByID_OK(t *testing.T) {
	venta := domain.Venta{IDVenta: 1, Serie: "F001", Correlativo: 10}
	s := &VentaServiceImpl{repo: &fakeRepo{venta: venta}}
	out, err := s.GetByID(1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if out.IDVenta != 1 || out.Serie != "F001" || out.Correlativo != 10 {
		t.Error("no se obtuvo la venta correcta")
	}
}

func TestVentaService_GetByID_Error(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{getByIDErr: errors.New("fail")}}
	_, err := s.GetByID(1)
	if err == nil {
		t.Error("debe retornar error si falla repo.GetByID")
	}
}

// ---------------------------------------------------------------------------
// Tests List
// ---------------------------------------------------------------------------

func TestVentaService_List_OK(t *testing.T) {
	ventas := []domain.Venta{{IDVenta: 1}, {IDVenta: 2}}
	s := &VentaServiceImpl{repo: &fakeRepo{ventas: ventas}}
	out, _, err := s.List(1, 15)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("esperado 2 ventas, obtenido %d", len(out))
	}
}

func TestVentaService_List_Error(t *testing.T) {
	s := &VentaServiceImpl{repo: &fakeRepo{listErr: errors.New("fail")}}
	_, _, err := s.List(1, 15)
	if err == nil {
		t.Error("debe retornar error si falla repo.List")
	}
}

// ---------------------------------------------------------------------------
// Test serieFromTipoComprobante
// ---------------------------------------------------------------------------

func TestSerieFromTipoComprobante(t *testing.T) {
	casos := []struct {
		id     int64
		espera string
		error  bool
	}{
		{1, "B001", false},
		{2, "F001", false},
		{3, "T001", false},
		{0, "", true},
		{99, "", true},
	}
	for _, c := range casos {
		serie, err := util.SerieFromTipoComprobante(c.id)
		if c.error && err == nil {
			t.Errorf("id=%d: se esperaba error", c.id)
		}
		if !c.error && serie != c.espera {
			t.Errorf("id=%d: esperada %s, obtenida %s", c.id, c.espera, serie)
		}
	}
}

// ---------------------------------------------------------------------------
// Test campo calculado NumeroComprobante
// ---------------------------------------------------------------------------

func TestToVentaOutput_NumeroComprobante(t *testing.T) {
	casos := []struct {
		serie       string
		correlativo uint
		esperado    string
	}{
		{"B001", 1, "B001-000001"},
		{"B001", 123, "B001-000123"},
		{"F001", 1, "F001-000001"},
		{"F001", 9999, "F001-009999"},
		{"T001", 1000000, "T001-1000000"}, // correlativo grande -> sin truncar
	}
	for _, c := range casos {
		v := &domain.Venta{Serie: c.serie, Correlativo: c.correlativo}
		out := toVentaOutput(v)
		if out.NumeroComprobante != c.esperado {
			t.Errorf("serie=%s correlativo=%d: esperado %s, obtenido %s",
				c.serie, c.correlativo, c.esperado, out.NumeroComprobante)
		}
		if out.Serie != c.serie {
			t.Errorf("serie no mapeada correctamente: %s != %s", out.Serie, c.serie)
		}
		if out.Correlativo != c.correlativo {
			t.Errorf("correlativo no mapeado correctamente: %d != %d", out.Correlativo, c.correlativo)
		}
	}
}
