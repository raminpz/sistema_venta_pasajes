package handler

import (
	"sistema_venta_pasajes/internal/terminal/repository"
	"sistema_venta_pasajes/internal/terminal/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewTerminalRepository(db)
	svc := service.NewTerminalService(repo)
	h := NewTerminalHandler(svc)
	r.HandleFunc("/terminal", h.Create).Methods("POST")
	r.HandleFunc("/terminal", h.List).Methods("GET")
	r.HandleFunc("/terminal/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/terminal/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/terminal/{id}", h.Delete).Methods("DELETE")
}
