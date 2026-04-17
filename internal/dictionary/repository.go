package dictionary

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

// Repository contains methods to interact with the data source.
type Repository interface {
	GetAlphabets() ([]Alphabet, error)
	GetWordsByAlphabet(alphabet string) ([]Word, error)
	GetWord(word string) (Word, error)
	Search(keyword string) (SearchResult, error)
}

// NewRepository returns a MySQL-backed Repository.
func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
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

func (r *MockRepository) Search(keyword string) (SearchResult, error) {
	args := r.Called(keyword)
	return args.Get(0).(SearchResult), args.Error(1)
}
