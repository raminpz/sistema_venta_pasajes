package handler

import (
	asientoTramoRepository "sistema_venta_pasajes/internal/asiento_tramo/repository"
	asientoTramoService "sistema_venta_pasajes/internal/asiento_tramo/service"
	"sistema_venta_pasajes/internal/venta/repository"
	"sistema_venta_pasajes/internal/venta/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterRoutes registra las rutas RESTful de venta en el router principal
func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewVentaRepository(db)

	// Inicializar servicio de asiento_tramo
	atRepo := asientoTramoRepository.NewAsientoTramoRepository(db)
	atSvc := asientoTramoService.NewAsientoTramoService(atRepo)

	svc := service.NewVentaService(repo, atSvc)
	h := NewVentaHandler(svc)
	r.HandleFunc("/venta", h.Create).Methods("POST")
	r.HandleFunc("/venta", h.List).Methods("GET")
	r.HandleFunc("/ventas", h.List).Methods("GET")
	r.HandleFunc("/venta/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/venta/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/venta/{id}", h.Delete).Methods("DELETE")
}
