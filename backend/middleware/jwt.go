package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fiber3-app/backend/config"
	"fiber3-app/backend/database"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Claims struct {
	Subject  string `json:"sub"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Expires  int64  `json:"exp"`
}

func CreateJWT(subject, username, role string) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	expiresAt := time.Now().Add(8 * time.Hour)
	claims, err := json.Marshal(Claims{Subject: subject, Username: username, Role: role, Expires: expiresAt.Unix()})
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(claims)
	signingInput := header + "." + payload
	token := signingInput + "." + jwtSign(config.Env.JWTSecret, signingInput)
	if err := database.SaveSession(token, expiresAt); err != nil {
		return "", err
	}
	return token, nil
}

func RequireJWT(c fiber.Ctx) error {
	authorization := c.Get("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "JWT de administrador requerido"})
	}
	token := strings.TrimPrefix(authorization, "Bearer ")
	parts := strings.Split(token, ".")
	if len(parts) != 3 || subtle.ConstantTimeCompare([]byte(jwtSign(config.Env.JWTSecret, parts[0]+"."+parts[1])), []byte(parts[2])) != 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "JWT inválido"})
	}
	bytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "JWT inválido"})
	}
	var claims Claims
	if json.Unmarshal(bytes, &claims) != nil || claims.Expires < time.Now().Unix() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "JWT expirado"})
	}
	active, err := database.SessionIsActive(token)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "No se pudo validar la sesión"})
	}
	if !active {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Sesión no válida"})
	}
	c.Locals("claims", claims)
	return c.Next()
}

func RequireAdminJWT(c fiber.Ctx) error {
	if err := RequireJWT(c); err != nil {
		return err
	}
	if c.Locals("claims").(Claims).Role != "ADMIN" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Acceso administrativo requerido"})
	}
	return nil
}

func jwtSign(secret, value string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
