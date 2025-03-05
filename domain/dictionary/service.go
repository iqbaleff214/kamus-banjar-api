package dictionary

import (
	"errors"
	"strings"

	"github.com/iqbaleff214/kamus-banjar-api/domain/dictionary/utils"
)

// Service contains all business logic method
type Service interface {
	GetAlphabets() ([]Alphabet, error)
	GetWord(word string) (Word, error)
	GetWordsByAlphabet(alphabet string) (Alphabet, []string, error)
	Search(keyword string) (SearchResult, error)
}

// service as a class
type service struct {
	repository Repository
}

// NewService is a function to instantiate new service object
func NewService(repository Repository) Service {
	return service{repository}
}

func (s service) GetAlphabets() ([]Alphabet, error) {
	return s.repository.GetAlphabets()
}

func (s service) GetWord(word string) (Word, error) {
	return s.repository.GetWord(word)
}

func (s service) GetWordsByAlphabet(alphabet string) (Alphabet, []string, error) {
	// Check if the requested alphabet is only one character
	if len(alphabet) != 1 {
		return Alphabet{}, nil, errors.New("alphabet only has one character")
	}

	// Check if the requested alphabet is correct and lower cased
	letter := []rune(alphabet)[0]
	if letter < 'a' || letter > 'z' {
		return Alphabet{}, nil, errors.New("symbols are not allowed")
	}

	// Get all available alphabets
	alphabets, err := s.repository.GetAlphabets()
	if err != nil {
		return Alphabet{}, nil, err
	}

	// Matching the requested alphabet with the available
	var alphabetObj Alphabet
	for _, a := range alphabets {
		if alphabet == a.Letter {
			alphabetObj = a
		}
	}

	// Get all words from the requested alphabet
	wordsByAlphabet, err := s.repository.GetWordsByAlphabet(alphabet)
	if err != nil {
		return Alphabet{}, nil, errors.New("the alphabet is not available")
	}

	// Rearrange words as string
	var words []string
	for _, w := range wordsByAlphabet {
		words = append(words, w.Word)
	}

	return alphabetObj, words, nil
}

func (s service) Search(keyword string) (SearchResult, error) {
	transform := strings.ToLower(strings.TrimSpace(keyword))
	if len(transform) == 0 {
		return SearchResult{}, errors.New("search keyword is required")
	}

	if !utils.IsAlphabet(transform) {
		return SearchResult{}, errors.New("symbols are not allowed")
	}

	return s.repository.Search(transform)
}
