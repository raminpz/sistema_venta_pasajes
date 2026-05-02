package handler

import (
	"sistema_venta_pasajes/internal/tramo/repository"
	"sistema_venta_pasajes/internal/tramo/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewTramoRepository(db)
	svc := service.NewTramoService(repo)
	h := NewTramoHandler(svc)

	r.HandleFunc("/tramo", h.Create).Methods("POST")
	r.HandleFunc("/tramos", h.List).Methods("GET")
	r.HandleFunc("/tramo/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/tramo/{id}", h.Update).Methods("PATCH", "PUT")
	r.HandleFunc("/tramo/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/tramos/ruta/{id_ruta}", h.ListByRuta).Methods("GET")
}

