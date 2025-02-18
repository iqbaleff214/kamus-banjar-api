package dictionary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/ksckaan1/gokachu"
)

var expireTime time.Duration = 1 * time.Hour

// jsonRepository as a class
type jsonRepository struct {
	alphabets  []string
	wordsCache *gokachu.Gokachu[string, []Word]
	dictCache  *gokachu.Gokachu[string, map[string]Word]
}

func newJsonRepository(alphabets []string) jsonRepository {
	config := gokachu.Config{}
	wordsCache := gokachu.New[string, []Word](config)
	dictCache := gokachu.New[string, map[string]Word](config)

	return jsonRepository{
		alphabets:  alphabets,
		wordsCache: wordsCache,
		dictCache:  dictCache,
	}
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
	result := []Word{}

	if !slices.Contains(r.alphabets, alphabet) {
		return result, errors.New("the alphabet is not available")
	}

	// get from cache
	if words, ok := r.wordsCache.Get(alphabet); ok {
		return words, nil
	}

	url := r.dataSourceUrl(alphabet)
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

	// expire in 1 hour
	r.wordsCache.SetWithTTL(alphabet, result, expireTime)

	return result, nil

}

func (r jsonRepository) GetWord(word string) (Word, error) {
	alphabet := string(word[0])

	// check caches first
	if dict, ok := r.dictCache.Get(alphabet); ok {
		if w, ok := dict[word]; ok {
			return w, nil
		}
	}

	if words, ok := r.wordsCache.Get(alphabet); ok {
		dict := createDict(words)
		r.dictCache.SetWithTTL(alphabet, dict, expireTime)

		if w, ok := dict[word]; ok {
			return w, nil
		}
	}

	// fetch them
	words, err := r.GetWordsByAlphabet(alphabet)
	if err != nil {
		return Word{}, errors.New("the word is not found")
	}

	// set cache
	dict := createDict(words)
	r.wordsCache.SetWithTTL(alphabet, words, expireTime)
	r.dictCache.SetWithTTL(alphabet, dict, expireTime)

	if w, ok := dict[word]; ok {
		return w, nil
	}

	return Word{}, errors.New("the word is not found")
}

func createDict(words []Word) map[string]Word {
	dict := make(map[string]Word, len(words))
	for _, w := range words {
		dict[w.Word] = w
	}

	return dict
}

func (r jsonRepository) dataSourceUrl(letter string) string {
	return "https://raw.githubusercontent.com/iqbaleff214/kamus-banjar-api/refs/heads/main/data/" + letter + ".json"
}
