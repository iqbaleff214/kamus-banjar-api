package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iqbaleff214/kamus-banjar-api/pkg/pagination"
)

// OK sends a 200 success response with data payload.
func OK(c *fiber.Ctx, msg string, data any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": msg,
		"data":    data,
	})
}

// Created sends a 201 success response with data payload.
func Created(c *fiber.Ctx, msg string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"status":  "success",
		"message": msg,
		"data":    data,
	})
}

// NoContent sends a 200 success response with no data field.
func NoContent(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": msg,
	})
}

// Paginated sends a 200 success response with data and pagination meta.
func Paginated(c *fiber.Ctx, msg string, data any, meta pagination.Meta) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": msg,
		"data":    data,
		"meta":    meta,
	})
}
