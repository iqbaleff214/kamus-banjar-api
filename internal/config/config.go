package config

import (
	"database/sql"
	"log"
	"os"
	"time"

	mysqldrv "github.com/go-sql-driver/mysql"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Port          string
	MySQLDSN      string
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
	AdminEmail    string
	AdminPassword string
}

// Load reads env vars, validates required fields, and returns a Config.
func Load() Config {
	cfg := Config{
		Port:          getEnv("PORT", ":8001"),
		MySQLDSN:      os.Getenv("MYSQL_DSN"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTAccessTTL:  parseDuration(getEnv("JWT_ACCESS_TTL", "15m")),
		JWTRefreshTTL: parseDuration(getEnv("JWT_REFRESH_TTL", "168h")),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}

	if cfg.MySQLDSN != "" {
		if cfg.JWTSecret == "" {
			log.Fatal("JWT_SECRET environment variable is required when MYSQL_DSN is set")
		}
		if len(cfg.JWTSecret) < 32 {
			log.Fatal("JWT_SECRET must be at least 32 characters")
		}
	}

	return cfg
}

// OpenDB opens and validates a MySQL connection, injecting parseTime=true.
func (c Config) OpenDB() *sql.DB {
	dsn := withParseTime(c.MySQLDSN)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	if err = db.Ping(); err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	return db
}

// withParseTime ensures parseTime=true is set in the DSN so that DATETIME
// columns scan directly into time.Time values.
func withParseTime(dsn string) string {
	cfg, err := mysqldrv.ParseDSN(dsn)
	if err != nil {
		return dsn
	}
	cfg.ParseTime = true
	return cfg.FormatDSN()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("invalid duration %q: %v", s, err)
	}
	return d
}
