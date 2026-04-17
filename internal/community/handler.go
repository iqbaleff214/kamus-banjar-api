package community

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iqbaleff214/kamus-banjar-api/internal/dictionary"
)

// Handler exposes all HTTP handlers for the community domain.
type Handler interface {
	// Votes
	Vote(c *fiber.Ctx) error
	// Comments
	ListComments(c *fiber.Ctx) error
	PostComment(c *fiber.Ctx) error
	DeleteComment(c *fiber.Ctx) error
	// Bookmarks
	ListBookmarks(c *fiber.Ctx) error
	AddBookmark(c *fiber.Ctx) error
	RemoveBookmark(c *fiber.Ctx) error
	// Word of the day
	GetWordOfTheDay(c *fiber.Ctx) error
	SetWordOfTheDay(c *fiber.Ctx) error
}

type handler struct {
	svc     Service
	dictSvc dictionary.Service
}

// NewHandler creates a community handler. dictSvc is used to fetch full word
// details for the word-of-the-day endpoint.
func NewHandler(svc Service, dictSvc dictionary.Service) Handler {
	return &handler{svc: svc, dictSvc: dictSvc}
}

// ─────────────────────────────────────────────────────────────
// Votes
// ─────────────────────────────────────────────────────────────

// POST /api/v1/entries/:word/votes
func (h *handler) Vote(c *fiber.Ctx) error {
	word := strings.ToLower(c.Params("word"))

	var req VoteRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.svc.CastVote(word, callerID(c), req.Vote); err != nil {
		return mapErr(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Vote recorded.",
	})
}

// ─────────────────────────────────────────────────────────────
// Comments
// ─────────────────────────────────────────────────────────────

// GET /api/v1/entries/:word/comments
func (h *handler) ListComments(c *fiber.Ctx) error {
	word := strings.ToLower(c.Params("word"))

	comments, err := h.svc.ListComments(word)
	if err != nil {
		return mapErr(err)
	}
	if comments == nil {
		comments = []Comment{}
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Comments retrieved.",
		"data":    comments,
	})
}

// POST /api/v1/entries/:word/comments
func (h *handler) PostComment(c *fiber.Ctx) error {
	word := strings.ToLower(c.Params("word"))

	var req CommentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	userName, _ := c.Locals("userName").(string)
	comment, err := h.svc.PostComment(word, callerID(c), userName, req.Body, req.ParentID)
	if err != nil {
		return mapErr(err)
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]any{
		"code":    fiber.StatusCreated,
		"status":  "success",
		"message": "Comment posted.",
		"data":    comment,
	})
}

// DELETE /api/v1/entries/:word/comments/:id  (own comment)
// DELETE /api/v1/admin/entries/:word/comments/:id  (admin)
func (h *handler) DeleteComment(c *fiber.Ctx) error {
	isAdmin := c.Locals("userRole") == "admin"

	if err := h.svc.DeleteComment(c.Params("id"), callerID(c), isAdmin); err != nil {
		return mapErr(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Comment deleted.",
	})
}

// ─────────────────────────────────────────────────────────────
// Bookmarks
// ─────────────────────────────────────────────────────────────

// GET /api/v1/me/bookmarks
func (h *handler) ListBookmarks(c *fiber.Ctx) error {
	page, limit := pagination(c)
	words, total, err := h.svc.ListBookmarks(callerID(c), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	totalPages := (total + limit - 1) / limit
	return c.JSON(map[string]any{
		"data": words,
		"meta": map[string]any{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// POST /api/v1/me/bookmarks/:word
func (h *handler) AddBookmark(c *fiber.Ctx) error {
	word := strings.ToLower(c.Params("word"))

	if err := h.svc.AddBookmark(word, callerID(c)); err != nil {
		return mapErr(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Bookmark added.",
	})
}

// DELETE /api/v1/me/bookmarks/:word
func (h *handler) RemoveBookmark(c *fiber.Ctx) error {
	word := strings.ToLower(c.Params("word"))

	if err := h.svc.RemoveBookmark(word, callerID(c)); err != nil {
		return mapErr(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Bookmark removed.",
	})
}

// ─────────────────────────────────────────────────────────────
// Word of the day
// ─────────────────────────────────────────────────────────────

// GET /api/v1/word-of-the-day
func (h *handler) GetWordOfTheDay(c *fiber.Ctx) error {
	word, err := h.svc.GetWordOfTheDay()
	if err != nil {
		return mapErr(err)
	}

	result, err := h.dictSvc.GetWord(word)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch word detail")
	}

	c.Set("Cache-Control", "public, max-age=3600")
	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Word of the day retrieved.",
		"data":    result,
	})
}

// PUT /api/v1/admin/word-of-the-day
func (h *handler) SetWordOfTheDay(c *fiber.Ctx) error {
	var req WordOfTheDayRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.svc.SetWordOfTheDay(req.Word, callerID(c), req.Date); err != nil {
		return mapErr(err)
	}

	return c.JSON(map[string]any{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Word of the day set.",
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

func mapErr(err error) error {
	switch {
	case errors.Is(err, ErrWordNotFound), errors.Is(err, ErrCommentNotFound), errors.Is(err, ErrNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, ErrForbidden):
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	case errors.Is(err, ErrInvalidVote), errors.Is(err, ErrEmptyBody):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNoWordToday):
		return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}
