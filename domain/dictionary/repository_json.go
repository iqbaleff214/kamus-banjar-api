package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// jsonRepository as a class
type jsonRepository struct {
	alphabets []string
}

// GetAlphabets to retrieve all alphabets
func (r jsonRepository) GetAlphabets() ([]Alphabet, error) {
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
			defer func() { _ = resp.Body.Close() }()

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

	return alphabets, nil
}

// GetWordsByAlphabet to retrieve all words that associate with certain alphabet
func (r jsonRepository) GetWordsByAlphabet(alphabet string) ([]Word, error) {
	for _, letter := range r.alphabets {
		if letter == alphabet {
			var result []Word

			url := r.dataSourceUrl(letter)
			resp, err := http.Get(url)
			if err != nil {
				return result, errors.New("the alphabet is not available")
			}
			defer func() { _ = resp.Body.Close() }()

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

func (r jsonRepository) GetWord(word string) (Word, error) {
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

func (r jsonRepository) dataSourceUrl(letter string) string {
	return "https://raw.githubusercontent.com/iqbaleff214/kamus-banjar-api/refs/heads/main/data/" + letter + ".json"
}

func (r jsonRepository) Search(keyword string) (SearchResult, error) {
	return SearchResult{}, errors.New("method Search isn't implemented for JSON repository")
}
