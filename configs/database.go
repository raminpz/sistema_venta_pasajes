package configs

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(cfg Config) (*gorm.DB, error) {
	dsn := buildDSN(cfg.DB)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir la conexion MySQL con GORM: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la conexion SQL nativa desde GORM: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("no se pudo conectar a MySQL: %w", err)
	}

	return db, nil
}

func buildDSN(cfg DBConfig) string {
	params := strings.TrimSpace(cfg.Params)
	if params == "" {
		params = "parseTime=true&loc=Local&charset=utf8mb4"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		normalizeParams(params),
	)
}

func normalizeParams(params string) string {
	values, err := url.ParseQuery(params)
	if err != nil {
		return params
	}

	if values.Get("parseTime") == "" {
		values.Set("parseTime", "true")
	}
	if values.Get("loc") == "" {
		values.Set("loc", "Local")
	}
	if values.Get("charset") == "" {
		values.Set("charset", "utf8mb4")
	}

	return values.Encode()
}
