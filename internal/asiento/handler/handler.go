package handler

import (
	"encoding/json"
	"net/http"
	"sistema_venta_pasajes/internal/asiento/input"
	"sistema_venta_pasajes/internal/asiento/service"
	"sistema_venta_pasajes/internal/asiento/util"
	"sistema_venta_pasajes/pkg"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service service.AsientoService
}

func New(service service.AsientoService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateAsientoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.MsgErrorDecodificarJSON+err.Error()))
		return
	}
	asiento, err := h.service.Create(in)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MsgAsientoCreado, asiento, nil)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.MsgInvalidID))
		return
	}
	asiento, err := h.service.GetByID(id)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MsgOK, asiento, nil)
}

func (h *Handler) ListByVehiculo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idVehiculo, err := strconv.ParseInt(vars["id_vehiculo"], 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.MsgInvalidIDVehiculo))
		return
	}
	asientos, err := h.service.ListByVehiculo(idVehiculo)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MsgOK, asientos, nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.MsgInvalidID))
		return
	}
	var in input.UpdateAsientoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.MsgErrorDecodificarJSON+err.Error()))
		return
	}
	if err := h.service.Update(id, in); err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MsgAsientoActualizado, nil, nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.MsgInvalidID))
		return
	}
	if err := h.service.Delete(id); err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MsgAsientoEliminado, nil, nil)
}

func (h *Handler) CambiarEstado(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.MsgInvalidID))
		return
	}
	var in input.CambiarEstadoAsientoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.MsgErrorDecodificarJSON+err.Error()))
		return
	}
	if err := in.Validate(); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_estado", err.Error()))
		return
	}
	if err := h.service.CambiarEstado(id, in.Estado); err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MsgEstadoActualizado, nil, nil)
}
