package routes

import (
	"context"
	"io"
	"log"
	"net/http"
	middleware2 "sistema_venta_pasajes/configs/http/middleware"
	"sistema_venta_pasajes/pkg"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	asientohandler "sistema_venta_pasajes/internal/asiento/handler"
	authhandler "sistema_venta_pasajes/internal/auth/handler"
	conductorhandler "sistema_venta_pasajes/internal/conductor/handler"
	licenciahandler "sistema_venta_pasajes/internal/control_acceso/handler"
	empresahandler "sistema_venta_pasajes/internal/empresa/handler"
	encomiendahandler "sistema_venta_pasajes/internal/encomienda/handler"
	liquidacionhandler "sistema_venta_pasajes/internal/liquidacion/handler"
	pagohandler "sistema_venta_pasajes/internal/pago/handler"
	pasajerohandler "sistema_venta_pasajes/internal/pasajero/handler"
	programacionhandler "sistema_venta_pasajes/internal/programacion/handler"
	proveedorsistemahandler "sistema_venta_pasajes/internal/proveedor/handler"
	rutahandler "sistema_venta_pasajes/internal/ruta/handler"
	terminalhandler "sistema_venta_pasajes/internal/terminal/handler"
	usuariohandler "sistema_venta_pasajes/internal/usuario/handler"
	vehiculohandler "sistema_venta_pasajes/internal/vehiculo/handler"
	ventahandler "sistema_venta_pasajes/internal/venta/handler"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func NewRouter(db *gorm.DB, providerAPIKey string, jwtSecret string) *mux.Router {
	router := mux.NewRouter()

	// ── Middlewares globales ──────────────────────────────────────────────────
	router.Use(loggingMiddleware)
	router.Use(middleware2.RequestID)
	router.Use(middleware2.Recoverer)

	router.NotFoundHandler = adapt(func(w http.ResponseWriter, r *http.Request) error {
		return pkg.NotFound("route_not_found", "ruta no encontrada")
	})
	router.MethodNotAllowedHandler = adapt(func(w http.ResponseWriter, r *http.Request) error {
		return pkg.MethodNotAllowed("método HTTP no permitido para esta ruta")
	})

	// ── Health checks (públicos) ──────────────────────────────────────────────
	router.Handle("/health", adapt(healthHandler)).Methods("GET")
	router.Handle("/ready", adapt(readyHandler(db))).Methods("GET")

	// ── Rate limiters ─────────────────────────────────────────────────────────
	// Login: 5 intentos por minuto por IP (anti brute-force)
	loginLimiter := middleware2.NewRateLimiter(5.0/60.0, 5)
	// API general: 120 req/minuto por IP
	apiLimiter := middleware2.NewRateLimiter(120.0/60.0, 20)

	// ── Estado del sistema (siempre público) ──────────────────────────────────
	router.HandleFunc("/api/v1/control-acceso/status",
		func(w http.ResponseWriter, r *http.Request) {
			licenciahandler.ServeStatus(w, r, db)
		}).Methods(http.MethodGet)

	// ── Rutas públicas de autenticación (con rate limit) ─────────────────────
	authhandler.RegisterRoutes(router, db, jwtSecret, loginLimiter)

	// ── Subrouter protegido: JWT + control de acceso por estado ───────────────
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(apiLimiter.Middleware())
	api.Use(middleware2.JWTAuth(jwtSecret))
	api.Use(middleware2.ControlAccesoSistema(db))

	api.Handle("/health", adapt(healthHandler)).Methods("GET")
	api.Handle("/ready", adapt(readyHandler(db))).Methods("GET")

	// ── Control de acceso del sistema: solo PROVEEDOR ─────────────────────────
	proveedorRouter := api.PathPrefix("/control-acceso").Subrouter()
	proveedorRouter.Use(middleware2.RequireRole("PROVEEDOR"))
	licenciahandler.RegisterProtectedRoutes(proveedorRouter, db, providerAPIKey)

	// ── Datos maestros: ADMIN + PROVEEDOR ─────────────────────────────────────
	adminRouter := api.NewRoute().Subrouter()
	adminRouter.Use(middleware2.RequireRole("ADMIN", "PROVEEDOR"))

	terminalhandler.RegisterRoutes(adminRouter, db)
	empresahandler.RegisterRoutes(adminRouter, db)
	conductorhandler.RegisterRoutes(adminRouter, db)
	rutahandler.RegisterRutaRoutes(adminRouter, db)
	vehiculohandler.RegisterRoutes(adminRouter, db)
	asientohandler.RegisterAsientoRoutes(adminRouter, db)
	proveedorsistemahandler.RegisterRoutes(adminRouter, db)
	usuariohandler.RegisterUsuarioHandlers(adminRouter, db)

	// ── Operaciones: ADMIN + VENDEDOR + PROVEEDOR ─────────────────────────────
	operRouter := api.NewRoute().Subrouter()
	operRouter.Use(middleware2.RequireRole("ADMIN", "VENDEDOR", "PROVEEDOR"))

	pasajerohandler.RegisterRoutes(operRouter, db)
	programacionhandler.RegisterRoutes(operRouter, db)
	pagohandler.RegisterRoutes(operRouter, db)
	encomiendahandler.RegisterRoutes(operRouter, db)
	ventahandler.RegisterRoutes(operRouter, db)
	liquidacionhandler.RegisterRoutes(operRouter, db)

	// ── OPTIONS preflight CORS ────────────────────────────────────────────────
	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return router
}

// Middleware de logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Context().Value("request_id")
		if requestID == nil {
			requestID = "-"
		}
		logPrefix := func() string {
			return "[REQUEST_ID=" + requestID.(string) + "]"
		}

		log.Printf("%s %s %s %s", logPrefix(), r.Method, r.RequestURI, r.RemoteAddr)

		// Loguear el body solo para POST, PUT, PATCH
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil && len(bodyBytes) > 0 {
				log.Printf("%s Body: %s", logPrefix(), strings.ReplaceAll(string(bodyBytes), "\n", " "))
				// Restaurar el body para el siguiente handler
				r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
			}
		}

		next.ServeHTTP(w, r)
		log.Printf("%s %s %s completed in %v", logPrefix(), r.Method, r.RequestURI, time.Since(start))
	})
}

