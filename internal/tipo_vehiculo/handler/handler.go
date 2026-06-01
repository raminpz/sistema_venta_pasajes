package handler

import (
	"encoding/json"
	"net/http"
	"sistema_venta_pasajes/internal/tipo_vehiculo/input"
	"sistema_venta_pasajes/internal/tipo_vehiculo/service"
	"sistema_venta_pasajes/internal/tipo_vehiculo/util"
	"sistema_venta_pasajes/pkg"
	"strconv"

	"github.com/gorilla/mux"
)

type TipoVehiculoHandler struct {
	service service.TipoVehiculoService
}

func NewTipoVehiculoHandler(s service.TipoVehiculoService) *TipoVehiculoHandler {
	return &TipoVehiculoHandler{service: s}
}

func (h *TipoVehiculoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateTipoVehiculoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}

	out, err := h.service.Create(in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_CREATED, out, nil)
}

func (h *TipoVehiculoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALIDID, util.ERR_INVALID_ID).WithCause(err))
		return
	}

	var in input.UpdateTipoVehiculoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.HandleDecodeError(w, err)
		return
	}

	out, err := h.service.Update(id, in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	pkg.WriteSuccess(w, http.StatusOK, util.MSG_UPDATED, out, nil)
}

func (h *TipoVehiculoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALIDID, util.ERR_INVALID_ID).WithCause(err))
		return
	}

	if err := h.service.Delete(id); err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	pkg.WriteSuccess(w, http.StatusOK, util.MSG_DELETED, nil, nil)
}

func (h *TipoVehiculoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_INVALIDID, util.ERR_INVALID_ID).WithCause(err))
		return
	}

	out, err := h.service.GetByID(id)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	pkg.WriteSuccess(w, http.StatusOK, util.MSG_GET, out, nil)
}

func (h *TipoVehiculoHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.List()
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	if out == nil {
		out = []input.TipoVehiculoOutput{}
	}

	pkg.WriteSuccess(w, http.StatusOK, util.MSG_LIST, out, map[string]any{"count": len(out)})
}

