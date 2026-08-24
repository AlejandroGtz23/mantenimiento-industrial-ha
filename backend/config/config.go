package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	Port                   string
	DNS                    string
	JWTSecret              string
	QRSecret               string
	EvidenceDirectory      string
	RedisHost              string
	RedisPassword          string
	BootstrapAdminUser     string
	BootstrapAdminPassword string
}

var Env *env

func LoadConfig() {
	// Admite ejecutar desde la raíz (go run ./backend/cmd) o desde backend/cmd.
	if err := godotenv.Load(".env", "../../.env"); err != nil {
		// Las variables de entorno pueden ser inyectadas por producción sin archivo .env.
		if os.Getenv("APP_PORT") == "" || os.Getenv("DB_DSN") == "" {
			log.Fatal(err.Error())
		}
	}

	Env = &env{
		Port:              os.Getenv("APP_PORT"),
		DNS:               os.Getenv("DB_DSN"),
		JWTSecret:         envOrDefault("JWT_SECRET", "cambia-este-secreto-jwt-en-produccion"),
		QRSecret:          envOrDefault("QR_SECRET", "cambia-este-secreto-qr-en-produccion"),
		EvidenceDirectory: envOrDefault("EVIDENCE_DIRECTORY", "./uploads/mantenimiento"),
		// En Compose basta REDIS_HOST=redis; localmente es opcional para no romper
		// el flujo de desarrollo existente si Redis no está iniciado.
		RedisHost:              os.Getenv("REDIS_HOST"),
		RedisPassword:          os.Getenv("REDIS_PASSWORD"),
		BootstrapAdminUser:     os.Getenv("BOOTSTRAP_ADMIN_USER"),
		BootstrapAdminPassword: os.Getenv("BOOTSTRAP_ADMIN_PASSWORD"),
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
