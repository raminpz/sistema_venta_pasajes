package handler

import (
	"encoding/json"
	"net/http"
	"sistema_venta_pasajes/internal/pasajero/input"
	"sistema_venta_pasajes/internal/pasajero/service"
	"sistema_venta_pasajes/internal/pasajero/util"
	"sistema_venta_pasajes/pkg"
	"strconv"

	"github.com/gorilla/mux"
)

type PasajeroHandler struct {
	service service.PasajeroService
}

func NewPasajeroHandler(s service.PasajeroService) *PasajeroHandler {
	return &PasajeroHandler{service: s}
}

func (h *PasajeroHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreatePasajeroInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("json_invalid", util.MSG_JSON_INVALID).WithCause(err))
		return
	}
	out, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("create_error", util.MSG_CREATE_ERROR).WithCause(err).WithDetails(err.Error()))
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_CREATE_SUCCESS, out, nil)
}

func (h *PasajeroHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		pkg.WriteError(w, r, pkg.BadRequest("missing_id", util.MSG_MISSING_ID))
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", util.MSG_INVALID_ID).WithCause(err))
		return
	}
	var in input.UpdatePasajeroInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("json_invalid", util.MSG_JSON_INVALID).WithCause(err))
		return
	}
	out, err := h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("update_error", util.MSG_UPDATE_ERROR).WithCause(err).WithDetails(err.Error()))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_UPDATE_SUCCESS, out, nil)
}

func (h *PasajeroHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		pkg.WriteError(w, r, pkg.BadRequest("missing_id", util.MSG_MISSING_ID))
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", util.MSG_INVALID_ID).WithCause(err))
		return
	}
	err = h.service.Delete(id)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("delete_error", util.MSG_DELETE_ERROR).WithCause(err).WithDetails(err.Error()))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_DELETE_SUCCESS, nil, nil)
}

func (h *PasajeroHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		pkg.WriteError(w, r, pkg.BadRequest("missing_id", util.MSG_MISSING_ID))
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", util.MSG_INVALID_ID).WithCause(err))
		return
	}
	out, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, pkg.NotFound("not_found", util.MSG_NOT_FOUND).WithCause(err).WithDetails(err.Error()))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_FOUND_SUCCESS, out, nil)
}

func (h *PasajeroHandler) List(w http.ResponseWriter, r *http.Request) {
	page := 1
	size := 10
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if s := r.URL.Query().Get("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			size = v
		}
	}
	   out, meta, err := h.service.List(page, size)
	   if err != nil {
			   pkg.WriteError(w, r, pkg.Internal(util.MSG_LIST_ERROR).WithCause(err).WithDetails(err.Error()))
			   return
	   }
	   // Si no hay resultados, devolver array vacío
	   if out == nil || len(out) == 0 {
			   out = []input.PasajeroOutput{}
	   }
	   pkg.WriteSuccess(w, http.StatusOK, util.MSG_LIST_SUCCESS, out, meta)
}

func (h *PasajeroHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		pkg.WriteError(w, r, pkg.BadRequest("missing_query", util.MSG_MISSING_QUERY))
		return
	}
	out, err := h.service.Search(query)
	if err != nil {
		pkg.WriteError(w, r, pkg.Internal(util.MSG_SEARCH_ERROR).WithCause(err).WithDetails(err.Error()))
		return
	}
	if out == nil || len(out) == 0 {
		out = []input.PasajeroOutput{}
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_SEARCH_SUCCESS, out, nil)
}
