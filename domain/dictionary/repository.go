package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// Repository contains method to interact with data source
type Repository interface {
	GetAlphabets() []Alphabet
	GetWordsByAlphabet(alphabet string) ([]Word, error)
}

// repository as a class
type repository struct {
	alphabets []string
}

// NewRepository is a function to instantiate new repository object
//
//	this function also retrieve all content of json file that containing Banjar word dictinories.
//	retrieved contents stored to repository object properties and used as the data source.
func NewRepository() Repository {
	return repository{
		alphabets: []string{
			"a", "b", "c", "d", "g", "h",
			"i", "j", "k", "l", "m", "n",
			"p", "r", "s", "t", "u", "w",
			"y",
		},
	}
}

// GetAlphabets to retrieve all alphabets
func (r repository) GetAlphabets() []Alphabet {
	var mu sync.Mutex
	var wg sync.WaitGroup

	alphabets := make([]Alphabet, len(r.alphabets))

	for i, letter := range r.alphabets {
		endpoint := r.dataSourceUrl(letter)
		wg.Add(1)
		alphabets[i].Letter = letter
		go func(url, letter string, index int) {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(letter, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(letter, err)
				return
			}

			var result []any
			err = json.Unmarshal(body, &result)
			if err != nil {
				fmt.Println(letter, err)
				return
			}

			mu.Lock()
			alphabets[index].Total = len(result)
			mu.Unlock()
		}(endpoint, letter, i)
	}

	wg.Wait()

	return alphabets
}

// GetWordsByAlphabet to retrieve all words that associate with certain alphabet
func (r repository) GetWordsByAlphabet(alphabet string) ([]Word, error) {
	for _, letter := range r.alphabets {
		if letter == alphabet {
			result := []Word{}

			url := r.dataSourceUrl(letter)
			resp, err := http.Get(url)
			if err != nil {
				return result, errors.New("the alphabet is not available")
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return result, errors.New("the alphabet is not available")
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return result, errors.New("the alphabet is not available")
			}

			if err := json.Unmarshal(body, &result); err != nil {
				return result, errors.New("the alphabet is not available")
			}

			return result, nil
		}
	}

	return []Word{}, errors.New("the alphabet is not available")
}

func (r repository) dataSourceUrl(letter string) string {
	return "https://raw.githubusercontent.com/iqbaleff214/kamus-banjar-api/refs/heads/main/data/" + letter + ".json"
}
