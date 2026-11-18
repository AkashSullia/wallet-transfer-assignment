package handlers

import (
	"wallet-transfer-assignment/internal/usecase"

	"github.com/gofiber/fiber/v3"
)

type APIHandler interface {
	RegisterRoutes(fiber.Router)
}

type RouterMap map[string]APIHandler

type WalletHandler struct {
	w usecase.WalletService
}

type TransferHandler struct {
	t usecase.TransferService
}

type HealthCheck struct{}

func (h *HealthCheck) RegisterRoutes(router fiber.Router) {
	router.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})
}
