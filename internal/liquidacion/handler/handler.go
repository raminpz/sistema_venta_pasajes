package handler

import (
	"encoding/json"
	"net/http"
	"sistema_venta_pasajes/internal/liquidacion/input"
	"sistema_venta_pasajes/internal/liquidacion/service"
	"sistema_venta_pasajes/internal/liquidacion/util"
	"sistema_venta_pasajes/pkg"
	"strconv"

	"github.com/gorilla/mux"
)

// LiquidacionHandler expone los endpoints HTTP del módulo.
type LiquidacionHandler struct {
	service service.LiquidacionService
}

// NewLiquidacionHandler crea el handler.
func NewLiquidacionHandler(s service.LiquidacionService) *LiquidacionHandler {
	return &LiquidacionHandler{service: s}
}

// Generar genera una liquidación de caja para una programación.
func (h *LiquidacionHandler) Generar(w http.ResponseWriter, r *http.Request) {
	var in input.GenerarLiquidacionInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.ERR_INVALID_JSON).WithCause(err))
		return
	}
	liq, err := h.service.Generar(in)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_CREATED, liq, nil)
}

// ActualizarEstado actualiza el estado de la liquidación (PENDIENTE → ENTREGADO).
func (h *LiquidacionHandler) ActualizarEstado(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	var in input.ActualizarEstadoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.ERR_INVALID_JSON).WithCause(err))
		return
	}
	liq, err := h.service.ActualizarEstado(id, in)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_UPDATED, liq, nil)
}

// Delete elimina una liquidación por ID.
func (h *LiquidacionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	if err := h.service.Delete(id); err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_DELETED, nil, nil)
}

// GetByID obtiene una liquidación por su ID.
func (h *LiquidacionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	liq, err := h.service.GetByID(id)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_GET, liq, nil)
}

// List lista las liquidaciones con paginación.
func (h *LiquidacionHandler) List(w http.ResponseWriter, r *http.Request) {
	page, size, err := pkg.ParsePaginationParams(r, 1, 15)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_pagination", util.ERR_INVALID_PAGE).WithCause(err))
		return
	}
	liqs, total, err := h.service.List(page, size)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	_, _, meta := pkg.Paginate(page, size, total)
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_LIST, liqs, meta)
}

// ObtenerResumenCaja previsualiza la caja de una programación sin persistir.
func (h *LiquidacionHandler) ObtenerResumenCaja(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id_programacion"], 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	resumen, err := h.service.ObtenerResumenCaja(id)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_RESUMEN_CAJA, resumen, nil)
}
