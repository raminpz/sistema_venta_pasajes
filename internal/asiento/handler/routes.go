package handler

import (
	"github.com/gorilla/mux"
)

// RegisterRoutes registra las rutas HTTP para el handler de asiento
func RegisterRoutes(r *mux.Router, h *Handler) {
	r.HandleFunc("/asiento", h.Create).Methods("POST")
	r.HandleFunc("/asiento/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/vehiculo/{id_vehiculo}/asientos", h.ListByVehiculo).Methods("GET")
	r.HandleFunc("/asiento/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/asiento/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/asiento/{id}/estado", h.CambiarEstado).Methods("PATCH")
}
