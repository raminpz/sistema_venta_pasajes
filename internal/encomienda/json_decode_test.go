package encomienda

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"sistema_venta_pasajes/internal/encomienda/handler"
	"sistema_venta_pasajes/internal/encomienda/service"
)

// Test: JSON decode errors
func TestCreateEncomiendaJSONDecodeErrors(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		expectCode    int
		expectErrCode string
	}{
		{
			name:          "JSON inválido",
			body:          `{invalid json}`,
			expectCode:    400,
			expectErrCode: "invalid_json",
		},
		{
			name:          "Tipo inválido para campo numérico",
			body:          `{"id_venta": "not a number", "id_programacion": 1, "costo": 35.0}`,
			expectCode:    400,
			expectErrCode: "invalid_json_type",
		},
		{
			name:          "Body vacío",
			body:          ``,
			expectCode:    400,
			expectErrCode: "empty_body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewFakeEncomiendaRepo()
			svc := service.NewEncomiendaService(repo)
			h := handler.NewEncomiendaHandler(svc)

			req := httptest.NewRequest("POST", "/api/v1/encomienda", bytes.NewReader([]byte(tt.body)))
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, w.Code)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			if err, ok := resp["error"].(string); ok && err != tt.expectErrCode {
				t.Errorf("Expected error code '%s', got '%s'", tt.expectErrCode, err)
			}
		})
	}
}
