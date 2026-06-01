package handler

import (
	"sistema_venta_pasajes/internal/tipo_vehiculo/repository"
	"sistema_venta_pasajes/internal/tipo_vehiculo/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewTipoVehiculoRepository(db)
	svc := service.NewTipoVehiculoService(repo)
	h := NewTipoVehiculoHandler(svc)

	r.HandleFunc("/tipo-vehiculo", h.Create).Methods("POST")
	r.HandleFunc("/tipo-vehiculo", h.List).Methods("GET")
	r.HandleFunc("/tipo-vehiculos", h.List).Methods("GET")
	r.HandleFunc("/tipo-vehiculo/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/tipo-vehiculo/{id}", h.Update).Methods("PUT", "PATCH")
	r.HandleFunc("/tipo-vehiculo/{id}", h.Delete).Methods("DELETE")
}

