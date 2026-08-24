package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fiber3-app/backend/config"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis guarda la lista de tokens de sesión válidos. Así las dos réplicas del
// backend consultan el mismo estado y una sesión no depende del contenedor que
// atendió el inicio de sesión.
var Redis *redis.Client

func ConnectRedis() {
	if config.Env.RedisHost == "" {
		log.Println("Redis no configurado: se usarán JWT sin estado en desarrollo local")
		return
	}
	address := config.Env.RedisHost
	if !strings.Contains(address, ":") {
		address += ":6379"
	}
	Redis = redis.NewClient(&redis.Options{Addr: address, Password: config.Env.RedisPassword})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Redis.Ping(ctx).Err(); err != nil {
		log.Fatal("Error de conexión a Redis: ", err)
	}
	log.Println("¡Conectado a Redis para sesiones compartidas!")
}

func SaveSession(token string, expiresAt time.Time) error {
	if Redis == nil {
		return nil
	}
	return Redis.Set(context.Background(), sessionKey(token), "active", time.Until(expiresAt)).Err()
}

func SessionIsActive(token string) (bool, error) {
	if Redis == nil {
		return true, nil
	}
	count, err := Redis.Exists(context.Background(), sessionKey(token)).Result()
	return count == 1, err
}

func sessionKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return fmt.Sprintf("maintenance:session:%s", hex.EncodeToString(sum[:]))
}
