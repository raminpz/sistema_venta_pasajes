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
	RegisterRoutes(r, h)
}
