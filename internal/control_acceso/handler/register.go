package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"sistema_venta_pasajes/internal/control_acceso/repository"
	"sistema_venta_pasajes/internal/control_acceso/service"
)

// ServeStatus es una función auxiliar para exponer el endpoint de estado en una ruta pública.
func ServeStatus(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	repo := repository.NewRepository(db)
	svc := service.New(repo)
	h := New(svc)
	h.GetStatus(w, r)
}

// RegisterProtectedRoutes registra las rutas administrativas de control de acceso.
// Solo accesibles con rol PROVEEDOR (el middleware se aplica en el router principal).
//
//	GET  /api/v1/control-acceso         → detalle completo
//	POST /api/v1/control-acceso         → crear nuevo control
//	PUT  /api/v1/control-acceso/{id}/activar
//	PUT  /api/v1/control-acceso/{id}/bloquear
//	PUT  /api/v1/control-acceso/{id}/renovar
func RegisterProtectedRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewRepository(db)
	svc := service.New(repo)
	h := New(svc)

	r.HandleFunc("", h.GetLatest).Methods(http.MethodGet)
	r.HandleFunc("", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/{id:[0-9]+}/activar", h.Activar).Methods(http.MethodPut)
	r.HandleFunc("/{id:[0-9]+}/bloquear", h.Bloquear).Methods(http.MethodPut)
	r.HandleFunc("/{id:[0-9]+}/renovar", h.Renovar).Methods(http.MethodPut)
}
