package configs

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sistema_venta_pasajes/configs/http/routes"

	"github.com/gorilla/handlers"
)

func Run() error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	db, err := NewMySQL(cfg)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("no se pudo obtener la conexión SQL nativa desde GORM: %w", err)
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	log.Printf("conexion a MySQL OK: %s:%s/%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	router := routes.NewRouter(db)

	// Middleware CORS para desarrollo local
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppPort),
		Handler:      corsHandler,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	log.Printf("servidor HTTP iniciado en http://localhost:%s (%s)", cfg.AppPort, cfg.AppEnv)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("no se pudo iniciar el servidor HTTP: %w", err)
	}

	return nil
}
