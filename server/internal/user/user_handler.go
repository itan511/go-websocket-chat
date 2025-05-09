package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var u CreateUserReq
	if err := c.BodyParser(&u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.Service.CreateUser(c.Context(), &u)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var user LoginUserReq
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	u, err := h.Service.Login(c.Context(), &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    u.accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		Domain:   "localhost",
		HTTPOnly: true,
		Secure:   false,
	})

	return c.Status(fiber.StatusOK).JSON(u)
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Прошедшая дата для удаления cookie
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logout successful"})
}
