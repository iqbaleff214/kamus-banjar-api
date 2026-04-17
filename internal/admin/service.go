package admin

// Service exposes admin dashboard operations.
type Service interface {
	GetStats() (Stats, error)
}

type service struct {
	repo Repository
}

// NewService creates a Service backed by repo.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetStats() (Stats, error) {
	return s.repo.GetStats()
}
