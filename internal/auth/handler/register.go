package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	mw "sistema_venta_pasajes/configs/http/middleware"
	"sistema_venta_pasajes/internal/auth/repository"
	"sistema_venta_pasajes/internal/auth/service"
)

// RegisterRoutes registra las rutas públicas de autenticación.
//
//	POST /api/v1/auth/login   → sin auth (con rate limit externo)
//	POST /api/v1/auth/refresh → sin auth
//	POST /api/v1/auth/logout  → requiere JWT (aplicado externamente en el router)
func RegisterRoutes(router *mux.Router, db *gorm.DB, jwtSecret string, loginLimiter *mw.RateLimiter) {
	repo := repository.NewAuthRepository(db)
	svc := service.NewAuthService(repo, jwtSecret)
	h := NewAuthHandler(svc)

	// Login con rate limit estricto (5 intentos/minuto por IP)
	router.Handle("/api/v1/auth/login",
		loginLimiter.Middleware()(http.HandlerFunc(h.Login)),
	).Methods(http.MethodPost)

	// Refresh sin rate limit agresivo
	router.HandleFunc("/api/v1/auth/refresh", h.Refresh).Methods(http.MethodPost)

	// Logout: requiere JWT (el subrouter protegido lo aplica)
	router.HandleFunc("/api/v1/auth/logout", h.Logout).Methods(http.MethodPost)
}
