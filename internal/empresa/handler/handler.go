package handler

import (
	"encoding/json"
	"net/http"
	"sistema_venta_pasajes/internal/empresa/input"
	"sistema_venta_pasajes/internal/empresa/service"
	"sistema_venta_pasajes/pkg"
	"strconv"

	"github.com/gorilla/mux"
)

type EmpresaHandler struct {
	service service.EmpresaService
}

func NewEmpresaHandler(s service.EmpresaService) *EmpresaHandler {
	return &EmpresaHandler{service: s}
}

func (h *EmpresaHandler) Create(w http.ResponseWriter, r *http.Request) {
	   var in input.CreateEmpresaInput
	   if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			   pkg.HandleDecodeError(w, err)
			   return
	   }
	   out, err := h.service.Create(in)
	   if err != nil {
			   pkg.WriteError(w, r, err)
			   return
	   }
	   pkg.WriteSuccess(w, http.StatusCreated, "Empresa creada correctamente", out, nil)
}

func (h *EmpresaHandler) Update(w http.ResponseWriter, r *http.Request) {
	   vars := mux.Vars(r)
	   idStr := vars["id"]
	   id, err := strconv.ParseInt(idStr, 10, 64)
	   if err != nil {
			   pkg.WriteError(w, r, pkg.BadRequest("invalid_id", "ID inválido"))
			   return
	   }
	   var in input.UpdateEmpresaInput
	   if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			   pkg.HandleDecodeError(w, err)
			   return
	   }
	   out, err := h.service.Update(id, in)
	   if err != nil {
			   pkg.WriteError(w, r, err)
			   return
	   }
	   pkg.WriteSuccess(w, http.StatusOK, "Empresa actualizada correctamente", out, nil)
}

func (h *EmpresaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	   vars := mux.Vars(r)
	   idStr := vars["id"]
	   id, err := strconv.ParseInt(idStr, 10, 64)
	   if err != nil {
			   pkg.WriteError(w, r, pkg.BadRequest("invalid_id", "ID inválido"))
			   return
	   }
	   if err := h.service.Delete(id); err != nil {
			   pkg.WriteError(w, r, err)
			   return
	   }
	   pkg.WriteSuccess(w, http.StatusOK, "Empresa eliminada correctamente", nil, nil)
}

func (h *EmpresaHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	   vars := mux.Vars(r)
	   idStr := vars["id"]
	   id, err := strconv.ParseInt(idStr, 10, 64)
	   if err != nil {
			   pkg.WriteError(w, r, pkg.BadRequest("invalid_id", "ID inválido"))
			   return
	   }
	   out, err := h.service.GetByID(id)
	   if err != nil {
			   pkg.WriteError(w, r, err)
			   return
	   }
	   pkg.WriteSuccess(w, http.StatusOK, "Empresa obtenida correctamente", out, nil)
}

func (h *EmpresaHandler) List(w http.ResponseWriter, r *http.Request) {
	   out, err := h.service.List()
	   if err != nil {
			   pkg.WriteError(w, r, err)
			   return
	   }
	   // Si no hay resultados, devolver array vacío
	   if out == nil {
			   out = []input.EmpresaOutput{}
	   }
	   pkg.WriteSuccess(w, http.StatusOK, "Empresas obtenidas correctamente", out, map[string]any{"count": len(out)})
}
