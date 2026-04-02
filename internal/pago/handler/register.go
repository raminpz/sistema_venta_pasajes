package handler

import (
	"sistema_venta_pasajes/internal/pago/repository"
	"sistema_venta_pasajes/internal/pago/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewPagoRepository(db)
	svc := service.NewPagoService(repo)
	h := NewPagoHandler(svc)
	r.HandleFunc("/pago", h.Create).Methods("POST")
	r.HandleFunc("/pago", h.List).Methods("GET")
	r.HandleFunc("/pago/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/pago/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/pago/{id}", h.Delete).Methods("DELETE")
}
