package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"sistema_venta_pasajes/internal/asiento/domain"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../config/.env")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	params := os.Getenv("DB_PARAMS")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, pass, host, port, name, params)
	fmt.Printf("DB_HOST=%s\nDB_PORT=%s\nDB_NAME=%s\nDB_USER=%s\nDB_PASS=%s\nDB_PARAMS=%s\n", host, port, name, user, pass, params)
	fmt.Printf("DSN: %s\n", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	testDB = db
	os.Exit(m.Run())
}

func cleanTable(t *testing.T) {
	// Borrar primero de las tablas hijas para respetar claves foráneas
	err := testDB.Exec("DELETE FROM DETALLE_PASAJE").Error
	if err != nil {
		t.Fatalf("failed to clean DETALLE_PASAJE: %v", err)
	}
	err = testDB.Exec("DELETE FROM ASIENTO").Error
	if err != nil {
		t.Fatalf("failed to clean ASIENTO: %v", err)
	}
}

func TestAsientoRepository_Create(t *testing.T) {
	cleanTable(t)
	repo := NewAsientoRepository(testDB)
	a := &domain.Asiento{IDVehiculo: 1, NumeroAsiento: "A1", Estado: "ACTIVO"}
	err := repo.Create(a)
	if err != nil || a.IDAsiento == 0 {
		t.Errorf("unexpected result: err=%v, asiento=%+v", err, a)
	}
	if a.Estado != "ACTIVO" {
		t.Errorf("expected Estado=ACTIVO, got %s", a.Estado)
	}
}

func TestAsientoRepository_GetByID(t *testing.T) {
	cleanTable(t)
	repo := NewAsientoRepository(testDB)
	a := &domain.Asiento{IDVehiculo: 2, NumeroAsiento: "B2", Estado: "RESERVADO"}
	repo.Create(a)
	got, err := repo.GetByID(int64(a.IDAsiento))
	if err != nil || got.IDAsiento != a.IDAsiento {
		t.Errorf("unexpected result: err=%v, asiento=%+v", err, got)
	}
	if got.Estado != "RESERVADO" {
		t.Errorf("expected Estado=RESERVADO, got %s", got.Estado)
	}
}

func TestAsientoRepository_ListByVehiculo(t *testing.T) {
	cleanTable(t)
	repo := NewAsientoRepository(testDB)
	repo.Create(&domain.Asiento{IDVehiculo: 2, NumeroAsiento: "A1", Estado: "ACTIVO"})
	repo.Create(&domain.Asiento{IDVehiculo: 2, NumeroAsiento: "A2", Estado: "OCUPADO"})
	asientos, err := repo.ListByVehiculo(2)
	if err != nil || len(asientos) != 2 {
		t.Errorf("unexpected result: err=%v, asientos=%+v", err, asientos)
	}
}

func TestAsientoRepository_Update(t *testing.T) {
	cleanTable(t)
	repo := NewAsientoRepository(testDB)
	a := &domain.Asiento{IDVehiculo: 1, NumeroAsiento: "C3", Estado: "ACTIVO"}
	repo.Create(a)
	a.NumeroAsiento = "C4"
	a.Estado = "OCUPADO"
	err := repo.Update(a)
	if err != nil {
		t.Errorf("unexpected result: err=%v", err)
	}
	got, _ := repo.GetByID(int64(a.IDAsiento))
	if got.NumeroAsiento != "C4" {
		t.Errorf("update failed: got=%+v", got)
	}
	if got.Estado != "OCUPADO" {
		t.Errorf("expected Estado=OCUPADO, got %s", got.Estado)
	}
}

func TestAsientoRepository_Delete(t *testing.T) {
	cleanTable(t)
	repo := NewAsientoRepository(testDB)
	a := &domain.Asiento{IDVehiculo: 1, NumeroAsiento: "D1"}
	repo.Create(a)
	err := repo.Delete(int64(a.IDAsiento))
	if err != nil {
		t.Errorf("unexpected result: err=%v", err)
	}
	_, err = repo.GetByID(int64(a.IDAsiento))
	if err == nil {
		t.Errorf("expected error for deleted record, got nil")
	}
}

func TestAsientoRepository_CambiarEstado(t *testing.T) {
	cleanTable(t)
	repo := NewAsientoRepository(testDB)
	a := &domain.Asiento{IDVehiculo: 1, NumeroAsiento: "F1", Estado: "ACTIVO"}
	err := repo.Create(a)
	if err != nil {
		t.Fatalf("no se pudo crear asiento: %v", err)
	}
	// Cambiar a RESERVADO
	err = repo.CambiarEstado(int64(a.IDAsiento), "RESERVADO")
	if err != nil {
		t.Errorf("error al cambiar a RESERVADO: %v", err)
	}
	got, _ := repo.GetByID(int64(a.IDAsiento))
	if got.Estado != "RESERVADO" {
		t.Errorf("expected Estado=RESERVADO, got %s", got.Estado)
	}
	// Cambiar a OCUPADO
	err = repo.CambiarEstado(int64(a.IDAsiento), "OCUPADO")
	if err != nil {
		t.Errorf("error al cambiar a OCUPADO: %v", err)
	}
	got, _ = repo.GetByID(int64(a.IDAsiento))
	if got.Estado != "OCUPADO" {
		t.Errorf("expected Estado=OCUPADO, got %s", got.Estado)
	}
	// Cambiar a ACTIVO
	err = repo.CambiarEstado(int64(a.IDAsiento), "ACTIVO")
	if err != nil {
		t.Errorf("error al cambiar a ACTIVO: %v", err)
	}
	got, _ = repo.GetByID(int64(a.IDAsiento))
	if got.Estado != "ACTIVO" {
		t.Errorf("expected Estado=ACTIVO, got %s", got.Estado)
	}
}