func adapt(handler AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value("request_id")
		if requestID == nil {
			requestID = "-"
		}
		if err := handler(w, r); err != nil {
			appErr := pkg.AsAppError(err)
			log.Printf("[REQUEST_ID=%v] [ERROR] %s %s: %v | details: %+v", requestID, r.Method, r.RequestURI, err, appErr.Details)
			pkg.WriteError(w, r, err)
		}
	})
}

func healthHandler(w http.ResponseWriter, _ *http.Request) error {
	pkg.WriteSuccess(w, http.StatusOK, "servicio disponible", map[string]any{
		"status":  "ok",
		"service": "sistema_venta_pasajes",
		"time":    time.Now().Format(time.RFC3339),
	}, nil)
	return nil
}

func readyHandler(db *gorm.DB) AppHandler {
	return func(w http.ResponseWriter, _ *http.Request) error {
		if db == nil {
			return pkg.ServiceUnavailable("database_unavailable", "sin conexión a base de datos")
		}

		sqlDB, err := db.DB()
		if err != nil {
			return pkg.ServiceUnavailable("database_unavailable", "sin conexión a base de datos").WithCause(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			return pkg.ServiceUnavailable("database_unavailable", "sin conexión a base de datos").WithCause(err)
		}

		pkg.WriteSuccess(w, http.StatusOK, "conexión a base de datos disponible", map[string]any{
			"status":  "ok",
			"service": "sistema_venta_pasajes",
			"time":    time.Now().Format(time.RFC3339),
		}, nil)

		return nil
	}
}
