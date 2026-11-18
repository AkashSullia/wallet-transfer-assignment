package handlers

import (
	"errors"
	"uuid"
	"wallet-transfer-assignment/internal/domain"
	"wallet-transfer-assignment/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func NewWalletHandler(w usecase.WalletService) *WalletHandler {
	return &WalletHandler{w}
}

func (h *WalletHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/create", h.CreateWallet)
	router.Get("/balance", h.GetBalance)
}

func (h *WalletHandler) CreateWallet(c fiber.Ctx) error {
	req := domain.CreateWalletReq{}
	if err := c.Bind().Body(&req); err != nil {
		log.Error().Ctx(c.Context()).Err(err).Msg("Body parsing failed.")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(map[string]string{"msg": "JSON validation error"})
	}

	wallet, err := h.w.CreateWallet(c.Context(), req.OpeningBalance)
	if err != nil {
		if errors.Is(err, domain.ErrNegativeOpeningBalance) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(map[string]string{"msg": "opening balance can't be negative"})
		}
		log.Error().Ctx(c.Context()).Err(err).Msg("Err getting wallet balance")
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"msg": "err getting balance"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "Wallet created succesfully.", "data": wallet})
}

func (h *WalletHandler) GetBalance(c fiber.Ctx) error {
	walletID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Ctx(c.Context()).Err(err).Msg("Err processing request")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(map[string]string{"msg": "uuid is not valid"})
	}

	wallet, err := h.w.GetWalletBalance(c.Context(), walletID)
	if err != nil {
		log.Error().Ctx(c.Context()).Err(err).Msg("Err getting wallet balance")
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"msg": "err getting balance"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"balance": wallet.Balance})
}
