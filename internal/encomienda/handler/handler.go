package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sistema_venta_pasajes/internal/encomienda/input"
	"sistema_venta_pasajes/internal/encomienda/service"
	"sistema_venta_pasajes/internal/encomienda/util"
	"sistema_venta_pasajes/pkg"

	"github.com/gorilla/mux"
)

type EncomiendaHandler struct {
	service service.EncomiendaService
}

func NewEncomiendaHandler(s service.EncomiendaService) *EncomiendaHandler {
	return &EncomiendaHandler{service: s}
}

func (h *EncomiendaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateEncomiendaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}
	out, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_ENCOMIENDA_CREATED, out, nil)
}

func (h *EncomiendaHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_ENCOMIENDA_INVALID_ID))
		return
	}

	var in input.UpdateEncomiendaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}

	out, err := h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_ENCOMIENDA_UPDATED, out, nil)
}

func (h *EncomiendaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_ENCOMIENDA_INVALID_ID))
		return
	}

	if err := h.service.Delete(id); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_ENCOMIENDA_DELETED, nil, nil)
}

func (h *EncomiendaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_ENCOMIENDA_INVALID_ID))
		return
	}

	out, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_ENCOMIENDA_GET, out, nil)
}

func (h *EncomiendaHandler) List(w http.ResponseWriter, r *http.Request) {
	page, size, err := util.ParsePaginationParams(r)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_PAGINATION, util.MSG_ENCOMIENDA_INVALID_PAGE).WithCause(err))
		return
	}

	outs, total, err := h.service.List(page, size)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	if outs == nil {
		outs = []input.EncomiendaOutput{}
	}

	_, _, meta := pkg.Paginate(page, size, total)
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_ENCOMIENDA_LIST, outs, meta)
}
