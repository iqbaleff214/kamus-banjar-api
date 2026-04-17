package admin

import "github.com/stretchr/testify/mock"

// Repository fetches aggregated platform statistics.
type Repository interface {
	GetStats() (Stats, error)
}

// MockRepository is a testify mock for Repository.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetStats() (Stats, error) {
	args := m.Called()
	return args.Get(0).(Stats), args.Error(1)
}
