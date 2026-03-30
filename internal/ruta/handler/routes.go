package handler

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, h *Handler) {
	r.HandleFunc("/ruta", h.List).Methods("GET")
	r.HandleFunc("/ruta/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/ruta", h.Create).Methods("POST")
	r.HandleFunc("/ruta/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/ruta/{id}", h.Delete).Methods("DELETE")
}
