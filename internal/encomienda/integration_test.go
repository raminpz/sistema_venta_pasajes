package encomienda

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sistema_venta_pasajes/internal/encomienda/domain"
	"sistema_venta_pasajes/internal/encomienda/handler"
	"sistema_venta_pasajes/internal/encomienda/input"
	"sistema_venta_pasajes/internal/encomienda/service"

	"gorm.io/gorm"
)

// FakeRepo para pruebas
type FakeEncomiendaRepo struct {
	encomiendas map[int64]*domain.Encomienda
	nextID      int64
}

func NewFakeEncomiendaRepo() *FakeEncomiendaRepo {
	return &FakeEncomiendaRepo{
		encomiendas: make(map[int64]*domain.Encomienda),
		nextID:      1,
	}
}

func (r *FakeEncomiendaRepo) Create(e *domain.Encomienda) error {
	e.IDEncomienda = r.nextID
	r.encomiendas[e.IDEncomienda] = e
	r.nextID++
	return nil
}

func (r *FakeEncomiendaRepo) Update(e *domain.Encomienda) error {
	if _, exists := r.encomiendas[e.IDEncomienda]; !exists {
		return gorm.ErrRecordNotFound
	}
	r.encomiendas[e.IDEncomienda] = e
	return nil
}

func (r *FakeEncomiendaRepo) Delete(id int64) error {
	if _, exists := r.encomiendas[id]; !exists {
		return gorm.ErrRecordNotFound
	}
	delete(r.encomiendas, id)
	return nil
}

func (r *FakeEncomiendaRepo) GetByID(id int64) (*domain.Encomienda, error) {
	e, exists := r.encomiendas[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return e, nil
}

func (r *FakeEncomiendaRepo) List(offset, limit int) ([]domain.Encomienda, int, error) {
	var results []domain.Encomienda
	for _, e := range r.encomiendas {
		results = append(results, *e)
	}
	total := len(results)
	if offset < len(results) {
		end := offset + limit
		if end > len(results) {
			end = len(results)
		}
		results = results[offset:end]
	}
	return results, total, nil
}

// Test: Validación de campos opcionales
func TestCreateEncomiendaOptionalFieldsValidation(t *testing.T) {
	tests := []struct {
		name       string
		payload    map[string]interface{}
		expectCode int
		expectMsg  string
	}{
		{
			name: "RemitenteDoc vacío",
			payload: map[string]interface{}{
				"id_venta":            1,
				"id_programacion":     1,
				"costo":               35.0,
				"remitente_nombre":    "Juan Perez",
				"destinatario_nombre": "Maria Lopez",
				"remitente_doc":       "", // Vacío, debe rechazarse
			},
			expectCode: 422,
			expectMsg:  "remitente_doc",
		},
		{
			name: "DestinatarioDoc vacío",
			payload: map[string]interface{}{
				"id_venta":            1,
				"id_programacion":     1,
				"costo":               35.0,
				"remitente_nombre":    "Juan Perez",
				"destinatario_nombre": "Maria Lopez",
				"destinatario_doc":    "", // Vacío, debe rechazarse
			},
			expectCode: 422,
			expectMsg:  "destinatario_doc",
		},
		{
			name: "DestinatarioTel vacío",
			payload: map[string]interface{}{
				"id_venta":            1,
				"id_programacion":     1,
				"costo":               35.0,
				"remitente_nombre":    "Juan Perez",
				"destinatario_nombre": "Maria Lopez",
				"destinatario_tel":    "", // Vacío, debe rechazarse
			},
			expectCode: 422,
			expectMsg:  "destinatario_tel",
		},
		{
			name: "RemitenteDoc con valor válido - Debe pasar",
			payload: map[string]interface{}{
				"id_venta":            1,
				"id_programacion":     1,
				"costo":               35.0,
				"remitente_nombre":    "Juan Perez",
				"remitente_doc":       "12345678",
				"destinatario_nombre": "Maria Lopez",
				"destinatario_tel":    "987654321",
			},
			expectCode: 201,
			expectMsg:  "Encomienda creada",
		},
		{
			name: "RemitenteDoc faltante",
			payload: map[string]interface{}{
				"id_venta":            1,
				"id_programacion":     1,
				"costo":               35.0,
				"remitente_nombre":    "Juan Perez",
				"destinatario_nombre": "Maria Lopez",
				"destinatario_tel":    "987654321",
			},
			expectCode: 422,
			expectMsg:  "remitente_doc",
		},
		{
			name: "DestinatarioTel faltante",
			payload: map[string]interface{}{
				"id_venta":            1,
				"id_programacion":     1,
				"costo":               35.0,
				"remitente_nombre":    "Juan Perez",
				"remitente_doc":       "12345678",
				"destinatario_nombre": "Maria Lopez",
			},
			expectCode: 422,
			expectMsg:  "destinatario_tel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewFakeEncomiendaRepo()
			svc := service.NewEncomiendaService(repo)
			h := handler.NewEncomiendaHandler(svc)

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/v1/encomienda", bytes.NewReader(body))
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, w.Code)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			if tt.expectCode == 422 {
				if details, ok := resp["details"].(map[string]interface{}); !ok || len(details) == 0 {
					t.Errorf("Expected details with field information, got: %v", resp)
				}
				if details, ok := resp["details"].(map[string]interface{}); ok {
					if _, hasField := details[tt.expectMsg]; !hasField {
						t.Errorf("Expected field '%s' in details, got: %v", tt.expectMsg, details)
					}
				}
			}
		})
	}
}

// Test: Creación exitosa
func TestCreateEncomiendaSuccess(t *testing.T) {
	repo := NewFakeEncomiendaRepo()
	svc := service.NewEncomiendaService(repo)
	h := handler.NewEncomiendaHandler(svc)

	payload := map[string]interface{}{
		"id_venta":            1,
		"id_programacion":     1,
		"costo":               35.0,
		"remitente_nombre":    "Juan Perez",
		"remitente_doc":       "12345678",
		"destinatario_nombre": "Maria Lopez",
		"destinatario_tel":    "987654321",
		"estado":              "PENDIENTE",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/encomienda", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["message"] != "Encomienda creada correctamente" {
		t.Errorf("Expected success message, got: %v", resp)
	}
}

// Test: Estado por defecto
func TestCreateEncomiendaDefaultEstado(t *testing.T) {
	repo := NewFakeEncomiendaRepo()
	svc := service.NewEncomiendaService(repo)

	in := input.CreateEncomiendaInput{
		IDVenta:            1,
		IDProgramacion:     1,
		Costo:              35.0,
		RemitenteNombre:    "Juan Perez",
		RemitenteDoc:       "12345678",
		DestinatarioNombre: "Maria Lopez",
		DestinatarioTel:    "987654321",
		// Estado no especificado, debe ser PENDIENTE por defecto
	}

	out, err := svc.Create(in)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if out.Estado != "PENDIENTE" {
		t.Errorf("Expected estado PENDIENTE by default, got: %s", out.Estado)
	}
}
