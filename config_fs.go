//go:build fs

package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
		if os.Getenv("SOURCE_PATH") != "" {
			return os.Getenv("SOURCE_PATH")
		}

		curDir, _ := os.Getwd()
		return curDir
	}

	return nil
}
