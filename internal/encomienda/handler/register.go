package handler

import (
	"sistema_venta_pasajes/internal/encomienda/repository"
	"sistema_venta_pasajes/internal/encomienda/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewEncomiendaRepository(db)
	svc := service.NewEncomiendaService(repo)
	h := NewEncomiendaHandler(svc)

	r.HandleFunc("/encomienda", h.Create).Methods("POST")
	r.HandleFunc("/encomienda", h.List).Methods("GET")
	r.HandleFunc("/encomienda/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/encomienda/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/encomienda/{id}", h.Delete).Methods("DELETE")
}
