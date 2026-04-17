package admin_test

import (
	"errors"
	"testing"

	"github.com/iqbaleff214/kamus-banjar-api/internal/admin"
	"github.com/stretchr/testify/assert"
)

func newService(repo admin.Repository) admin.Service {
	return admin.NewService(repo)
}

func TestGetStats_Success(t *testing.T) {
	expected := admin.Stats{
		Words: admin.WordStats{
			Official:  1200,
			Community: 45,
			Pending:   12,
			Rejected:  8,
		},
		Users: admin.UserStats{
			Total:    320,
			Active:   310,
			Inactive: 10,
		},
		Contributions: admin.ContributionStats{
			Pending:   12,
			ThisWeek:  5,
			ThisMonth: 18,
		},
		TopContributors: []admin.TopContributor{
			{UserID: "u1", Name: "Alice", ApprovedCount: 10},
			{UserID: "u2", Name: "Bob", ApprovedCount: 7},
		},
	}

	repo := new(admin.MockRepository)
	repo.On("GetStats").Return(expected, nil)

	stats, err := newService(repo).GetStats()

	assert.NoError(t, err)
	assert.Equal(t, expected, stats)
	repo.AssertExpectations(t)
}

func TestGetStats_RepoError(t *testing.T) {
	repo := new(admin.MockRepository)
	repo.On("GetStats").Return(admin.Stats{}, errors.New("db down"))

	_, err := newService(repo).GetStats()

	assert.EqualError(t, err, "db down")
}

func TestGetStats_EmptyTopContributors(t *testing.T) {
	expected := admin.Stats{
		TopContributors: []admin.TopContributor{},
	}

	repo := new(admin.MockRepository)
	repo.On("GetStats").Return(expected, nil)

	stats, err := newService(repo).GetStats()

	assert.NoError(t, err)
	assert.Empty(t, stats.TopContributors)
}
