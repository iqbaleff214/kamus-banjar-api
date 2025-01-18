package main

import (
	"errors"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

// GET /api/v1/
func rootV1Handler(c *fiber.Ctx) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"owner":   "M. Iqbal Effendi <iqbaleff214@gmail.com>",
		"license": "MIT",
		"source":  "https://github.com/iqbaleff214/kamus-banjar-api",
		"version": "1",
		"endpoints": []map[string]string{
			{
				"path":        "/api/v1/alphabets",
				"method":      "GET",
				"description": "Returning a list of alphabet information.",
			},
			{
				"path":        "/api/v1/alphabets/{letter}",
				"method":      "GET",
				"description": "Returning a list of Banjar words according to the given letter.",
			},
			{
				"path":        "/api/v1/entries/{word}",
				"method":      "GET",
				"description": "Returning the definition and meaning according to the given Banjar word.",
			},
		},
		"stats": map[string]float64{
			"Alloc (MB)": float64(memStats.Alloc)/1024/1024,
			"TotalAlloc (MB)": float64(memStats.TotalAlloc)/1024/1024,
			"Sys (MB)": float64(memStats.Sys)/1024/1024,
			"NumGC": float64(memStats.NumGC),
		},
	})
}

// Error handler response
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	err = c.Status(code).JSON(map[string]any{
		"code":    code,
		"message": e.Message,
		"status":  "error",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{
			"code":    fiber.StatusInternalServerError,
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return nil
}
