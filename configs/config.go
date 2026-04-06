package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppEnv         string
	ProviderAPIKey string // PROVIDER_API_KEY: clave secreta exclusiva del proveedor
	JWTSecret      string // JWT_SECRET: clave para firmar tokens JWT (HS256)
	DB             DBConfig
	HTTP           HTTPConfig
}

type DBConfig struct {
	Host            string
	Port            string
	Name            string
	User            string
	Pass            string
	Params          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type HTTPConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load("config/.env")
	_ = godotenv.Load(".env")

	cfg := Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		ProviderAPIKey: getEnv("PROVIDER_API_KEY", ""),
		JWTSecret:      getEnv("JWT_SECRET", "cambiar_en_produccion_secreto_jwt_256bits"),
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnv("DB_PORT", "3306"),
			Name:            strings.TrimSpace(os.Getenv("DB_NAME")),
			User:            strings.TrimSpace(os.Getenv("DB_USER")),
			Pass:            os.Getenv("DB_PASS"),
			Params:          getEnv("DB_PARAMS", "parseTime=true&loc=Local&charset=utf8mb4"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 30)) * time.Minute,
		},
		HTTP: HTTPConfig{
			ReadTimeout:  time.Duration(getEnvAsInt("HTTP_READ_TIMEOUT", 10)) * time.Second,
			WriteTimeout: time.Duration(getEnvAsInt("HTTP_WRITE_TIMEOUT", 10)) * time.Second,
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var missing []string

	if strings.TrimSpace(c.DB.Name) == "" {
		missing = append(missing, "DB_NAME")
	}
	if strings.TrimSpace(c.DB.User) == "" {
		missing = append(missing, "DB_USER")
	}
	if strings.TrimSpace(c.DB.Host) == "" {
		missing = append(missing, "DB_HOST")
	}
	if strings.TrimSpace(c.DB.Port) == "" {
		missing = append(missing, "DB_PORT")
	}

	if len(missing) > 0 {
		return fmt.Errorf("faltan variables de entorno requeridas: %s", strings.Join(missing, ", "))
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
