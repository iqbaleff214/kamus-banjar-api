package pagination

import "github.com/gofiber/fiber/v2"

// Meta holds page metadata returned in paginated responses.
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Parse extracts page and limit from query params with safe defaults and bounds.
func Parse(c *fiber.Ctx) (page, limit int) {
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

// NewMeta builds a Meta from page, limit, and total item count.
func NewMeta(page, limit, total int) Meta {
	pages := 0
	if limit > 0 {
		pages = (total + limit - 1) / limit
	}
	return Meta{Page: page, Limit: limit, Total: total, TotalPages: pages}
}
