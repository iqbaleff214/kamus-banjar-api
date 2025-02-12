package dictionary

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
)

type embedRepository struct {
	alphabets []Alphabet
	words     map[string][]Word
}

func (r *embedRepository) init(data embed.FS) {
	letters := []string{
		"a", "b", "c", "d", "g", "h",
		"i", "j", "k", "l", "m", "n",
		"p", "r", "s", "t", "u", "w",
		"y",
	}

	r.alphabets = make([]Alphabet, len(letters))
	r.words = make(map[string][]Word)
	for i, letter := range letters {
		r.alphabets[i].Letter = letter
		if b, err := fs.ReadFile(data, "data/"+letter+".json"); err == nil {
			var result []Word
			if err = json.Unmarshal(b, &result); err == nil {
				r.words[letter] = result
				r.alphabets[i].Total = len(result)
			}
		}
	}
}

func (r *embedRepository) GetAlphabets() []Alphabet {
	return r.alphabets
}

func (r *embedRepository) GetWordsByAlphabet(alphabet string) ([]Word, error) {
	words, ok := r.words[alphabet]
	if ok {
		return words, nil
	}

	return words, errors.New("the alphabet is not available")
}
