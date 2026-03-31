package handler

import (
	"sistema_venta_pasajes/internal/usuario/repository"
	"sistema_venta_pasajes/internal/usuario/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UsuarioRegister struct{}

func RegisterUsuarioHandlers(router *mux.Router, db *gorm.DB) {
	repo := repository.NewUsuarioRepository(db) // Usar repo real, no el mock
	serv := service.NewUsuarioService(repo)
	h := NewUsuarioHandler(serv)
	registerRoutesWithHandler(router, h)
}

func registerRoutesWithHandler(r *mux.Router, h *UsuarioHandler) {
	r.HandleFunc("/usuario", h.CrearUsuario).Methods("POST")
	r.HandleFunc("/usuario", h.ListarUsuarios).Methods("GET")
	r.HandleFunc("/usuario/{id}", h.ObtenerUsuarioPorID).Methods("GET")
	r.HandleFunc("/usuario/{id}", h.ActualizarUsuario).Methods("PUT")
	r.HandleFunc("/usuario/{id}", h.EliminarUsuario).Methods("DELETE")
}
