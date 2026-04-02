package handler

import (
	"encoding/json"
	"net/http"
	"sistema_venta_pasajes/internal/vehiculo/input"
	"sistema_venta_pasajes/internal/vehiculo/service"
	"sistema_venta_pasajes/internal/vehiculo/util"
	"sistema_venta_pasajes/pkg"
	"strconv"

	"github.com/gorilla/mux"
)

type VehiculoHandler struct {
	service service.VehiculoService
}

func NewVehiculoHandler(s service.VehiculoService) *VehiculoHandler {
	return &VehiculoHandler{service: s}
}

func (h *VehiculoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in input.CreateVehiculoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.ERR_INVALID_JSON).WithCause(err))
		return
	}
	vehiculo, err := h.service.Create(in)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, util.MSG_CREATED, vehiculo, nil)
}

func (h *VehiculoHandler) Update(w http.ResponseWriter, r *http.Request) {
	var in input.UpdateVehiculoInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_json", util.ERR_INVALID_JSON).WithCause(err))
		return
	}
	vehiculo, err := h.service.Update(in)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_UPDATED, vehiculo, nil)
}

func (h *VehiculoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
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

func (h *VehiculoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_id", util.ERR_INVALID_ID))
		return
	}
	vehiculo, err := h.service.GetByID(id)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_GET, vehiculo, nil)
}

func (h *VehiculoHandler) List(w http.ResponseWriter, r *http.Request) {
	page, size, err := pkg.ParsePaginationParams(r, 1, 15)
	if err != nil {
		pkg.Error(w, pkg.BadRequest("invalid_pagination", util.ERR_INVALID_PAGE).WithCause(err))
		return
	}
	vehiculos, total, err := h.service.List(page, size)
	if err != nil {
		pkg.Error(w, err)
		return
	}
	_, _, meta := pkg.Paginate(page, size, total)
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_LIST, vehiculos, meta)
}
