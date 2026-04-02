package handler

import (
	"sistema_venta_pasajes/internal/vehiculo/repository"
	"sistema_venta_pasajes/internal/vehiculo/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterRoutes registra las rutas RESTful de vehiculo en el router principal
func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewVehiculoRepository(db)
	svc := service.NewVehiculoService(repo)
	h := NewVehiculoHandler(svc)
	r.HandleFunc("/vehiculo", h.Create).Methods("POST")
	r.HandleFunc("/vehiculos", h.List).Methods("GET")
	r.HandleFunc("/vehiculo/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/vehiculo/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/vehiculo/{id}", h.Delete).Methods("DELETE")
}
