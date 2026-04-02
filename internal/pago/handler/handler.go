package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sistema_venta_pasajes/internal/pago/input"
	"sistema_venta_pasajes/internal/pago/service"
	"sistema_venta_pasajes/internal/pago/util"
	"sistema_venta_pasajes/pkg"

	"github.com/gorilla/mux"
)

type PagoHandler struct {
	service service.PagoService
}

func NewPagoHandler(s service.PagoService) *PagoHandler {
	return &PagoHandler{service: s}
}

func (h *PagoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreatePagoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}
	out, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_PAGO_CREATED, out, nil)
}

func (h *PagoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PAGO_INVALID_ID))
		return
	}
	var in input.UpdatePagoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}
	out, err := h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_PAGO_UPDATED, out, nil)
}

func (h *PagoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PAGO_INVALID_ID))
		return
	}
	if err := h.service.Delete(id); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_PAGO_DELETED, nil, nil)
}

func (h *PagoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PAGO_INVALID_ID))
		return
	}
	out, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_PAGO_GET, out, nil)
}

func (h *PagoHandler) List(w http.ResponseWriter, r *http.Request) {
	page, size, err := util.ParsePaginationParams(r)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_PAGINATION, util.MSG_PAGO_INVALID_PAGE).WithCause(err))
		return
	}

	var idVenta *int64
	if idVentaStr := r.URL.Query().Get("id_venta"); idVentaStr != "" {
		v, parseErr := strconv.ParseInt(idVentaStr, 10, 64)
		if parseErr != nil || v <= 0 {
			pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PAGO_INVALID_ID))
			return
		}
		idVenta = &v
	}

	outs, total, err := h.service.List(page, size, idVenta)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	if outs == nil {
		outs = []input.PagoOutput{}
	}
	_, _, meta := pkg.Paginate(page, size, total)
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_PAGO_LIST, outs, meta)
}
