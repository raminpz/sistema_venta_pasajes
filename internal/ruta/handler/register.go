package handler

import (
	"sistema_venta_pasajes/internal/ruta/repository"
	"sistema_venta_pasajes/internal/ruta/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRutaRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewRutaRepository(db)
	svc := service.New(repo)
	h := New(svc)
	r.HandleFunc("/ruta", h.List).Methods("GET")
	r.HandleFunc("/ruta/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/ruta", h.Create).Methods("POST")
	r.HandleFunc("/ruta/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/ruta/{id}", h.Delete).Methods("DELETE")
}
