package handler

import (
	"sistema_venta_pasajes/internal/parada/repository"
	"sistema_venta_pasajes/internal/parada/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewParadaRepository(db)
	svc := service.NewParadaService(repo)
	h := NewParadaHandler(svc)

	r.HandleFunc("/parada", h.Create).Methods("POST")
	r.HandleFunc("/parada/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/parada/{id}", h.Update).Methods("PATCH", "PUT")
	r.HandleFunc("/parada/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/paradas/ruta/{id_ruta}", h.ListByRuta).Methods("GET")
}

