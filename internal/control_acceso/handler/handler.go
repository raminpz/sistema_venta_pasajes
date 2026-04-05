package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"sistema_venta_pasajes/internal/control_acceso/input"
	"sistema_venta_pasajes/internal/control_acceso/service"
	"sistema_venta_pasajes/internal/control_acceso/util"
	"sistema_venta_pasajes/pkg"
)

type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// GetStatus es público: devuelve el estado actual del sistema al frontend.
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.svc.GetStatus()
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_ACCESO_STATUS, status, nil)
}

// GetLatest devuelve los detalles completos del control de acceso (solo proveedor).
func (h *Handler) GetLatest(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.GetLatest()
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_ACCESO_DETALLES, out, nil)
}

// Create crea y activa un nuevo control de acceso (solo proveedor).
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.ActivarControlAccesoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}
	out, err := h.svc.Create(in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_ACCESO_CREADO, out, nil)
}

// Activar reactiva el sistema (solo proveedor).
func (h *Handler) Activar(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	if err := h.svc.Activar(id); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_ACCESO_ACTIVADO, nil, nil)
}

// Bloquear bloquea el sistema manualmente (solo proveedor).
func (h *Handler) Bloquear(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	if err := h.svc.Bloquear(id); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_ACCESO_BLOQUEADO, nil, nil)
}

// Renovar renueva el control de acceso con una nueva fecha de expiración (solo proveedor).
func (h *Handler) Renovar(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	var in input.RenovarControlAccesoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}
	out, err := h.svc.Renovar(id, in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_ACCESO_RENOVADO, out, nil)
}

func parseID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, pkg.BadRequest("id_invalido", "El ID de control de acceso es inválido.")
	}
	return id, nil
}
