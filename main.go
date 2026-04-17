package main

import (
	"database/sql"
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

	var db *sql.DB
	if cfg.MySQLDSN != "" {
		db = cfg.OpenDB()
		defer db.Close()

		seeder.Seed(db, seedData)
		seeder.SeedAdmin(db, cfg.AdminEmail, cfg.AdminPassword)
	} else {
		log.Println("MYSQL_DSN not set — starting in read-only dictionary mode (auth, community, and admin routes disabled)")
	}

	app := server.New(db, cfg)
	log.Fatal(app.Listen(cfg.Port))
}
