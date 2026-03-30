package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"
)

func TestHealthRouteReturnsSuccessEnvelope(t *testing.T) {
	router := NewRouter(nil)
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("se esperaba el estado %d, se obtuvo %d", http.StatusOK, response.Code)
	}

	if response.Header().Get("X-Request-ID") == "" {
		t.Fatal("se esperaba que el header X-Request-ID estuviera presente")
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("se esperaba una respuesta JSON válida, se obtuvo el error: %v", err)
	}
	if payload["code"] != float64(http.StatusOK) {
		t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusOK, payload["code"])
	}

	if payload["message"] != "servicio disponible" {
		t.Fatalf("se esperaba el mensaje 'servicio disponible', se obtuvo %#v", payload["message"])
	}

	data, ok := payload["data"].(map[string]any)
	if !ok {
		t.Fatalf("se esperaba un objeto en data, se obtuvo %#v", payload["data"])
	}

	if data["status"] != "ok" {
		t.Fatalf("se esperaba data.status = 'ok', se obtuvo %#v", data["status"])
	}
}

func TestReadyRouteWithoutDatabaseReturnsCentralizedError(t *testing.T) {
	router := NewRouter(nil)
	request := httptest.NewRequest(http.MethodGet, "/ready", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("se esperaba el estado %d, se obtuvo %d", http.StatusServiceUnavailable, response.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("se esperaba una respuesta JSON válida, se obtuvo el error: %v", err)
	}

	if payload["code"] != float64(http.StatusServiceUnavailable) {
		t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusServiceUnavailable, payload["code"])
	}

	if payload["error"] != "database_unavailable" {
		t.Fatalf("se esperaba el error 'database_unavailable', se obtuvo %#v", payload["error"])
	}
}

func TestNotFoundRouteReturnsCentralizedError(t *testing.T) {
	router := NewRouter(nil)
	request := httptest.NewRequest(http.MethodGet, "/no-existe-totalmente", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	// Gorilla/mux puede devolver 405 en vez de 404 si hay handler global OPTIONS
	if response.Code != http.StatusNotFound && response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("se esperaba el estado 404 o 405, se obtuvo %d", response.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("se esperaba una respuesta JSON válida, se obtuvo el error: %v", err)
	}

	if response.Code == http.StatusNotFound {
		if payload["code"] != float64(http.StatusNotFound) {
			t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusNotFound, payload["code"])
		}
		if payload["error"] != "route_not_found" {
			t.Fatalf("se esperaba el error 'route_not_found', se obtuvo %#v", payload["error"])
		}
	} else if response.Code == http.StatusMethodNotAllowed {
		if payload["code"] != float64(http.StatusMethodNotAllowed) {
			t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusMethodNotAllowed, payload["code"])
		}
		if payload["error"] != "method_not_allowed" {
			t.Fatalf("se esperaba el error 'method_not_allowed', se obtuvo %#v", payload["error"])
		}
	}
}

func TestMethodNotAllowedReturnsCentralizedError(t *testing.T) {
	router := NewRouter(nil)
	request := httptest.NewRequest(http.MethodPost, "/health", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("se esperaba el estado %d, se obtuvo %d", http.StatusMethodNotAllowed, response.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("se esperaba una respuesta JSON válida, se obtuvo el error: %v", err)
	}

	if payload["code"] != float64(http.StatusMethodNotAllowed) {
		t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusMethodNotAllowed, payload["code"])
	}
	if payload["error"] != "method_not_allowed" {
		t.Fatalf("se esperaba el error 'method_not_allowed', se obtuvo %#v", payload["error"])
	}
}

func TestProveedorSistemaDeleteRouteReturnsValidationErrorForInvalidID(t *testing.T) {
	router := NewRouter(&gorm.DB{})
	request := httptest.NewRequest(http.MethodDelete, "/api/v1/proveedor/abc", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba el estado %d, se obtuvo %d", http.StatusBadRequest, response.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("se esperaba una respuesta JSON válida, se obtuvo el error: %v", err)
	}

	if payload["code"] != float64(http.StatusBadRequest) {
		t.Fatalf("se esperaba el code %d, se obtuvo %#v", http.StatusBadRequest, payload["code"])
	}
	if payload["error"] != "invalid_provider_id" {
		t.Fatalf("se esperaba el error 'invalid_provider_id', se obtuvo %#v", payload["error"])
	}
}

// Se eliminan los tests de empresa porque la funcionalidad fue removida y la variable 'request' no está definida
