package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"sistema_venta_pasajes/internal/venta/input"
	"sistema_venta_pasajes/internal/venta/service"
	"sistema_venta_pasajes/internal/venta/util"
	"sistema_venta_pasajes/pkg"
)

type VentaHandler struct {
	service service.VentaService
}

func NewVentaHandler(s service.VentaService) *VentaHandler {
	return &VentaHandler{service: s}
}

func (h *VentaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.VentaCreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ErrCodeInvalidBody, util.MsgVentaErrorValidacion).WithCause(err))
		return
	}
	venta, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusBadRequest, util.ErrCodeCreateError, err.Error()))
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MsgVentaCreada, venta, nil)
}

func (h *VentaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ErrCodeInvalidID, util.MsgVentaNoEncontrada))
		return
	}
	var in input.VentaUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ErrCodeInvalidBody, util.MsgVentaErrorValidacion).WithCause(err))
		return
	}
	venta, err := h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusBadRequest, util.ErrCodeUpdateError, err.Error()))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MsgVentaActualizada, venta, nil)
}

func (h *VentaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ErrCodeInvalidID, util.MsgVentaNoEncontrada))
		return
	}
	if err := h.service.Delete(id); err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusInternalServerError, util.ErrCodeDeleteError, util.MsgVentaErrorEliminar).WithCause(err))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MsgVentaEliminada, nil, nil)
}

func (h *VentaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ErrCodeInvalidID, util.MsgVentaNoEncontrada))
		return
	}
	venta, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusNotFound, util.ErrCodeNotFound, util.MsgVentaNoEncontrada).WithCause(err))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MsgVentaObtener, venta, nil)
}

func (h *VentaHandler) List(w http.ResponseWriter, r *http.Request) {
	page, size, err := util.ParsePaginationParams(r)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ErrCodeInvalidPage, util.MsgVentaPaginacionInvalida).WithCause(err))
		return
	}

	ventas, total, err := h.service.List(page, size)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusInternalServerError, util.ErrCodeListError, util.MsgVentaErrorListar).WithCause(err))
		return
	}
	if ventas == nil {
		ventas = []input.VentaOutput{}
	}
	_, _, meta := pkg.Paginate(page, size, total)
	pkg.WriteSuccess(w, http.StatusOK, util.MsgVentaListada, ventas, meta)
}
