package admin

import "github.com/gofiber/fiber/v2"

// Handler exposes admin HTTP handlers.
type Handler interface {
	GetStats(c *fiber.Ctx) error
}

type handler struct {
	svc Service
}

// NewHandler creates an admin handler.
func NewHandler(svc Service) Handler {
	return &handler{svc: svc}
}

// GET /api/v1/admin/stats
func (h *handler) GetStats(c *fiber.Ctx) error {
	stats, err := h.svc.GetStats()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"code":    fiber.StatusOK,
		"message": "Stats retrieved successfully.",
		"status":  "success",
		"data":    stats,
	})
}
