package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sistema_venta_pasajes/internal/usuario/input"
	"sistema_venta_pasajes/internal/usuario/util"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type fakeUsuarioService struct {
	createFn  func(input.UsuarioCreateInput) (*input.UsuarioOutput, error)
	updateFn  func(int, input.UsuarioUpdateInput) (*input.UsuarioOutput, error)
	deleteFn  func(int) error
	getByIDFn func(int) (*input.UsuarioOutput, error)
	listFn    func(int, int) ([]input.UsuarioOutput, int, error)
}

func (f *fakeUsuarioService) Create(in input.UsuarioCreateInput) (*input.UsuarioOutput, error) {
	return f.createFn(in)
}

func (f *fakeUsuarioService) Update(id int, in input.UsuarioUpdateInput) (*input.UsuarioOutput, error) {
	return f.updateFn(id, in)
}

func (f *fakeUsuarioService) Delete(id int) error {
	return f.deleteFn(id)
}

func (f *fakeUsuarioService) GetByID(id int) (*input.UsuarioOutput, error) {
	return f.getByIDFn(id)
}

func (f *fakeUsuarioService) List(page, size int) ([]input.UsuarioOutput, int, error) {
	return f.listFn(page, size)
}

func TestUsuarioHandler_CrearUsuario(t *testing.T) {
	t.Run("json invalido", func(t *testing.T) {
		h := NewUsuarioHandler(&fakeUsuarioService{})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/usuario", strings.NewReader("{"))
		res := httptest.NewRecorder()

		h.CrearUsuario(res, req)

		if res.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", res.Code)
		}
	})

	t.Run("duplicado", func(t *testing.T) {
		h := NewUsuarioHandler(&fakeUsuarioService{
			createFn: func(input.UsuarioCreateInput) (*input.UsuarioOutput, error) {
				return nil, errors.New(util.MSG_DNI_DUPLICATE)
			},
		})
		body := `{"id_rol":1,"nombre":"Juan","apellidos":"Perez","dni":"12345678","email":"juan@mail.com","password":"1234","telefono":"999999999"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/usuario", strings.NewReader(body))
		res := httptest.NewRecorder()

		h.CrearUsuario(res, req)

		if res.Code != http.StatusConflict {
			t.Fatalf("status esperado 409, obtenido %d", res.Code)
		}
	})

	t.Run("ok", func(t *testing.T) {
		h := NewUsuarioHandler(&fakeUsuarioService{
			createFn: func(in input.UsuarioCreateInput) (*input.UsuarioOutput, error) {
				return &input.UsuarioOutput{IDUsuario: 1, Nombre: in.Nombre}, nil
			},
		})
		body := `{"id_rol":1,"nombre":"Juan","apellidos":"Perez","dni":"12345678","email":"juan@mail.com","password":"1234","telefono":"999999999"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/usuario", strings.NewReader(body))
		res := httptest.NewRecorder()

		h.CrearUsuario(res, req)

		if res.Code != http.StatusCreated {
			t.Fatalf("status esperado 201, obtenido %d", res.Code)
		}
	})
}

func TestUsuarioHandler_Actualizar_Eliminar_Obtener_Listar(t *testing.T) {
	h := NewUsuarioHandler(&fakeUsuarioService{
		updateFn: func(id int, in input.UsuarioUpdateInput) (*input.UsuarioOutput, error) {
			return &input.UsuarioOutput{IDUsuario: id, Nombre: in.Nombre}, nil
		},
		deleteFn: func(int) error { return nil },
		getByIDFn: func(id int) (*input.UsuarioOutput, error) {
			return &input.UsuarioOutput{IDUsuario: id, Nombre: "Juan"}, nil
		},
		listFn: func(page, size int) ([]input.UsuarioOutput, int, error) {
			return []input.UsuarioOutput{{IDUsuario: 1, Nombre: "Juan"}}, 1, nil
		},
	})

	t.Run("actualizar id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/usuario/abc", strings.NewReader(`{"nombre":"X"}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		res := httptest.NewRecorder()
		h.ActualizarUsuario(res, req)
		if res.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", res.Code)
		}
	})

	t.Run("actualizar ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/usuario/1", strings.NewReader(`{"nombre":"Luis","apellidos":"Diaz","email":"l@mail.com","telefono":"999999999","estado":"ACTIVO"}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		res := httptest.NewRecorder()
		h.ActualizarUsuario(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("status esperado 200, obtenido %d", res.Code)
		}
	})

	t.Run("eliminar id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/usuario/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		res := httptest.NewRecorder()
		h.EliminarUsuario(res, req)
		if res.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", res.Code)
		}
	})

	t.Run("obtener id invalido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/usuario/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		res := httptest.NewRecorder()
		h.ObtenerUsuarioPorID(res, req)
		if res.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", res.Code)
		}
	})

	t.Run("obtener ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/usuario/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		res := httptest.NewRecorder()
		h.ObtenerUsuarioPorID(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("status esperado 200, obtenido %d", res.Code)
		}
	})

	t.Run("listar paginacion invalida", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/usuario?page=abc", nil)
		res := httptest.NewRecorder()
		h.ListarUsuarios(res, req)
		if res.Code != http.StatusBadRequest {
			t.Fatalf("status esperado 400, obtenido %d", res.Code)
		}
	})

	t.Run("listar ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/usuario?page=1&size=10", nil)
		res := httptest.NewRecorder()
		h.ListarUsuarios(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("status esperado 200, obtenido %d", res.Code)
		}
	})
}

func TestMapUsuarioServiceError(t *testing.T) {
	if got := mapUsuarioServiceError(nil, "x", "y"); got != nil {
		t.Fatal("se esperaba nil cuando err es nil")
	}

	err := mapUsuarioServiceError(errors.New(util.MSG_USER_NOT_FOUND), "c", "f")
	if err == nil {
		t.Fatal("se esperaba error mapeado")
	}
}
