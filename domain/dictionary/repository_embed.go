package dictionary

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
)

type embedRepository struct {
	alphabets  []Alphabet
	dictionary Dictionary
}

func (r *embedRepository) init(data embed.FS) {
	letters := []string{
		"a", "b", "c", "d", "g", "h",
		"i", "j", "k", "l", "m", "n",
		"p", "r", "s", "t", "u", "w",
		"y",
	}

	r.alphabets = make([]Alphabet, len(letters))
	r.dictionary = NewDictionary()
	for i, letter := range letters {
		r.alphabets[i].Letter = letter
		b, err := fs.ReadFile(data, "data/"+letter+".json")
		if err != nil {
			continue
		}

		var result []Word
		if err = json.Unmarshal(b, &result); err == nil {
			r.dictionary.AddEntries(result)
			r.alphabets[i].Total = len(result)
		}
	}
}

func (r *embedRepository) GetAlphabets() ([]Alphabet, error) {
	return r.alphabets, nil
}

func (r *embedRepository) GetWordsByAlphabet(alphabet string) ([]Word, error) {
	entries := r.dictionary.GetEntries(alphabet)
	if len(entries) > 0 {
		return entries, nil
	}

	return entries, errors.New("the alphabet is not available")
}

func (r *embedRepository) GetWord(word string) (Word, error) {
	if word, ok := r.dictionary.GetWord(word); ok {
		return word, nil
	}

	return Word{}, errors.New("the word is not found")
}
