package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sistema_venta_pasajes/internal/venta/input"
	"sistema_venta_pasajes/internal/venta/service"
	"sistema_venta_pasajes/internal/venta/util"
	"sistema_venta_pasajes/pkg"

	"github.com/gorilla/mux"
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
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_BODY, util.MSG_VENTA_VALIDATION_ERROR).WithCause(err))
		return
	}
	venta, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusBadRequest, util.ERR_CODE_CREATE, err.Error()))
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_VENTA_CREATED, venta, nil)
}

func (h *VentaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_VENTA_NOT_FOUND))
		return
	}
	var in input.VentaUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_BODY, util.MSG_VENTA_VALIDATION_ERROR).WithCause(err))
		return
	}
	venta, err := h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusBadRequest, util.ERR_CODE_UPDATE, err.Error()))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_VENTA_UPDATED, venta, nil)
}

func (h *VentaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_VENTA_NOT_FOUND))
		return
	}
	if err := h.service.Delete(id); err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_DELETE, util.MSG_VENTA_DELETE_ERROR).WithCause(err))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_VENTA_DELETED, nil, nil)
}

func (h *VentaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_VENTA_NOT_FOUND))
		return
	}
	venta, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusNotFound, util.ERR_CODE_NOT_FOUND, util.MSG_VENTA_NOT_FOUND).WithCause(err))
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_VENTA_GET, venta, nil)
}

func (h *VentaHandler) List(w http.ResponseWriter, r *http.Request) {
	page, size, err := util.ParsePaginationParams(r)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_PAGE, util.MSG_VENTA_PAGINATION_INVALID).WithCause(err))
		return
	}

	ventas, total, err := h.service.List(page, size)
	if err != nil {
		pkg.WriteError(w, r, pkg.NewAppError(http.StatusInternalServerError, util.ERR_CODE_LIST, util.MSG_VENTA_LIST_ERROR).WithCause(err))
		return
	}
	if ventas == nil {
		ventas = []input.VentaOutput{}
	}
	_, _, meta := pkg.Paginate(page, size, total)
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_VENTA_LIST, ventas, meta)
}
