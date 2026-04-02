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
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_BODY, util.MSG_INVALID_DATA).WithCause(err))
		return
	}
	usuario, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, mapUsuarioServiceError(err, util.ERR_CODE_CREATE, util.MSG_CREATE_ERROR))
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_USER_CREATED, usuario, nil)
}

func (h *UsuarioHandler) ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_INVALID_ID).WithCause(err))
		return
	}
	var in input.UsuarioUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_BODY, util.MSG_INVALID_DATA).WithCause(err))
		return
	}
	_, err = h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, mapUsuarioServiceError(err, util.ERR_CODE_UPDATE, util.MSG_UPDATE_ERROR))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_USER_UPDATED, nil, nil)
}

func (h *UsuarioHandler) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_INVALID_ID).WithCause(err))
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		pkg.WriteError(w, r, mapUsuarioServiceError(err, util.ERR_CODE_DELETE, util.MSG_DELETE_ERROR))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_USER_DELETED, nil, nil)
}

func (h *UsuarioHandler) ObtenerUsuarioPorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_INVALID_ID).WithCause(err))
		return
	}
	usuario, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, mapUsuarioServiceError(err, util.ERR_CODE_NOT_FOUND, util.MSG_USER_NOT_FOUND))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_GET_OK, usuario, nil)
}

func (h *UsuarioHandler) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	page, size, err := util.ParsePaginationParams(r)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_PAGE, util.MSG_INVALID_DATA).WithCause(err))
		return
	}
	usuarios, total, err := h.service.List(page, size)
	if err != nil {
		pkg.WriteError(w, r, mapUsuarioServiceError(err, util.ERR_CODE_LIST, util.MSG_LIST_ERROR))
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
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_USER_LIST, usuarios, meta)
}

func mapUsuarioServiceError(err error, code, fallback string) error {
	if err == nil {
		return nil
	}
	switch err.Error() {
	case util.MSG_INVALID_DATA:
		return pkg.BadRequest(code, util.MSG_INVALID_DATA).WithCause(err)
	case util.MSG_EMAIL_DUPLICATE, util.MSG_DNI_DUPLICATE:
		return pkg.Conflict("duplicate_resource", err.Error()).WithCause(err)
	case util.MSG_USER_NOT_FOUND:
		return pkg.NotFound(util.ERR_CODE_NOT_FOUND, util.MSG_USER_NOT_FOUND).WithCause(err)
	default:
		return pkg.NewAppError(http.StatusInternalServerError, code, fallback).WithCause(err)
	}
}
