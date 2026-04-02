package handler

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"sistema_venta_pasajes/internal/venta/repository"
	"sistema_venta_pasajes/internal/venta/service"
)

type VentaHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

// RegisterRoutes registra las rutas RESTful de venta en el router principal
func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewVentaRepository(db)
	svc := service.NewVentaService(repo)
	h := NewVentaHandler(svc)
	registerRoutesWithHandler(r, h)
}

// registerRoutesWithHandler permite registrar rutas con un handler ya creado (para tests)
func registerRoutesWithHandler(r *mux.Router, h *VentaHandler) {
	r.HandleFunc("/venta", h.Create).Methods("POST")
	r.HandleFunc("/venta", h.List).Methods("GET")
	r.HandleFunc("/ventas", h.List).Methods("GET")
	r.HandleFunc("/venta/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/venta/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/venta/{id}", h.Delete).Methods("DELETE")
}
