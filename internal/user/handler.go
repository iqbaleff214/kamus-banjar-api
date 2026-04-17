package user

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iqbaleff214/kamus-banjar-api/pkg/pagination"
	"github.com/iqbaleff214/kamus-banjar-api/pkg/response"
)

type Handler interface {
	// auth
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	// admin
	ListUsers(c *fiber.Ctx) error
	Deactivate(c *fiber.Ctx) error
	Activate(c *fiber.Ctx) error
	Promote(c *fiber.Ctx) error
}

type handler struct {
	svc Service
}

func NewHandler(svc Service) Handler {
	return &handler{svc: svc}
}

// ─────────────────────────────────────────────────────────────
// Auth
// ─────────────────────────────────────────────────────────────

// POST /api/v1/auth/register
func (h *handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	u, err := h.svc.Register(req)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return response.Created(c, "Account created successfully.", u)
}

// POST /api/v1/auth/login
func (h *handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	pair, err := h.svc.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCreds):
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		case errors.Is(err, ErrAccountInactive):
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		default:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	return response.OK(c, "Login successful.", pair)
}

// POST /api/v1/auth/refresh
func (h *handler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "refresh_token is required")
	}

	pair, err := h.svc.RefreshTokens(req.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return response.OK(c, "Tokens refreshed.", pair)
}

// POST /api/v1/auth/logout  (requires auth)
func (h *handler) Logout(c *fiber.Ctx) error {
	var req LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.svc.Logout(req.RefreshToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.NoContent(c, "Logged out successfully.")
}

// GET /api/v1/auth/me  (requires auth)
func (h *handler) Me(c *fiber.Ctx) error {
	userID := callerID(c)
	u, err := h.svc.Me(userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "Profile retrieved.", u)
}

// PUT /api/v1/auth/me  (requires auth)
func (h *handler) UpdateProfile(c *fiber.Ctx) error {
	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	u, err := h.svc.UpdateProfile(callerID(c), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrWrongPassword):
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		case errors.Is(err, ErrNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	return response.OK(c, "Profile updated.", u)
}

// ─────────────────────────────────────────────────────────────
// Admin — Users
// ─────────────────────────────────────────────────────────────

// GET /api/v1/admin/users
func (h *handler) ListUsers(c *fiber.Ctx) error {
	page, limit := pagination.Parse(c)
	role := c.Query("role")

	var active *bool
	if raw := c.Query("active"); raw != "" {
		v := raw == "true" || raw == "1"
		active = &v
	}

	users, total, err := h.svc.ListUsers(page, limit, role, active)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if users == nil {
		users = []User{}
	}

	return response.Paginated(c, "Users retrieved.", users, pagination.NewMeta(page, limit, total))
}

// PATCH /api/v1/admin/users/:id/deactivate
func (h *handler) Deactivate(c *fiber.Ctx) error {
	return h.setActive(c, false)
}

// PATCH /api/v1/admin/users/:id/activate
func (h *handler) Activate(c *fiber.Ctx) error {
	return h.setActive(c, true)
}

func (h *handler) setActive(c *fiber.Ctx, active bool) error {
	targetID := c.Params("id")
	err := h.svc.SetActive(targetID, callerID(c), active)
	if err != nil {
		switch {
		case errors.Is(err, ErrSelfModify):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case errors.Is(err, ErrNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return response.NoContent(c, "User updated.")
}

// PATCH /api/v1/admin/users/:id/promote
func (h *handler) Promote(c *fiber.Ctx) error {
	if err := h.svc.Promote(c.Params("id")); err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.NoContent(c, "User promoted to admin.")
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────

func callerID(c *fiber.Ctx) string {
	id, _ := c.Locals("userID").(string)
	return id
}

