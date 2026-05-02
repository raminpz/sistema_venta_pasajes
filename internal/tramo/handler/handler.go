package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"sistema_venta_pasajes/internal/tramo/input"
	"sistema_venta_pasajes/internal/tramo/service"
	"sistema_venta_pasajes/internal/tramo/util"
	"sistema_venta_pasajes/pkg"
)

type TramoHandler struct {
	service service.TramoService
}

func NewTramoHandler(s service.TramoService) *TramoHandler {
	return &TramoHandler{service: s}
}

func (h *TramoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateTramoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.ERR_INVALID_JSON).WithCause(err))
		return
	}
	tramo, err := h.service.Create(in)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_CREATED, tramo, nil)
}

func (h *TramoHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	var in input.UpdateTramoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.ERR_INVALID_JSON).WithCause(err))
		return
	}
	in.IDTramo = id
	tramo, err := h.service.Update(in)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_UPDATED, tramo, nil)
}

func (h *TramoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	if err := h.service.Delete(id); err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_DELETED, nil, nil)
}

func (h *TramoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	tramo, err := h.service.GetByID(id)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_GET, tramo, nil)
}

func (h *TramoHandler) List(w http.ResponseWriter, r *http.Request) {
	page, size, err := pkg.ParsePaginationParams(r, 1, 15)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_pagination", util.ERR_INVALID_PAGE).WithCause(err))
		return
	}
	tramos, total, err := h.service.List(page, size)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	_, _, meta := pkg.Paginate(page, size, total)
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_LIST, tramos, meta)
}

func (h *TramoHandler) ListByRuta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idRuta, err := strconv.ParseInt(vars["id_ruta"], 10, 64)
	if err != nil || idRuta <= 0 {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	tramos, err := h.service.ListByRuta(idRuta)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_LIST, tramos, nil)
}

