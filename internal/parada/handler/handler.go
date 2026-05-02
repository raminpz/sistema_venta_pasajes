package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"sistema_venta_pasajes/internal/parada/input"
	"sistema_venta_pasajes/internal/parada/service"
	"sistema_venta_pasajes/internal/parada/util"
	"sistema_venta_pasajes/pkg"
)

type ParadaHandler struct {
	service service.ParadaService
}

func NewParadaHandler(s service.ParadaService) *ParadaHandler {
	return &ParadaHandler{service: s}
}

func (h *ParadaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateParadaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.ERR_INVALID_JSON).WithCause(err))
		return
	}
	parada, err := h.service.Create(in)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_CREATED, parada, nil)
}

func (h *ParadaHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	var in input.UpdateParadaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.ERR_INVALID_JSON).WithCause(err))
		return
	}
	in.IDParada = id
	parada, err := h.service.Update(in)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_UPDATED, parada, nil)
}

func (h *ParadaHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *ParadaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	parada, err := h.service.GetByID(id)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_GET, parada, nil)
}

func (h *ParadaHandler) ListByRuta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idRuta, err := strconv.ParseInt(vars["id_ruta"], 10, 64)
	if err != nil || idRuta <= 0 {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	paradas, err := h.service.ListByRuta(idRuta)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_LIST, paradas, nil)
}

