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
	words := s.repository.GetWordsByAlphabet(alphabet)

	for _, w := range words {
		if word == w.Word {
			return w, nil
		}
	}

	return Word{}, errors.New("the word is not found")
}

func (s service) GetWordsByAlphabet(alphabet string) (Alphabet, []string, error) {
	if len(alphabet) != 1 {
		return Alphabet{}, nil, errors.New("alphabet only has one character")
	}

	alphabets := s.repository.GetAlphabets()
	letter := []rune(alphabet)[0]
	if letter < 'a' || letter > 'z' {
		return Alphabet{}, nil, errors.New("symbols are not allowed")
	}

	alphabetIndex := int(letter - 'a')

	words := []string{}
	for _, w := range s.repository.GetWordsByAlphabet(alphabet) {
		words = append(words, w.Word)
	}

	return alphabets[alphabetIndex], words, nil
}
