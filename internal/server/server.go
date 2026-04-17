package server

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iqbaleff214/kamus-banjar-api/internal/config"
	"github.com/iqbaleff214/kamus-banjar-api/internal/dictionary"
	"github.com/iqbaleff214/kamus-banjar-api/internal/middleware"
	"github.com/iqbaleff214/kamus-banjar-api/internal/user"
)

// New creates and configures the Fiber application with all routes registered.
func New(db *sql.DB, cfg config.Config) *fiber.App {
	dictionaryRepository := dictionary.NewRepository(db)
	dictionaryService := dictionary.NewService(dictionaryRepository)
	dictionaryHandler := dictionary.NewHandler(dictionaryService)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	userHandler := user.NewHandler(userService)

	app := fiber.New(fiber.Config{
		AppName:      "Kamus Banjar API",
		ErrorHandler: errorHandler,
	})

	app.Use(compress.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")

	// Public dictionary endpoints
	api.Get("/", rootV1Handler)
	api.Get("/alphabets", dictionaryHandler.GetAlphabets)
	api.Get("/alphabets/:letter", dictionaryHandler.GetWordsByAlphabet)
	api.Get("/entries", dictionaryHandler.Search)
	api.Get("/entries/:word", dictionaryHandler.GetWord)

	// Auth endpoints
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.Refresh)
	auth.Post("/logout", middleware.Auth(cfg.JWTSecret), userHandler.Logout)
	auth.Get("/me", middleware.Auth(cfg.JWTSecret), userHandler.Me)
	auth.Put("/me", middleware.Auth(cfg.JWTSecret), userHandler.UpdateProfile)

	// Admin user management endpoints
	admin := api.Group("/admin", middleware.Auth(cfg.JWTSecret), middleware.Role("admin"))
	admin.Get("/users", userHandler.ListUsers)
	admin.Patch("/users/:id/deactivate", userHandler.Deactivate)
	admin.Patch("/users/:id/activate", userHandler.Activate)
	admin.Patch("/users/:id/promote", userHandler.Promote)

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
	message := "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	if jsonErr := c.Status(code).JSON(map[string]any{
		"code":    code,
		"message": message,
		"status":  "error",
	}); jsonErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{
			"code":    fiber.StatusInternalServerError,
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return nil
}
