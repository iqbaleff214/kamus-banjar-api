package main

import (
	"embed"
	"log"

	"github.com/iqbaleff214/kamus-banjar-api/internal/config"
	"github.com/iqbaleff214/kamus-banjar-api/internal/seeder"
	"github.com/iqbaleff214/kamus-banjar-api/internal/server"
)

//go:embed data/*.json
var seedData embed.FS

func main() {
	cfg := config.Load()

	db := cfg.OpenDB()
	defer db.Close()

	seeder.Seed(db, seedData)
	seeder.SeedAdmin(db, cfg.AdminEmail, cfg.AdminPassword)

	app := server.New(db, cfg)

	log.Fatal(app.Listen(cfg.Port))
}
