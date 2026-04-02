package handler

import (
	"sistema_venta_pasajes/internal/empresa/repository"
	"sistema_venta_pasajes/internal/empresa/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewEmpresaRepository(db)
	svc := service.NewEmpresaService(repo)
	h := NewEmpresaHandler(svc)
	r.HandleFunc("/empresa", h.Create).Methods("POST")
	r.HandleFunc("/empresa", h.List).Methods("GET")
	r.HandleFunc("/empresa/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/empresa/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/empresa/{id}", h.Delete).Methods("DELETE")
}
