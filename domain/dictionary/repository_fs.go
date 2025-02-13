package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type fsRepository struct {
	alphabets  []string
	index      map[string]Alphabet
	dictionary map[string][]Word
	path       string
}

func (r *fsRepository) GetAlphabets() ([]Alphabet, error) {
	var alphabets []Alphabet
	for _, letter := range r.alphabets {
		if v, ok := r.index[letter]; ok {
			alphabets = append(alphabets, v)
		}
	}

	return alphabets, nil
}

func (r *fsRepository) GetWordsByAlphabet(alphabet string) ([]Word, error) {
	var result []Word
	if !slices.Contains(r.alphabets, alphabet) {
		return result, errors.New("the alphabet is not available")
	}

	entries, ok := r.dictionary[alphabet]
	if !ok {
		return result, errors.New("the alphabet is not available")
	}

	return entries, nil
}

func (r *fsRepository) GetWord(word string) (Word, error) {
	alphabet := string(word[0])
	words, err := r.GetWordsByAlphabet(alphabet)
	if err != nil {
		return Word{}, errors.New("the word is not found")
	}

	for _, w := range words {
		if word == w.Word {
			return w, nil
		}
	}

	return Word{}, errors.New("the word is not found")
}

func (r *fsRepository) init() {
	if err := r.loadLetters(); err != nil {
		panic(err)
	}

	if err := r.loadWords(); err != nil {
		panic(err)
	}
}

// load letters
func (r *fsRepository) loadLetters() error {
	for _, letter := range r.alphabets {
		letterFile := filepath.Join(r.path, "data", letter+".json")

		content, err := os.ReadFile(letterFile)
		if err != nil {
			fmt.Println(letter, err)
			return err
		}

		var result []any
		err = json.Unmarshal(content, &result)
		if err != nil {
			fmt.Println(letter, err)
			return err
		}

		r.index[letter] = Alphabet{
			Letter: letter,
			Total:  len(result),
		}
	}

	return nil
}

func (r *fsRepository) loadWords() error {
	for _, letter := range r.alphabets {
		letterFile := filepath.Join(r.path, "data", letter+".json")

		content, err := os.ReadFile(letterFile)
		if err != nil {
			return errors.New("the alphabet is not available")
		}

		var result []Word
		if err := json.Unmarshal(content, &result); err != nil {
			return errors.New("the alphabet is not available")
		}

		for _, word := range result {
			r.dictionary[letter] = append(r.dictionary[letter], word)
		}
	}

	return nil
}
