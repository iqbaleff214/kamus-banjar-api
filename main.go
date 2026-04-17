package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iqbaleff214/kamus-banjar-api/domain/dictionary"
)

func main() {
	db := openDB()
	defer db.Close()

	seed(db)

	app := setup(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8001"
	}

	log.Fatal(app.Listen(port))
}

func setup(db *sql.DB) *fiber.App {
	dictionaryRepository := dictionary.NewRepository(db)
	dictionaryService := dictionary.NewService(dictionaryRepository)
	dictionaryHandler := dictionary.NewHandler(dictionaryService)

	app := fiber.New(config())
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

func config() fiber.Config {
	return fiber.Config{
		AppName:      "Kamus Banjar API",
		ErrorHandler: errorHandler,
	}
}
