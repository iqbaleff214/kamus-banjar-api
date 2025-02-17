package dictionary

import (
	"database/sql"
	"embed"
	"log"

	"github.com/stretchr/testify/mock"
)

// Repository contains method to interact with data source
type Repository interface {
	GetAlphabets() ([]Alphabet, error)
	GetWordsByAlphabet(alphabet string) ([]Word, error)
	GetWord(word string) (Word, error)
}

// NewRepository is a function to instantiate new Repository object
func NewRepository(config any) Repository {
	alphabets := []string{
		"a", "b", "c", "d", "g", "h",
		"i", "j", "k", "l", "m", "n",
		"p", "r", "s", "t", "u", "w",
		"y",
	}
	switch config.(type) {
	case embed.FS:
		log.Println("Using embedded filesystem")
		repo := embedRepository{}
		repo.init(config.(embed.FS))
		return &repo
	case string:
		log.Println("Using read file")
		repo := fsRepository{
			alphabets:  alphabets,
			index:      make(map[string]Alphabet),
			dictionary: NewDictionary(),
			path:       config.(string),
		}
		repo.init()
		return &repo
	case *sql.DB:
		log.Println("Using database connection")
		return &mysqlRepository{db: config.(*sql.DB)}
	default:
		log.Println("Using json from GitHUB", config)
		return jsonRepository{alphabets}
	}
}

// --------------------------------| MOCK OBJECT |-------------------------------- //

type MockRepository struct {
	mock.Mock
}

func (r *MockRepository) GetAlphabets() ([]Alphabet, error) {
	args := r.Called()
	return args.Get(0).([]Alphabet), args.Error(1)
}

func (r *MockRepository) GetWordsByAlphabet(alphabet string) ([]Word, error) {
	args := r.Called(alphabet)
	return args.Get(0).([]Word), args.Error(1)
}

func (r *MockRepository) GetWord(word string) (Word, error) {
	args := r.Called(word)
	return args.Get(0).(Word), args.Error(1)
}
