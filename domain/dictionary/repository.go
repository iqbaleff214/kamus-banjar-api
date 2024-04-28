package dictionary

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"path"
)

// Repository contains method to interact with data source
type Repository interface {
	GetAlphabets() []Alphabet
	GetWordsByAlphabet(alphabet string) ([]Word, error)
}

// repository as a class
type repository struct {
	alphabets []Alphabet
	words     AlphabeticWord
}

// NewRepository is a function to instantiate new repository object
//
//	this function also retrieve all content of json file that containing Banjar word dictinories.
//	retrieved contents stored to repository object properties and used as the data source.
func NewRepository(data embed.FS) Repository {
	availables := []string{
		"a", "b", "c", "d", "g", "h",
		"i", "j", "k", "l", "m", "n",
		"p", "r", "s", "t", "u", "w",
		"y",
	}

	alphabets := make([]Alphabet, len(availables))
	words := make(AlphabeticWord)

	for i, letter := range availables {
		alphabets[i].Letter = letter
		words[letter] = []Word{}
	}

	if err := fs.WalkDir(data, "data", func(filename string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		b, err := data.ReadFile(filename)
		if err != nil {
			return err
		}

		filename = path.Base(filename)
		alphabet := filename[:len(filename)-5]

		allWords := []Word{}

		if err := json.Unmarshal(b, &allWords); err != nil {
			return err
		}

		words[alphabet] = allWords

		for i, letter := range availables {
			if letter == alphabet {
				alphabets[i].Total = len(allWords)
				return nil
			}
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return repository{
		alphabets: alphabets,
		words:     words,
	}
}

// GetAlphabets to retrieve all alphabets
func (r repository) GetAlphabets() []Alphabet {
	return r.alphabets
}

// GetWordsByAlphabet to retrieve all words that associate with certain alphabet
func (r repository) GetWordsByAlphabet(alphabet string) ([]Word, error) {
	words, ok := r.words[alphabet]
	if ok {
		return words, nil
	}

	return words, errors.New("the alphabet is not available")
}
