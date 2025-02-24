package dictionary

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type mysqlRepository struct {
	db *sql.DB
}

func (r *mysqlRepository) GetAlphabets() ([]Alphabet, error) {
	var result []Alphabet

	query := "SELECT l.letter, COUNT(w.id) as total FROM letters l JOIN words w ON l.letter = w.letter GROUP BY l.letter;"
	rows, err := r.db.Query(query)
	if err != nil {
		return result, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		alphabet := Alphabet{}
		if err = rows.Scan(&alphabet.Letter, &alphabet.Total); err != nil {
			return result, err
		}
		result = append(result, alphabet)
	}

	return result, nil
}

func (r *mysqlRepository) GetWordsByAlphabet(alphabet string) ([]Word, error) {
	var result []Word

	// not a fuzzy search, but eh not bad
	query := "SELECT w.word FROM words w WHERE w.letter = ?;"
	rows, err := r.db.Query(query, alphabet)
	if err != nil {
		return result, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		word := Word{}
		if err = rows.Scan(&word.Word); err != nil {
			return result, err
		}
		result = append(result, word)
	}

	if len(result) == 0 {
		return result, errors.New("the alphabet is not available")
	}

	return result, nil
}

func (r *mysqlRepository) GetWord(word string) (Word, error) {
	var result Word

	query := "SELECT w.data FROM words w WHERE w.word = ?;"

	var data string
	if err := r.db.QueryRow(query, word).Scan(&data); err != nil {
		return result, errors.New("the word is not found")
	}

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return result, err
	}

	return result, nil
}

func (r *mysqlRepository) Search(word string) (SearchResult, error) {
	result := SearchResult{
		Search: word,
		Words:  []string{},
	}
	query := "SELECT w.word FROM words w WHERE w.word LIKE ?;"
	rows, err := r.db.Query(query, "%"+word+"%")
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		word := Word{}
		if err = rows.Scan(&word.Word); err != nil {
			return result, err
		}
		result.Words = append(result.Words, word.Word)
	}

	result.Total = len(result.Words)
	if result.Total == 0 {
		return result, errors.New("no matching word found")
	}

	return result, nil
}
