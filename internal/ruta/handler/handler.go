package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sistema_venta_pasajes/internal/ruta/input"
	"sistema_venta_pasajes/internal/ruta/service"
	"sistema_venta_pasajes/pkg"

	"github.com/gorilla/mux"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	rutas, err := h.service.List(r.Context())
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "Rutas obtenidas correctamente", rutas, map[string]any{"count": len(rutas)})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", "ID inválido"))
		return
	}
	ruta, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "Ruta obtenida correctamente", ruta, nil)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateRutaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}
	ruta, err := h.service.Create(r.Context(), in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, "Ruta creada correctamente", ruta, nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", "ID inválido"))
		return
	}
	var in input.UpdateRutaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}
	ruta, err := h.service.Update(r.Context(), id, in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "Ruta actualizada correctamente", ruta, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest("invalid_id", "ID inválido"))
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "Ruta eliminada correctamente", nil, nil)
}
