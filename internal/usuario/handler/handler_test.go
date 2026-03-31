package handler

import (
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/usuario/repository"
	"sistema_venta_pasajes/internal/usuario/service"
	"strings"
	"testing"
)

func TestCrearUsuario(t *testing.T) {
	repo := repository.NewUsuarioRepositoryMock()
	serv := service.NewUsuarioService(repo)
	h := NewUsuarioHandler(serv)

	body := `{"id_rol":1,"nombre":"Juan","apellidos":"Perez","dni":"12345678","email":"juan@mail.com","password":"1234","telefono":"999999999"}`
	req := httptest.NewRequest("POST", "/api/v1/usuario", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CrearUsuario(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("esperado status 201, obtenido %d", resp.StatusCode)
	}
}
