package handler

import (
	"sistema_venta_pasajes/internal/pasajero/repository"
	"sistema_venta_pasajes/internal/pasajero/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewPasajeroRepository(db)
	svc := service.NewPasajeroService(repo)
	h := NewPasajeroHandler(svc)
	r.HandleFunc("/pasajero", h.Create).Methods("POST")
	r.HandleFunc("/pasajeros", h.List).Methods("GET")
	r.HandleFunc("/pasajero/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/pasajero/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/pasajero/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/pasajeros/search", h.Search).Methods("GET")
}
