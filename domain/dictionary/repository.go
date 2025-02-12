package dictionary

import (
	"database/sql"
	"embed"
	"log"
)

// Repository contains method to interact with data source
type Repository interface {
	GetAlphabets() []Alphabet
	GetWordsByAlphabet(alphabet string) ([]Word, error)
	GetWord(word string) (Word, error)
}

// NewRepository is a function to instantiate new Repository object
func NewRepository(config any) Repository {
	alphabets := []string{
		"a", "b", "c", "d", "g", "h",
		"i", "j", "k", "l", "m", "n",
		"p", "r", "s", "t", "u", "w",
		"y",
	}
	switch config.(type) {
	case embed.FS:
		log.Println("Using embedded filesystem")
		repo := embedRepository{}
		repo.init(config.(embed.FS))
		return &repo
	case bool:
		log.Println("Using read file")
		repo := fsRepository{
			alphabets:  alphabets,
			index:      make(map[string]Alphabet),
			dictionary: make(map[string][]Word),
		}
		repo.init()
		return &repo
	case sql.DB:
		log.Println("Using database connection")
		return nil
	default:
		log.Println("Using json from GitHUB")
		return jsonRepository{alphabets}
	}
}
