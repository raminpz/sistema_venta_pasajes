package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"sistema_venta_pasajes/configs/http/middleware"
)

// ─── Rate Limit ───────────────────────────────────────────────────────────────

func TestRateLimit_PermiteRequestsDentroDelLimite(t *testing.T) {
	// 5 req de ráfaga → las 5 primeras deben pasar
	rl := middleware.NewRateLimiter(5.0/60.0, 5)
	mw := rl.Middleware()

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		w := httptest.NewRecorder()
		mw(okHandlerJWT).ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "request %d debería pasar", i+1)
	}
}

func TestRateLimit_BloqueaCuandoExcedeLimite(t *testing.T) {
	// burst=3 → la 4ta solicitud debe ser rechazada
	rl := middleware.NewRateLimiter(0.001, 3) // recarga muy lenta
	mw := rl.Middleware()

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "10.0.0.2:5678"
		w := httptest.NewRecorder()
		mw(okHandlerJWT).ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// La 4ta debe ser bloqueada
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "10.0.0.2:5678"
	w := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "rate_limit_exceeded")
}

func TestRateLimit_IPsDiferentesNoComparten(t *testing.T) {
	// Cada IP tiene su propio limiter
	rl := middleware.NewRateLimiter(0.001, 2) // burst=2 por IP
	mw := rl.Middleware()

	// IP A: consume 2 (se agota)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "192.168.1.1:1111"
		w := httptest.NewRecorder()
		mw(okHandlerJWT).ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// IP A: bloqueada
	reqA := httptest.NewRequest(http.MethodPost, "/login", nil)
	reqA.RemoteAddr = "192.168.1.1:1111"
	wA := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(wA, reqA)
	assert.Equal(t, http.StatusTooManyRequests, wA.Code)

	// IP B: debe pasar (no afectada por IP A)
	reqB := httptest.NewRequest(http.MethodPost, "/login", nil)
	reqB.RemoteAddr = "192.168.1.2:2222"
	wB := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(wB, reqB)
	assert.Equal(t, http.StatusOK, wB.Code, "IP diferente no debe ser afectada")
}

func TestRateLimit_HeaderRetryAfter(t *testing.T) {
	// Cuando se bloquea, debe incluir el header Retry-After
	rl := middleware.NewRateLimiter(0.001, 1)
	mw := rl.Middleware()

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "10.0.0.3:9999"
		w := httptest.NewRecorder()
		mw(okHandlerJWT).ServeHTTP(w, req)
		if w.Code == http.StatusTooManyRequests {
			assert.NotEmpty(t, w.Header().Get("Retry-After"))
			return
		}
	}
}

func TestRateLimit_XForwardedFor(t *testing.T) {
	// Detecta IP real detrás de proxy
	rl := middleware.NewRateLimiter(0.001, 1)
	mw := rl.Middleware()

	req1 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req1.Header.Set("X-Forwarded-For", "203.0.113.5, 10.0.0.1")
	req1.RemoteAddr = "10.0.0.1:1234"
	w1 := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Segunda petición desde misma IP real → bloqueada
	req2 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req2.Header.Set("X-Forwarded-For", "203.0.113.5, 10.0.0.1")
	req2.RemoteAddr = "10.0.0.1:5678" // distinto puerto, misma IP real
	w2 := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}

// ─── Simulación de ataque brute-force al login ────────────────────────────────

func TestAtaque_BruteForce_Login(t *testing.T) {
	// Simula 20 intentos de login desde la misma IP → debe bloquear después del burst
	rl := middleware.NewRateLimiter(5.0/60.0, 5) // 5 req/min, burst=5
	mw := rl.Middleware()

	bloqueados := 0
	for i := 0; i < 20; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		req.RemoteAddr = "172.16.0.50:4444"
		w := httptest.NewRecorder()
		mw(okHandlerJWT).ServeHTTP(w, req)
		if w.Code == http.StatusTooManyRequests {
			bloqueados++
		}
	}
	assert.Greater(t, bloqueados, 0, "Al menos una solicitud debe ser bloqueada en un ataque de fuerza bruta")
}

func TestAtaque_BruteForce_Concurrente(t *testing.T) {
	// Múltiples goroutines intentando login simultáneamente desde la misma IP
	rl := middleware.NewRateLimiter(5.0/60.0, 5)
	mw := rl.Middleware()

	var wg sync.WaitGroup
	bloqueados := 0
	var mu sync.Mutex

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
			req.RemoteAddr = "10.20.30.40:8888"
			w := httptest.NewRecorder()
			mw(okHandlerJWT).ServeHTTP(w, req)
			if w.Code == http.StatusTooManyRequests {
				mu.Lock()
				bloqueados++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	assert.Greater(t, bloqueados, 0, "El ataque concurrente debe ser limitado")
}

func TestRateLimit_SeRecuperaConElTiempo(t *testing.T) {
	// rps=10 (muy rápido), burst=2 → se agota rápido pero recarga en ms
	rl := middleware.NewRateLimiter(10.0, 2)
	mw := rl.Middleware()

	// Agotar el burst
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "1.2.3.4:5678"
		w := httptest.NewRecorder()
		mw(okHandlerJWT).ServeHTTP(w, req)
	}

	// Esperar que el limiter se recargue (rps=10 → 1 token cada 100ms)
	time.Sleep(250 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "1.2.3.4:5678"
	w := httptest.NewRecorder()
	mw(okHandlerJWT).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Debe pasar después de que el limiter se recargue")
}
