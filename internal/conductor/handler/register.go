package handler

import (
	"sistema_venta_pasajes/internal/conductor/repository"
	"sistema_venta_pasajes/internal/conductor/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewConductorRepository(db)
	svc := service.New(repo)
	h := New(svc)
	r.HandleFunc("/conductor", h.Create).Methods("POST")
	r.HandleFunc("/conductor", h.List).Methods("GET")
	r.HandleFunc("/conductor/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/conductor/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/conductor/{id}", h.Delete).Methods("DELETE")
}
