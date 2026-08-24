package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

var store *session.Store

func InitSessionStore(s *session.Store) {
	store = s
}

// Verifica que exista una sesión válida
func Validate(c fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	// Verificar si existe userID en la sesión
	userID := sess.Get("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Pasa datos a los siguientes handlers mediante Locals
	c.Locals("userID", userID)
	c.Locals("firstName", sess.Get("firstName"))
	c.Locals("lastName", sess.Get("lastName"))
	c.Locals("userName", sess.Get("userName"))
	c.Locals("imageSrc", sess.Get("imageSrc"))

	return c.Next()
}
