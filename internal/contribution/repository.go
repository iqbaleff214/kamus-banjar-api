package contribution

import (
	"database/sql"
	"time"

	"github.com/iqbaleff214/kamus-banjar-api/internal/dictionary"
	"github.com/stretchr/testify/mock"
)

// Repository defines all persistence operations for the contribution domain.
type Repository interface {
	// Word writes (used by both contribution and admin flows)
	WordExists(word string) (bool, error)
	CreateWord(w dictionary.Word, contributorID, source, status string) (int64, error)
	UpdateWord(wordID int64, w dictionary.Word) error
	SetWordStatus(wordID int64, status string, reviewerID *string, approvedAt *time.Time) error
	GetWordByID(wordID int64) (dictionary.Word, string, error) // (word, status, error)

	// Contribution log
	LogAction(c Contribution) error
	GetByID(id string) (Contribution, error)
	ListByContributor(contributorID string, page, limit int) ([]Contribution, int, error)
	ListAll(status string, page, limit int) ([]Contribution, int, error)

	// Admin word queries
	ListWords(status, source string, page, limit int) ([]dictionary.Word, int, error)
	DeleteWord(wordID int64) error
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

func (m *MockRepository) WordExists(word string) (bool, error) {
	args := m.Called(word)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CreateWord(w dictionary.Word, contributorID, source, status string) (int64, error) {
	args := m.Called(w, contributorID, source, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) UpdateWord(wordID int64, w dictionary.Word) error {
	return m.Called(wordID, w).Error(0)
}

func (m *MockRepository) SetWordStatus(wordID int64, status string, reviewerID *string, approvedAt *time.Time) error {
	return m.Called(wordID, status, reviewerID, approvedAt).Error(0)
}

func (m *MockRepository) GetWordByID(wordID int64) (dictionary.Word, string, error) {
	args := m.Called(wordID)
	return args.Get(0).(dictionary.Word), args.String(1), args.Error(2)
}

func (m *MockRepository) LogAction(c Contribution) error {
	return m.Called(c).Error(0)
}

func (m *MockRepository) GetByID(id string) (Contribution, error) {
	args := m.Called(id)
	return args.Get(0).(Contribution), args.Error(1)
}

func (m *MockRepository) ListByContributor(contributorID string, page, limit int) ([]Contribution, int, error) {
	args := m.Called(contributorID, page, limit)
	return args.Get(0).([]Contribution), args.Int(1), args.Error(2)
}

func (m *MockRepository) ListAll(status string, page, limit int) ([]Contribution, int, error) {
	args := m.Called(status, page, limit)
	return args.Get(0).([]Contribution), args.Int(1), args.Error(2)
}

func (m *MockRepository) ListWords(status, source string, page, limit int) ([]dictionary.Word, int, error) {
	args := m.Called(status, source, page, limit)
	return args.Get(0).([]dictionary.Word), args.Int(1), args.Error(2)
}

func (m *MockRepository) DeleteWord(wordID int64) error {
	return m.Called(wordID).Error(0)
}
