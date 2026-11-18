package handlers

import (
	"errors"

	"wallet-transfer-assignment/internal/domain"
	"wallet-transfer-assignment/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func NewTransferHandler(t usecase.TransferService) *TransferHandler {
	return &TransferHandler{t}
}

func (h *TransferHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/do", h.Transfer)
}

func (h *TransferHandler) Transfer(c fiber.Ctx) error {
	req := domain.TransferReq{}
	if err := c.Bind().Body(&req); err != nil {
		log.Error().Ctx(c.Context()).Err(err).Msg("Body parsing failed.")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(map[string]string{"msg": "JSON validation error"})
	}

	if len(req.IdempotencyKey) > 100 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(map[string]string{"msg": "IdempotencyKey too large."})
	}

	t, err := h.t.Transfer(c.Context(), req)
	if err != nil {
		log.Error().Ctx(c.Context()).Err(err).Msg("error making transfer.")
		if errors.Is(err, domain.ErrInsufficientFunds) {
			return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"msg": "insufficient fund"})
		}
	}
	return c.Status(fiber.StatusOK).JSON(t)
}
