package server

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iqbaleff214/kamus-banjar-api/internal/dictionary"
)

// New creates and configures the Fiber application with all routes registered.
func New(db *sql.DB) *fiber.App {
	dictionaryRepository := dictionary.NewRepository(db)
	dictionaryService := dictionary.NewService(dictionaryRepository)
	dictionaryHandler := dictionary.NewHandler(dictionaryService)

	app := fiber.New(fiber.Config{
		AppName:      "Kamus Banjar API",
		ErrorHandler: errorHandler,
	})

	app.Use(compress.New())
	app.Use(cors.New())

	api := app.Group("/api")

	api.Get("/v1", rootV1Handler)
	api.Get("/v1/alphabets", dictionaryHandler.GetAlphabets)
	api.Get("/v1/alphabets/:letter", dictionaryHandler.GetWordsByAlphabet)
	api.Get("/v1/entries", dictionaryHandler.Search)
	api.Get("/v1/entries/:word", dictionaryHandler.GetWord)

	return app
}

func rootV1Handler(c *fiber.Ctx) error {
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
	})
}

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
