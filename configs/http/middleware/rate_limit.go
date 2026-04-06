package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"sistema_venta_pasajes/pkg"
)

const (
	cleanupInterval = 5 * time.Minute
	idleTimeout     = 5 * time.Minute
)

type ipEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter implementa un limitador de tasa por IP usando token bucket.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*ipEntry
	rps     rate.Limit // tokens por segundo
	burst   int        // ráfaga máxima
}

// NewRateLimiter crea un RateLimiter con la tasa y ráfaga indicadas.
// rps=1/12, burst=5 → 5 requests en ráfaga, recarga 1 cada 12 s (5/min)
// rps=100/60, burst=20 → ~100 req/min para uso general.
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*ipEntry),
		rps:     rate.Limit(rps),
		burst:   burst,
	}
	go rl.cleanup()
	return rl
}

// Middleware devuelve el http.Handler middleware que aplica el rate limiting por IP.
func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			limiter := rl.getLimiter(ip)

			if !limiter.Allow() {
				w.Header().Set("Retry-After", "60")
				pkg.WriteError(w, r, pkg.NewAppError(
					http.StatusTooManyRequests,
					"rate_limit_exceeded",
					"Demasiadas solicitudes. Por favor espere antes de reintentar.",
				))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if e, ok := rl.entries[ip]; ok {
		e.lastSeen = time.Now()
		return e.limiter
	}

	l := rate.NewLimiter(rl.rps, rl.burst)
	rl.entries[ip] = &ipEntry{limiter: l, lastSeen: time.Now()}
	return l
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(cleanupInterval)
		rl.mu.Lock()
		for ip, e := range rl.entries {
			if time.Since(e.lastSeen) > idleTimeout {
				delete(rl.entries, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// clientIP extrae la IP real del cliente respetando proxies.
func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
