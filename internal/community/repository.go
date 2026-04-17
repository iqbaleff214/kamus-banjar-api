package community

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

// Repository defines all persistence operations for the community domain.
type Repository interface {
	// Shared word existence check (only active words).
	ActiveWordExists(word string) (bool, error)

	// Votes
	GetUserVote(word, userID string) (int, error) // 0 = no vote
	UpsertVote(word, userID string, vote int) error
	DeleteVote(word, userID string) error
	GetVoteSummary(word string) (up, down int, err error)

	// Comments
	CreateComment(c Comment) error
	GetComment(id string) (Comment, error)
	ListComments(word string) ([]Comment, error)
	DeleteComment(id string) error

	// Bookmarks
	AddBookmark(word, userID string) error
	RemoveBookmark(word, userID string) error
	ListBookmarks(userID string, page, limit int) ([]string, int, error)

	// Word of the day
	GetWordOfTheDay(date string) (string, error)
	SetWordOfTheDay(word, setBy, date string) error
	GetRandomActiveWord() (string, error)
}

// NewRepository returns a MySQL-backed Repository.
func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

// ─────────────────────────────────────────────────────────────
// Mock
// ─────────────────────────────────────────────────────────────

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ActiveWordExists(word string) (bool, error) {
	args := m.Called(word)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetUserVote(word, userID string) (int, error) {
	args := m.Called(word, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) UpsertVote(word, userID string, vote int) error {
	return m.Called(word, userID, vote).Error(0)
}

func (m *MockRepository) DeleteVote(word, userID string) error {
	return m.Called(word, userID).Error(0)
}

func (m *MockRepository) GetVoteSummary(word string) (up, down int, err error) {
	args := m.Called(word)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockRepository) CreateComment(c Comment) error {
	return m.Called(c).Error(0)
}

func (m *MockRepository) GetComment(id string) (Comment, error) {
	args := m.Called(id)
	return args.Get(0).(Comment), args.Error(1)
}

func (m *MockRepository) ListComments(word string) ([]Comment, error) {
	args := m.Called(word)
	return args.Get(0).([]Comment), args.Error(1)
}

func (m *MockRepository) DeleteComment(id string) error {
	return m.Called(id).Error(0)
}

func (m *MockRepository) AddBookmark(word, userID string) error {
	return m.Called(word, userID).Error(0)
}

func (m *MockRepository) RemoveBookmark(word, userID string) error {
	return m.Called(word, userID).Error(0)
}

func (m *MockRepository) ListBookmarks(userID string, page, limit int) ([]string, int, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]string), args.Int(1), args.Error(2)
}

func (m *MockRepository) GetWordOfTheDay(date string) (string, error) {
	args := m.Called(date)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) SetWordOfTheDay(word, setBy, date string) error {
	return m.Called(word, setBy, date).Error(0)
}

func (m *MockRepository) GetRandomActiveWord() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
