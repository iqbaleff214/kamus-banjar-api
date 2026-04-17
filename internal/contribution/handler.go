package contribution

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handler exposes all HTTP handlers for the contribution domain.
type Handler interface {
	// User contributions
	Submit(c *fiber.Ctx) error
	Mine(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Edit(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error

	// Admin — contribution review
	AdminList(c *fiber.Ctx) error
	Approve(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error

	// Admin — word management
	AdminListWords(c *fiber.Ctx) error
	AdminCreateWord(c *fiber.Ctx) error
	AdminUpdateWord(c *fiber.Ctx) error
	AdminDeleteWord(c *fiber.Ctx) error
}

type handler struct {
	svc Service
}

func NewHandler(svc Service) Handler {
	return &handler{svc: svc}
}

// ─────────────────────────────────────────────────────────────
// User contributions
// ─────────────────────────────────────────────────────────────

// POST /api/v1/contributions
func (h *handler) Submit(c *fiber.Ctx) error {
	var req SubmitRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	contrib, err := h.svc.Submit(callerID(c), req)
	if err != nil {
		return mapError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]any{
		"code":    fiber.StatusCreated,
		"status":  "success",
		"message": "Submission created.",
		"data":    contrib,
	})
}

// GET /api/v1/contributions/mine
func (h *handler) Mine(c *fiber.Ctx) error {
	page, limit := pagination(c)

	list, total, err := h.svc.Mine(callerID(c), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(paginatedResponse(list, page, limit, total))
}

// GET /api/v1/contributions/:id
func (h *handler) GetByID(c *fiber.Ctx) error {
	isAdmin := c.Locals("userRole") == "admin"
	contrib, err := h.svc.GetByID(callerID(c), c.Params("id"), isAdmin)
	if err != nil {
		return mapError(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Contribution retrieved.",
		"data":    contrib,
	})
}

// PUT /api/v1/contributions/:id
func (h *handler) Edit(c *fiber.Ctx) error {
	var req SubmitRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	contrib, err := h.svc.Edit(callerID(c), c.Params("id"), req)
	if err != nil {
		return mapError(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Submission updated.",
		"data":    contrib,
	})
}

// DELETE /api/v1/contributions/:id
func (h *handler) Delete(c *fiber.Ctx) error {
	if err := h.svc.Delete(callerID(c), c.Params("id")); err != nil {
		return mapError(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Submission deleted.",
	})
}

// ─────────────────────────────────────────────────────────────
// Admin — contribution review
// ─────────────────────────────────────────────────────────────

// GET /api/v1/admin/contributions
func (h *handler) AdminList(c *fiber.Ctx) error {
	page, limit := pagination(c)
	status := c.Query("status")

	list, total, err := h.svc.ListAll(status, page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(paginatedResponse(list, page, limit, total))
}

// PATCH /api/v1/admin/contributions/:id/approve
func (h *handler) Approve(c *fiber.Ctx) error {
	if err := h.svc.Approve(callerID(c), c.Params("id")); err != nil {
		return mapError(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Contribution approved.",
	})
}

// PATCH /api/v1/admin/contributions/:id/reject
func (h *handler) Reject(c *fiber.Ctx) error {
	var req RejectRequest
	// body is optional — ignore parse error
	_ = c.BodyParser(&req)

	if err := h.svc.Reject(callerID(c), c.Params("id"), req.Notes); err != nil {
		return mapError(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Contribution rejected.",
	})
}

// ─────────────────────────────────────────────────────────────
// Admin — word management
// ─────────────────────────────────────────────────────────────

// GET /api/v1/admin/words
func (h *handler) AdminListWords(c *fiber.Ctx) error {
	page, limit := pagination(c)
	status := c.Query("status")
	source := c.Query("source")

	words, total, err := h.svc.AdminListWords(status, source, page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(paginatedResponse(words, page, limit, total))
}

// POST /api/v1/admin/words
func (h *handler) AdminCreateWord(c *fiber.Ctx) error {
	var req SubmitRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	w, err := h.svc.AdminCreateWord(req)
	if err != nil {
		return mapError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]any{
		"code":    fiber.StatusCreated,
		"status":  "success",
		"message": "Word created.",
		"data":    w,
	})
}

// PUT /api/v1/admin/words/:id
func (h *handler) AdminUpdateWord(c *fiber.Ctx) error {
	wordID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || wordID < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid word id")
	}

	var req SubmitRequest
	if err = c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	w, err := h.svc.AdminUpdateWord(wordID, req)
	if err != nil {
		return mapError(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Word updated.",
		"data":    w,
	})
}

// DELETE /api/v1/admin/words/:id
func (h *handler) AdminDeleteWord(c *fiber.Ctx) error {
	wordID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || wordID < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid word id")
	}

	if err = h.svc.AdminDeleteWord(wordID); err != nil {
		return mapError(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Word deactivated.",
	})
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────

func callerID(c *fiber.Ctx) string {
	id, _ := c.Locals("userID").(string)
	return id
}

func pagination(c *fiber.Ctx) (page, limit int) {
	page = c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit = c.QueryInt("limit", 20)
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return
}

func paginatedResponse(data any, page, limit, total int) map[string]any {
	totalPages := (total + limit - 1) / limit
	return map[string]any{
		"data": data,
		"meta": map[string]any{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrWordNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, ErrForbidden):
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	case errors.Is(err, ErrWordActive), errors.Is(err, ErrNotPending):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, ErrDuplicateWord):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	case IsValidationErr(err):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}
