//go:build fs

package main

import (
	"database/sql"
	"log"
	"os"
)

func repoConfig() any {
	if os.Getenv("MYSQL_DSN") != "" {
		db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
		if err != nil {
			log.Fatal(err)
		}
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(0)
		return db
	}

	if _, err := os.ReadDir("data"); err == nil {
		return true
	}

	return nil
}
