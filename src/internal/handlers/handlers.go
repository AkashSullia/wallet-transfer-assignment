package handlers

import (
	"github.com/gofiber/fiber/v3"
)

type APIHandler interface {
	RegisterRoutes(fiber.Router)
}

type RouterMap map[string]APIHandler
