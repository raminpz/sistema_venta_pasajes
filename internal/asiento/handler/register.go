package handler

import (
	"sistema_venta_pasajes/internal/asiento/repository"
	"sistema_venta_pasajes/internal/asiento/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterAsientoRoutes crea el handler y registra las rutas RESTful de asiento en el router principal
func RegisterAsientoRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewAsientoRepository(db)
	svc := service.NewAsientoService(repo)
	h := New(svc)
	r.HandleFunc("/asiento", h.Create).Methods("POST")
	r.HandleFunc("/asiento/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/vehiculo/{id_vehiculo}/asientos", h.ListByVehiculo).Methods("GET")
	r.HandleFunc("/asiento/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/asiento/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/asiento/{id}/estado", h.CambiarEstado).Methods("PATCH")
}
