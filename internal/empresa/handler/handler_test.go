package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sistema_venta_pasajes/internal/empresa/input"
	"sistema_venta_pasajes/internal/empresa/repository"
	"sistema_venta_pasajes/internal/empresa/service"
	"sistema_venta_pasajes/pkg"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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

func setupHandler() *EmpresaHandler {
	db := setupTestDB()
	repo := repository.NewEmpresaRepository(db)
	svc := service.NewEmpresaService(repo)
	return NewEmpresaHandler(svc)
}

// Reemplaza setupHandler para usar el mock
func setupMockHandler(svc *EmpresaServiceMock) *EmpresaHandler {
	return NewEmpresaHandler(svc)
}

func TestEmpresaHandler_Create_And_GetByID(t *testing.T) {
	svc := &EmpresaServiceMock{}
	h := setupMockHandler(svc)
	in := input.CreateEmpresaInput{
		RUC:           "33333333333",
		RazonSocial:   "Empresa Handler",
		Telefono:      "987654333",
		FechaCreacion: time.Now().Format("2006-01-02"),
	}
	out := input.EmpresaOutput{
		IDEmpresa:     1,
		RUC:           in.RUC,
		RazonSocial:   in.RazonSocial,
		Telefono:      in.Telefono,
		FechaCreacion: in.FechaCreacion,
	}
	svc.On("Create", in).Return(out, nil)
	body, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/empresa", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Create(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Obtener por ID
	svc.On("GetByID", int64(1)).Return(out, nil)
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/empresa/1", nil)
	getW := httptest.NewRecorder()
	vars := map[string]string{"id": "1"}
	getReq = mux.SetURLVars(getReq, vars)
	h.GetByID(getW, getReq)
	assert.Equal(t, http.StatusOK, getW.Code)
}

func TestEmpresaHandler_Create_Invalid(t *testing.T) {
	svc := &EmpresaServiceMock{}
	h := setupMockHandler(svc)
	in := input.CreateEmpresaInput{
		RUC:           "abc",
		RazonSocial:   "",
		Telefono:      "123",
		FechaCreacion: "2026-03-27",
	}
	svc.On("Create", in).Return(input.EmpresaOutput{}, pkg.Validation("error de validación", nil))
	body, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/empresa", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Create(w, req)
	assert.Contains(t, []int{http.StatusBadRequest, http.StatusUnprocessableEntity}, w.Code)
}

func TestEmpresaHandler_List(t *testing.T) {
	svc := &EmpresaServiceMock{}
	h := setupMockHandler(svc)
	out1 := input.EmpresaOutput{IDEmpresa: 1, RUC: "44444444440", RazonSocial: "Empresa L0", Telefono: "987654321", FechaCreacion: time.Now().Format("2006-01-02")}
	out2 := input.EmpresaOutput{IDEmpresa: 2, RUC: "44444444441", RazonSocial: "Empresa L1", Telefono: "987654322", FechaCreacion: time.Now().Format("2006-01-02")}
	svc.On("List").Return([]input.EmpresaOutput{out1, out2}, nil)
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/empresas", nil)
	listW := httptest.NewRecorder()
	h.List(listW, listReq)
	assert.Equal(t, http.StatusOK, listW.Code)
	   var resp struct {
			   Data []map[string]interface{} `json:"data"`
	   }
	   _ = json.Unmarshal(listW.Body.Bytes(), &resp)
	   assert.GreaterOrEqual(t, len(resp.Data), 2, "Debe listar al menos dos empresas")
}
