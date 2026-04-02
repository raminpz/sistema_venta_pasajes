package handler

import (
	"sistema_venta_pasajes/internal/usuario/repository"
	"sistema_venta_pasajes/internal/usuario/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterUsuarioHandlers(router *mux.Router, db *gorm.DB) {
	repo := repository.NewUsuarioRepository(db) // Usar repo real, no el mock
	serv := service.NewUsuarioService(repo)
	h := NewUsuarioHandler(serv)
	router.HandleFunc("/usuario", h.CrearUsuario).Methods("POST")
	router.HandleFunc("/usuario", h.ListarUsuarios).Methods("GET")
	router.HandleFunc("/usuario/{id}", h.ObtenerUsuarioPorID).Methods("GET")
	router.HandleFunc("/usuario/{id}", h.ActualizarUsuario).Methods("PUT")
	router.HandleFunc("/usuario/{id}", h.EliminarUsuario).Methods("DELETE")
}
