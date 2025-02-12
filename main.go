package main

import (
	"embed"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iqbaleff214/kamus-banjar-api/domain/dictionary"
)

//go:embed data/*.json
var embedded embed.FS

func main() {
	app := setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8001"
	}

	log.Fatal(app.Listen(port))
}

func setup() *fiber.App {
	repoConfig := make(map[string]any)
	repoConfig["embed"] = embedded
	dictionaryRepository := dictionary.NewRepository(repoConfig)

	dictionaryService := dictionary.NewService(dictionaryRepository)

	dictionaryHandler := dictionary.NewHandler(dictionaryService)

	// app setup
	app := fiber.New(config())
	app.Use(cors.New())

	// api route
	api := app.Group("/api")

	// api v1
	api.Get("/v1", rootV1Handler)

	api.Get("/v1/alphabets", dictionaryHandler.GetAlphabets)
	api.Get("/v1/alphabets/:letter", dictionaryHandler.GetWordsByAlphabet)
	api.Get("/v1/entries/:word", dictionaryHandler.GetWord)

	return app
}

func config() fiber.Config {
	return fiber.Config{
		AppName:      "Kamus Banjar API",
		ErrorHandler: errorHandler,
	}
}
