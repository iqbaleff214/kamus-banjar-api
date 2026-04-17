package main

import (
	"embed"
	"log"
	"os"

	"github.com/iqbaleff214/kamus-banjar-api/internal/config"
	"github.com/iqbaleff214/kamus-banjar-api/internal/seeder"
	"github.com/iqbaleff214/kamus-banjar-api/internal/server"
)

//go:embed data/*.json
var seedData embed.FS

func main() {
	db := config.OpenDB()
	defer db.Close()

	seeder.Seed(db, seedData)

	app := server.New(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8001"
	}

	log.Fatal(app.Listen(port))
}
