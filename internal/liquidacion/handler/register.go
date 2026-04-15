package handler
import (
"github.com/gorilla/mux"
"gorm.io/gorm"
"sistema_venta_pasajes/internal/liquidacion/repository"
"sistema_venta_pasajes/internal/liquidacion/service"
)
// RegisterRoutes registra las rutas RESTful de liquidacion en el router principal.
func RegisterRoutes(r *mux.Router, db *gorm.DB) {
repo := repository.NewLiquidacionRepository(db)
svc := service.NewLiquidacionService(repo)
h := NewLiquidacionHandler(svc)
r.HandleFunc("/liquidacion", h.Generar).Methods("POST")
r.HandleFunc("/liquidaciones", h.List).Methods("GET")
r.HandleFunc("/liquidacion/{id}", h.GetByID).Methods("GET")
r.HandleFunc("/liquidacion/{id}", h.ActualizarEstado).Methods("PUT")
r.HandleFunc("/liquidacion/{id}", h.Delete).Methods("DELETE")
r.HandleFunc("/programacion/{id_programacion}/caja", h.ObtenerResumenCaja).Methods("GET")
}
