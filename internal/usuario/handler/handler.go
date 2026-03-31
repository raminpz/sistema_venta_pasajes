package handler

import (
	"encoding/json"
	"net/http"
	"sistema_venta_pasajes/internal/usuario/input"
	"sistema_venta_pasajes/internal/usuario/service"
	"sistema_venta_pasajes/internal/usuario/util"
	"sistema_venta_pasajes/pkg"
	"strconv"

	"github.com/gorilla/mux"
)

type UsuarioHandler struct {
	service service.UsuarioService
}

func NewUsuarioHandler(s service.UsuarioService) *UsuarioHandler {
	return &UsuarioHandler{service: s}
}

func (h *UsuarioHandler) CrearUsuario(w http.ResponseWriter, r *http.Request) {
	var in input.UsuarioCreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_body", "Datos inválidos").WithCause(err))
		return
	}
	usuario, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusInternalServerError, "create_error", "Error al crear usuario").WithCause(err))
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, "Usuario creado", usuario, nil)
}

func (h *UsuarioHandler) ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", "ID inválido").WithCause(err))
		return
	}
	var in input.UsuarioUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_body", "Datos inválidos").WithCause(err))
		return
	}
	_, err = h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusInternalServerError, "update_error", "Error al actualizar usuario").WithCause(err))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "Usuario actualizado", nil, nil)
}

func (h *UsuarioHandler) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", "ID inválido").WithCause(err))
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusInternalServerError, "delete_error", "Error al eliminar usuario").WithCause(err))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "Usuario eliminado", nil, nil)
}

func (h *UsuarioHandler) ObtenerUsuarioPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", "ID inválido").WithCause(err))
		return
	}
	usuario, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, pkg.NotFound("usuario_not_found", "Usuario no encontrado").WithCause(err))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "OK", usuario, nil)
}

func (h *UsuarioHandler) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	page, size, err := util.ParsePaginationParams(r)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_pagination", err.Error()))
		return
	}
	usuarios, total, err := h.service.List(page, size)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusInternalServerError, "list_error", "Error al listar usuarios").WithCause(err))
		return
	}
	if usuarios == nil {
		usuarios = []input.UsuarioOutput{}
	}
	meta := pkg.PaginationMeta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: (total + size - 1) / size,
	}
	pkg.WriteSuccess(w, http.StatusOK, "OK", usuarios, meta)
}
