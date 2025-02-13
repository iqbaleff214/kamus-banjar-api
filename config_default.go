//go:build !fs

package main

import (
	"database/sql"
	"embed"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

//go:embed data/*.json
var embedded embed.FS

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

	return embedded
}
