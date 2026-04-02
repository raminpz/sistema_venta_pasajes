package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sistema_venta_pasajes/internal/programacion/input"
	"sistema_venta_pasajes/internal/programacion/service"
	"sistema_venta_pasajes/internal/programacion/util"
	"sistema_venta_pasajes/pkg"

	"github.com/gorilla/mux"
)

type ProgramacionHandler struct {
	service service.ProgramacionService
}

func NewProgramacionHandler(s service.ProgramacionService) *ProgramacionHandler {
	return &ProgramacionHandler{service: s}
}

func (h *ProgramacionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateProgramacionInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}

	out, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_PROGRAMACION_CREATED, out, nil)
}

func (h *ProgramacionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PROGRAMACION_INVALID_ID))
		return
	}

	var in input.UpdateProgramacionInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}

	out, err := h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	pkg.WriteSuccess(w, http.StatusOK, util.MSG_PROGRAMACION_UPDATED, out, nil)
}

func (h *ProgramacionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PROGRAMACION_INVALID_ID))
		return
	}

	if err := h.service.Delete(id); err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	pkg.WriteSuccess(w, http.StatusOK, util.MSG_PROGRAMACION_DELETED, nil, nil)
}

func (h *ProgramacionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_ID, util.MSG_PROGRAMACION_INVALID_ID))
		return
	}

	out, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	pkg.WriteSuccess(w, http.StatusOK, util.MSG_PROGRAMACION_GET, out, nil)
}

func (h *ProgramacionHandler) List(w http.ResponseWriter, r *http.Request) {
	page, size, err := util.ParsePaginationParams(r)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALID_PAGINATION, util.MSG_PROGRAMACION_INVALID_PAGE).WithCause(err))
		return
	}

	outs, total, err := h.service.List(page, size)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	if outs == nil {
		outs = []input.ProgramacionOutput{}
	}

	_, _, meta := pkg.Paginate(page, size, total)
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_PROGRAMACION_LIST, outs, meta)
}
