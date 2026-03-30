package handler

import (
	"sistema_venta_pasajes/internal/proveedor/repository"
	"sistema_venta_pasajes/internal/proveedor/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := NewHandler(svc)
	RegisterRoutesWithHandler(r, h)
}

// RegisterRoutesWithHandler permite registrar rutas con un handler ya creado (para tests)
func RegisterRoutesWithHandler(r *mux.Router, h *Handler) {
	r.HandleFunc("/proveedor", h.Create).Methods("POST")
	r.HandleFunc("/proveedor", h.List).Methods("GET")
	r.HandleFunc("/proveedor/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/proveedor/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/proveedor/{id}", h.Delete).Methods("DELETE")
}
