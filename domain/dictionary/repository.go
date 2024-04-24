package dictionary

import (
	"embed"
	"encoding/json"
	"io/fs"
)

// Repository contains method to interact with data source
type Repository interface {
	GetAlphabets() []Alphabet
	GetWordsByAlphabet(alphabet string) []Word
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
	alphabets := make([]Alphabet, 26)
	words := make(AlphabeticWord)

	for i := 0; i < 26; i++ {
		letter := string('a' + rune(i))

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

		allWords := []Word{}

		if err := json.Unmarshal(b, &allWords); err != nil {
			return err
		}

		for _, w := range allWords {
			alphabet := w.Alphabet
			words[alphabet] = append(words[alphabet], w)

			index := int([]rune(alphabet)[0] - 'a')
			alphabets[index].Total++
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
func (r repository) GetWordsByAlphabet(alphabet string) []Word {
	return r.words[alphabet]
}
