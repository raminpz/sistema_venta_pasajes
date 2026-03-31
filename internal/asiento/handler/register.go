package handler

import (
	"sistema_venta_pasajes/internal/asiento/repository"
	"sistema_venta_pasajes/internal/asiento/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterAsientoRoutes crea el handler y registra las rutas RESTful de asiento en el router principal
func RegisterAsientoRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewAsientoRepository(db)
	svc := service.NewAsientoService(repo)
	h := New(svc)
	RegisterRoutes(r, h)
}
