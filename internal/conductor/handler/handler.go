package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"sistema_venta_pasajes/internal/conductor/input"
	"sistema_venta_pasajes/internal/conductor/service"
	"sistema_venta_pasajes/internal/conductor/util"
	"sistema_venta_pasajes/pkg"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	conductores, err := h.service.List(r.Context())
	if err != nil {
		pkg.WriteError(w, r, pkg.Internal(util.MSG_LIST_ERROR, err))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "Conductores obtenidos correctamente", conductores, map[string]any{"count": len(conductores)})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", util.MSG_NOT_FOUND))
		return
	}
	conductor, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "Conductor obtenido correctamente", conductor, nil)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateConductorInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}
	conductor, err := h.service.Create(r.Context(), in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, "Conductor creado correctamente", conductor, nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", util.MSG_NOT_FOUND))
		return
	}
	var in input.UpdateConductorInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}
	conductor, err := h.service.Update(r.Context(), id, in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "Conductor actualizado correctamente", conductor, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", util.MSG_NOT_FOUND))
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_DELETE_SUCCESS, nil, nil)
}
