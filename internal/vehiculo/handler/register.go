package handler

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"sistema_venta_pasajes/internal/vehiculo/repository"
	"sistema_venta_pasajes/internal/vehiculo/service"
)

type VehiculoHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

// RegisterRoutes registra las rutas RESTful de vehiculo en el router principal
func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewVehiculoRepository(db)
	svc := service.NewVehiculoService(repo)
	h := NewVehiculoHandler(svc)
	registerRoutesWithHandler(r, h)
}

// registerRoutesWithHandler permite registrar rutas con un handler ya creado (para tests)
func registerRoutesWithHandler(r *mux.Router, h *VehiculoHandler) {
	r.HandleFunc("/vehiculo", h.Create).Methods("POST")
	r.HandleFunc("/vehiculos", h.List).Methods("GET")
	r.HandleFunc("/vehiculo/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/vehiculo/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/vehiculo/{id}", h.Delete).Methods("DELETE")
}
