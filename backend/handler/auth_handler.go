package handler

import (
	"fiber3-app/backend/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	sessStore  *session.Store
	repository *repository.AuthRepository
}

func NewAuthHandler(s *session.Store) *AuthHandler {
	return &AuthHandler{sessStore: s, repository: repository.NewAuthRepository()}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	type credentials struct {
		Username string
		Password string
		Role     string
	}
	var req credentials

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString("Invalid bad request")
	}

	acc := h.repository.UserAccount(req.Username)

	if acc.ID == "" {
		return c.Status(404).SendString("No existe ese nombre de usuario")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password)); err != nil {
		return c.Status(400).SendString("Contraseña inválida")
	}

	// Obtener sesión (la cookie se llamará por defecto "session_id")
	sessStore, err := h.sessStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error de sesión")
	}

	// Guardar información del usuario en la sesión
	sessStore.Set("userID", acc.ID)
	sessStore.Set("firstName", acc.Firstname)
	sessStore.Set("lastName", acc.Lastname)
	sessStore.Set("userName", req.Username)
	sessStore.Set("imageSrc", acc.ImageSrc)

	// Guardar la sesión (esto envía la cookie al cliente)
	if err := sessStore.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not save session")
	}

	return c.SendString("success")
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	sess, err := h.sessStore.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	// Destruir sesión
	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not destroy session")
	}

	return c.SendString("success")
}

func (h *AuthHandler) UpdatePasswd(c fiber.Ctx) error {
	var req map[string]string
	c.Bind().Body(&req)

	userName, ok := c.Locals("userName").(string) // get userName from session
	if !ok {
		return c.Status(500).SendString("Datos de sesión no válidos")
	}

	acc := h.repository.UserAccount(userName)
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req["currentPasswd"])); err != nil {
		return c.Status(400).SendString("Contraseña inválida")
	}

	if req["newPasswd"] != req["confirmPasswd"] {
		return c.Status(400).SendString("La nueva y la contraseña de confirmación no coinciden")
	}

	hashPasswd, _ := bcrypt.GenerateFromPassword([]byte(req["newPasswd"]), 14)
	if err := h.repository.UpdatePasswd(userName, string(hashPasswd)); err != nil {
		return c.Status(500).SendString("No se puedo cambiar contraseña")
	}

	return c.SendString("success")
}

func (h *AuthHandler) CurrentUserAcc(c fiber.Ctx) error {
	userName, ok := c.Locals("userName").(string) // get userName from session
	if !ok {
		return c.Status(500).SendString("Datos de sesión no válidos")
	}

	return c.JSON(h.repository.UserAccount(userName))
}
