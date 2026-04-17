package server

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iqbaleff214/kamus-banjar-api/internal/admin"
	"github.com/iqbaleff214/kamus-banjar-api/internal/community"
	"github.com/iqbaleff214/kamus-banjar-api/internal/config"
	"github.com/iqbaleff214/kamus-banjar-api/internal/contribution"
	"github.com/iqbaleff214/kamus-banjar-api/internal/dictionary"
	"github.com/iqbaleff214/kamus-banjar-api/internal/middleware"
	"github.com/iqbaleff214/kamus-banjar-api/internal/user"
)

// New creates and configures the Fiber application with all routes registered.
// When db is nil (no MYSQL_DSN configured), only the public dictionary endpoints
// are registered. Auth, community write, and admin routes require MySQL.
func New(db *sql.DB, cfg config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Kamus Banjar API",
		ErrorHandler: errorHandler,
	})

	app.Use(compress.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")
	api.Get("/", rootV1Handler)

	if db == nil {
		// Dictionary-only mode: no MySQL available.
		dictRepo := dictionary.NewRepository(db)
		dictSvc := dictionary.NewService(dictRepo)
		dictHandler := dictionary.NewHandler(dictSvc, nil)

		api.Get("/alphabets", dictHandler.GetAlphabets)
		api.Get("/alphabets/:letter", dictHandler.GetWordsByAlphabet)
		api.Get("/entries", dictHandler.Search)
		api.Get("/entries/:word", dictHandler.GetWord)

		return app
	}

	// ── MySQL-backed services ────────────────────────────────────

	dictionaryRepository := dictionary.NewRepository(db)
	dictionaryService := dictionary.NewService(dictionaryRepository)

	communityRepository := community.NewRepository(db)
	communityService := community.NewService(communityRepository)
	communityHandler := community.NewHandler(communityService, dictionaryService)

	// Pass communityService as VoteCounter so GetWord includes vote counts.
	dictionaryHandler := dictionary.NewHandler(dictionaryService, communityService)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	userHandler := user.NewHandler(userService)

	contribRepository := contribution.NewRepository(db)
	contribService := contribution.NewService(contribRepository)
	contribHandler := contribution.NewHandler(contribService)

	adminRepository := admin.NewRepository(db)
	adminService := admin.NewService(adminRepository)
	adminHandler := admin.NewHandler(adminService)

	// ── Public dictionary endpoints ──────────────────────────────
	api.Get("/alphabets", dictionaryHandler.GetAlphabets)
	api.Get("/alphabets/:letter", dictionaryHandler.GetWordsByAlphabet)
	api.Get("/entries", dictionaryHandler.Search)
	api.Get("/entries/:word", dictionaryHandler.GetWord)

	// ── Public community endpoints ───────────────────────────────
	api.Get("/word-of-the-day", communityHandler.GetWordOfTheDay)
	api.Get("/entries/:word/comments", communityHandler.ListComments)

	// ── Authenticated community endpoints ────────────────────────
	api.Post("/entries/:word/votes", middleware.Auth(cfg.JWTSecret), communityHandler.Vote)
	api.Post("/entries/:word/comments", middleware.Auth(cfg.JWTSecret), communityHandler.PostComment)
	api.Delete("/entries/:word/comments/:id", middleware.Auth(cfg.JWTSecret), communityHandler.DeleteComment)

	// ── Bookmark endpoints ───────────────────────────────────────
	me := api.Group("/me", middleware.Auth(cfg.JWTSecret))
	me.Get("/bookmarks", communityHandler.ListBookmarks)
	me.Post("/bookmarks/:word", communityHandler.AddBookmark)
	me.Delete("/bookmarks/:word", communityHandler.RemoveBookmark)

	// ── Auth endpoints ───────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.Refresh)
	auth.Post("/logout", middleware.Auth(cfg.JWTSecret), userHandler.Logout)
	auth.Get("/me", middleware.Auth(cfg.JWTSecret), userHandler.Me)
	auth.Put("/me", middleware.Auth(cfg.JWTSecret), userHandler.UpdateProfile)

	// ── User contribution endpoints ──────────────────────────────
	contrib := api.Group("/contributions", middleware.Auth(cfg.JWTSecret))
	contrib.Post("/", contribHandler.Submit)
	contrib.Get("/mine", contribHandler.Mine)
	contrib.Get("/:id", contribHandler.GetByID)
	contrib.Put("/:id", contribHandler.Edit)
	contrib.Delete("/:id", contribHandler.Delete)

	// ── Admin endpoints ──────────────────────────────────────────
	adminGrp := api.Group("/admin", middleware.Auth(cfg.JWTSecret), middleware.Role("admin"))
	// Admin — users
	adminGrp.Get("/users", userHandler.ListUsers)
	adminGrp.Patch("/users/:id/deactivate", userHandler.Deactivate)
	adminGrp.Patch("/users/:id/activate", userHandler.Activate)
	adminGrp.Patch("/users/:id/promote", userHandler.Promote)
	// Admin — contributions
	adminGrp.Get("/contributions", contribHandler.AdminList)
	adminGrp.Patch("/contributions/:id/approve", contribHandler.Approve)
	adminGrp.Patch("/contributions/:id/reject", contribHandler.Reject)
	// Admin — words
	adminGrp.Get("/words", contribHandler.AdminListWords)
	adminGrp.Post("/words", contribHandler.AdminCreateWord)
	adminGrp.Put("/words/:id", contribHandler.AdminUpdateWord)
	adminGrp.Delete("/words/:id", contribHandler.AdminDeleteWord)
	// Admin — community
	adminGrp.Delete("/entries/:word/comments/:id", communityHandler.DeleteComment)
	adminGrp.Put("/word-of-the-day", communityHandler.SetWordOfTheDay)
	// Admin — stats
	adminGrp.Get("/stats", adminHandler.GetStats)

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
