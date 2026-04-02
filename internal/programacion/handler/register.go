package handler

import (
	"sistema_venta_pasajes/internal/programacion/repository"
	"sistema_venta_pasajes/internal/programacion/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewProgramacionRepository(db)
	svc := service.NewProgramacionService(repo)
	h := NewProgramacionHandler(svc)
	r.HandleFunc("/programacion", h.Create).Methods("POST")
	r.HandleFunc("/programacion", h.List).Methods("GET")
	r.HandleFunc("/programacion/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/programacion/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/programacion/{id}", h.Delete).Methods("DELETE")
}
