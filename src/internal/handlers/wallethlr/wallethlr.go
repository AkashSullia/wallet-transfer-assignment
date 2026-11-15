package wallethlr

import (
	service "wallet-transfer-assignment/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type WalletHandler struct {
	service.Service
	logger *zerolog.Logger
}

func NewWalletHandler(InventorySvc WalletService, lg *zerolog.Logger) *WalletHandler {
	logger := lg.With().Str("logger", "user").Logger()
	return &WalletHandler{InventorySvc, &logger}
}

func (h *WalletHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/balance", h.GetWalletBalance)
}
