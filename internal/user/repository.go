package user

import (
	"database/sql"
	"time"

	"github.com/stretchr/testify/mock"
)

// Repository defines all persistence operations for the user domain.
type Repository interface {
	Create(u User, passwordHash string) error
	FindByEmail(email string) (User, string, error) // (user, passwordHash, error)
	FindByID(id string) (User, error)
	Update(id, name string) error
	UpdatePassword(id, hash string) error
	SaveRefreshToken(userID, tokenHash string, expiresAt time.Time) error
	RevokeRefreshToken(tokenHash string) error
	FindRefreshToken(tokenHash string) (userID string, err error)
	// admin
	ListUsers(page, limit int, role string, active *bool) ([]User, int, error)
	SetActive(id string, active bool) error
	SetRole(id, role string) error
}

// NewRepository returns a MySQL-backed Repository.
func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

// --------------------------------| MOCK |-------------------------------- //

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(u User, passwordHash string) error {
	return m.Called(u, passwordHash).Error(0)
}

func (m *MockRepository) FindByEmail(email string) (User, string, error) {
	args := m.Called(email)
	return args.Get(0).(User), args.Get(1).(string), args.Error(2)
}

func (m *MockRepository) FindByID(id string) (User, error) {
	args := m.Called(id)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockRepository) Update(id, name string) error {
	return m.Called(id, name).Error(0)
}

func (m *MockRepository) UpdatePassword(id, hash string) error {
	return m.Called(id, hash).Error(0)
}

func (m *MockRepository) SaveRefreshToken(userID, tokenHash string, expiresAt time.Time) error {
	return m.Called(userID, tokenHash, expiresAt).Error(0)
}

func (m *MockRepository) RevokeRefreshToken(tokenHash string) error {
	return m.Called(tokenHash).Error(0)
}

func (m *MockRepository) FindRefreshToken(tokenHash string) (string, error) {
	args := m.Called(tokenHash)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockRepository) ListUsers(page, limit int, role string, active *bool) ([]User, int, error) {
	args := m.Called(page, limit, role, active)
	return args.Get(0).([]User), args.Get(1).(int), args.Error(2)
}

func (m *MockRepository) SetActive(id string, active bool) error {
	return m.Called(id, active).Error(0)
}

func (m *MockRepository) SetRole(id, role string) error {
	return m.Called(id, role).Error(0)
}
