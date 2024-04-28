package dictionary

import "errors"

// Service contains all business logic method
type Service interface {
	GetAlphabets() ([]Alphabet, error)
	GetWord(word string) (Word, error)
	GetWordsByAlphabet(alphabet string) (Alphabet, []string, error)
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
	alphabets := s.repository.GetAlphabets()
	return alphabets, nil
}

func (s service) GetWord(word string) (Word, error) {
	alphabet := string(word[0])
	words, err := s.repository.GetWordsByAlphabet(alphabet)
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
	alphabets := s.repository.GetAlphabets()

	// Matching the requested alphabet with the availables
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
	words := []string{}
	for _, w := range wordsByAlphabet {
		words = append(words, w.Word)
	}

	return alphabetObj, words, nil
}
