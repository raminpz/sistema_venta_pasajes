package middleware

import (
	"net/http"
	"sync"
	"time"

	"gorm.io/gorm"

	"sistema_venta_pasajes/pkg"
)

const (
	diasGracia = 30
	ttlCache   = time.Minute // refresca el estado cada minuto
)

var (
	estadoMu        sync.RWMutex
	estadoCache     string
	expiracionCache time.Time
)

// ControlAccesoSistema aplica restricciones de acceso según el estado operativo del sistema,
// consultando la tabla LICENCIA_SISTEMA con caché de 1 minuto.
//   - OPERATIVO    → la solicitud pasa sin restricción
//   - SOLO_LECTURA → solo se permiten peticiones GET (escrituras bloqueadas)
//   - BLOQUEADO    → todas las peticiones son bloqueadas
//
// Si db es nil (ej. en tests sin BD) la solicitud pasa directamente.
// Las rutas de control_acceso (/api/v1/control_acceso/*) NO deben usar este middleware.
func ControlAccesoSistema(db *gorm.DB) func(http.Handler) http.Handler {
	if db == nil {
		return func(next http.Handler) http.Handler { return next }
	}
	return ControlAccesoConFetcher(func() string {
		return obtenerEstadoCacheado(db)
	})
}

// ControlAccesoConFetcher construye el middleware con un proveedor de estado personalizado.
// Útil para tests unitarios: basta con pasar una función que devuelva el estado deseado.
func ControlAccesoConFetcher(fetcher func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			estado := fetcher()
			switch estado {
			case "BLOQUEADO":
				pkg.WriteError(w, r, pkg.Forbidden(
					"sistema_bloqueado",
					"El sistema está bloqueado. Contacte al proveedor para reactivarlo.",
				))
				return
			case "SOLO_LECTURA":
				if r.Method != http.MethodGet {
					pkg.WriteError(w, r, pkg.Forbidden(
						"sistema_solo_lectura",
						"El sistema está en modo solo lectura. Renueve su suscripción para realizar esta operación.",
					))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// InvalidarCacheEstado fuerza una recarga del estado en la próxima solicitud.
// Útil cuando se requiere propagación inmediata tras un cambio de estado.
func InvalidarCacheEstado() {
	estadoMu.Lock()
	expiracionCache = time.Time{}
	estadoMu.Unlock()
}

func obtenerEstadoCacheado(db *gorm.DB) string {
	estadoMu.RLock()
	if estadoCache != "" && time.Now().Before(expiracionCache) {
		s := estadoCache
		estadoMu.RUnlock()
		return s
	}
	estadoMu.RUnlock()

	estadoMu.Lock()
	defer estadoMu.Unlock()
	// Double-check tras adquirir el lock de escritura
	if estadoCache != "" && time.Now().Before(expiracionCache) {
		return estadoCache
	}

	s := consultarEstadoEnBD(db)
	estadoCache = s
	expiracionCache = time.Now().Add(ttlCache)
	return s
}

func consultarEstadoEnBD(db *gorm.DB) string {
	var row struct {
		Estado          string    `gorm:"column:ESTADO"`
		FechaExpiracion time.Time `gorm:"column:FECHA_EXPIRACION"`
	}
	err := db.Table("CONTROL_ACCESO").
		Select("ESTADO, FECHA_EXPIRACION").
		Order("ID_ACCESO DESC").
		First(&row).Error
	if err != nil {
		// Sin registro en BD o error de conexión → BLOQUEADO (fail-safe)
		return "BLOQUEADO"
	}

	if row.Estado == "BLOQUEADO" {
		return "BLOQUEADO"
	}

	now := time.Now()
	if !now.After(row.FechaExpiracion) {
		return "OPERATIVO"
	}

	finGracia := row.FechaExpiracion.AddDate(0, 0, diasGracia)
	if !now.After(finGracia) {
		return "SOLO_LECTURA"
	}

	return "BLOQUEADO"
}
