package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"io/fs"
	"log"

	"github.com/iqbaleff214/kamus-banjar-api/domain/dictionary"
)

//go:embed data/*.json
var seedData embed.FS

var seedLetters = []string{
	"a", "b", "c", "d", "g", "h",
	"i", "j", "k", "l", "m", "n",
	"p", "r", "s", "t", "u", "w",
	"y",
}

// seed populates the database from embedded JSON files on first startup.
// It is a no-op if the words table already has rows.
func seed(db *sql.DB) {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&count); err != nil || count > 0 {
		return
	}

	log.Println("seeding database from JSON data files...")

	for _, letter := range seedLetters {
		b, err := fs.ReadFile(seedData, "data/"+letter+".json")
		if err != nil {
			log.Printf("seed: skip letter '%s': %v", letter, err)
			continue
		}

		var words []dictionary.Word
		if err = json.Unmarshal(b, &words); err != nil {
			log.Printf("seed: skip letter '%s': %v", letter, err)
			continue
		}

		if _, err = db.Exec("INSERT IGNORE INTO letters (letter) VALUES (?)", letter); err != nil {
			log.Printf("seed: insert letter '%s': %v", letter, err)
			continue
		}

		for _, w := range words {
			data, _ := json.Marshal(w)
			if _, err = db.Exec(
				"INSERT IGNORE INTO words (word, letter, source, data) VALUES (?, ?, 'official', ?)",
				w.Word, letter, string(data),
			); err != nil {
				log.Printf("seed: insert word '%s': %v", w.Word, err)
			}
		}

		log.Printf("seed: inserted %d words for letter '%s'", len(words), letter)
	}

	log.Println("seeding complete")
}
