package handler

import (
	//"fiber3-app/server/models"
	"fiber3-app/backend/domain"
	"fiber3-app/backend/repository"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type profileHandler struct {
	repository *repository.ProfileRepository
}

func NewProfileHandler() *profileHandler {
	return &profileHandler{repository: repository.NewProfileRepository()}
}

func (h *profileHandler) Create(c fiber.Ctx) error {
	var profile domain.Profile

	if err := c.Bind().Body(&profile); err != nil {
		return err
	}
	profile.ID = uuid.NewString()

	if err := h.repository.Create(&profile); err != nil {
		return c.Status(400).SendString(strings.Split(err.Error(), "\n")[2])
	}

	return c.SendString(profile.ID)
}

func (h *profileHandler) GetAll(c fiber.Ctx) error {
	data := h.repository.GetAll()

	return c.JSON(data)
}

func (h *profileHandler) GetByID(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	return ctx.JSON(h.repository.GetByID(id))
}

func (h *profileHandler) Update(c fiber.Ctx) error {
	var profile domain.Profile

	if err := c.Bind().Body(&profile); err != nil {
		return err
	}

	err := h.repository.Update(&profile)

	if err != nil {
		return c.Status(400).SendString(strings.Split(err.Error(), "\n")[2])
	}

	return c.SendString("success")
}
