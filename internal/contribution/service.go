package contribution

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iqbaleff214/kamus-banjar-api/internal/dictionary"
)

// Service defines business logic for the contribution domain.
type Service interface {
	// User-facing
	Submit(contributorID string, req SubmitRequest) (Contribution, error)
	Edit(contributorID, contributionID string, req SubmitRequest) (Contribution, error)
	Delete(contributorID, contributionID string) error
	Mine(contributorID string, page, limit int) ([]Contribution, int, error)
	GetByID(callerID, contributionID string, isAdmin bool) (Contribution, error)

	// Admin — contribution review
	ListAll(status string, page, limit int) ([]Contribution, int, error)
	Approve(reviewerID, contributionID string) error
	Reject(reviewerID, contributionID string, notes string) error

	// Admin — word management
	AdminListWords(status, source string, page, limit int) ([]dictionary.Word, int, error)
	AdminCreateWord(req SubmitRequest) (dictionary.Word, error)
	AdminUpdateWord(wordID int64, req SubmitRequest) (dictionary.Word, error)
	AdminDeleteWord(wordID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ─────────────────────────────────────────────────────────────
// User-facing
// ─────────────────────────────────────────────────────────────

func (s *service) Submit(contributorID string, req SubmitRequest) (Contribution, error) {
	if err := validateSubmit(req); err != nil {
		return Contribution{}, err
	}

	exists, err := s.repo.WordExists(req.Word)
	if err != nil {
		return Contribution{}, err
	}
	if exists {
		return Contribution{}, ErrDuplicateWord
	}

	w := reqToWord(req)
	wordID, err := s.repo.CreateWord(w, contributorID, "community", "pending")
	if err != nil {
		return Contribution{}, err
	}

	c := Contribution{
		ID:            uuid.New().String(),
		WordID:        wordID,
		ContributorID: contributorID,
		Action:        "submitted",
		CreatedAt:     time.Now(),
	}
	if err = s.repo.LogAction(c); err != nil {
		return Contribution{}, err
	}
	c.Word = &w
	return c, nil
}

func (s *service) Edit(contributorID, contributionID string, req SubmitRequest) (Contribution, error) {
	if err := validateSubmit(req); err != nil {
		return Contribution{}, err
	}

	existing, err := s.repo.GetByID(contributionID)
	if err != nil {
		return Contribution{}, err
	}
	if existing.ContributorID != contributorID {
		return Contribution{}, ErrForbidden
	}

	_, status, err := s.repo.GetWordByID(existing.WordID)
	if err != nil {
		return Contribution{}, err
	}
	if status == "active" {
		return Contribution{}, ErrWordActive
	}

	newWord := reqToWord(req)
	if err = s.repo.UpdateWord(existing.WordID, newWord); err != nil {
		return Contribution{}, err
	}
	if err = s.repo.SetWordStatus(existing.WordID, "pending", nil, nil); err != nil {
		return Contribution{}, err
	}

	c := Contribution{
		ID:            uuid.New().String(),
		WordID:        existing.WordID,
		ContributorID: contributorID,
		Action:        "revised",
		CreatedAt:     time.Now(),
	}
	if err = s.repo.LogAction(c); err != nil {
		return Contribution{}, err
	}
	c.Word = &newWord
	return c, nil
}

func (s *service) Delete(contributorID, contributionID string) error {
	existing, err := s.repo.GetByID(contributionID)
	if err != nil {
		return err
	}
	if existing.ContributorID != contributorID {
		return ErrForbidden
	}

	_, status, err := s.repo.GetWordByID(existing.WordID)
	if err != nil {
		return err
	}
	if status != "pending" {
		return ErrNotPending
	}

	// Hard-delete the word; FK CASCADE removes contribution records automatically.
	return s.repo.DeleteWord(existing.WordID)
}

func (s *service) Mine(contributorID string, page, limit int) ([]Contribution, int, error) {
	return s.repo.ListByContributor(contributorID, page, limit)
}

func (s *service) GetByID(callerID, contributionID string, isAdmin bool) (Contribution, error) {
	c, err := s.repo.GetByID(contributionID)
	if err != nil {
		return Contribution{}, err
	}
	if !isAdmin && c.ContributorID != callerID {
		return Contribution{}, ErrForbidden
	}
	return c, nil
}

// ─────────────────────────────────────────────────────────────
// Admin — contribution review
// ─────────────────────────────────────────────────────────────

func (s *service) ListAll(status string, page, limit int) ([]Contribution, int, error) {
	return s.repo.ListAll(status, page, limit)
}

func (s *service) Approve(reviewerID, contributionID string) error {
	c, err := s.repo.GetByID(contributionID)
	if err != nil {
		return err
	}

	now := time.Now()
	if err = s.repo.SetWordStatus(c.WordID, "active", &reviewerID, &now); err != nil {
		return err
	}

	return s.repo.LogAction(Contribution{
		ID:            uuid.New().String(),
		WordID:        c.WordID,
		ContributorID: c.ContributorID,
		ReviewerID:    &reviewerID,
		Action:        "approved",
		CreatedAt:     time.Now(),
	})
}

func (s *service) Reject(reviewerID, contributionID string, notes string) error {
	c, err := s.repo.GetByID(contributionID)
	if err != nil {
		return err
	}

	if err = s.repo.SetWordStatus(c.WordID, "rejected", &reviewerID, nil); err != nil {
		return err
	}

	logEntry := Contribution{
		ID:            uuid.New().String(),
		WordID:        c.WordID,
		ContributorID: c.ContributorID,
		ReviewerID:    &reviewerID,
		Action:        "rejected",
		CreatedAt:     time.Now(),
	}
	if n := strings.TrimSpace(notes); n != "" {
		logEntry.Notes = &n
	}
	return s.repo.LogAction(logEntry)
}

// ─────────────────────────────────────────────────────────────
// Admin — word management
// ─────────────────────────────────────────────────────────────

func (s *service) AdminListWords(status, source string, page, limit int) ([]dictionary.Word, int, error) {
	return s.repo.ListWords(status, source, page, limit)
}

func (s *service) AdminCreateWord(req SubmitRequest) (dictionary.Word, error) {
	if err := validateSubmit(req); err != nil {
		return dictionary.Word{}, err
	}

	exists, err := s.repo.WordExists(req.Word)
	if err != nil {
		return dictionary.Word{}, err
	}
	if exists {
		return dictionary.Word{}, ErrDuplicateWord
	}

	w := reqToWord(req)
	wordID, err := s.repo.CreateWord(w, "", "official", "active")
	if err != nil {
		return dictionary.Word{}, err
	}
	w.ID = wordID
	w.Source = "official"
	w.Status = "active"
	return w, nil
}

func (s *service) AdminUpdateWord(wordID int64, req SubmitRequest) (dictionary.Word, error) {
	if err := validateSubmit(req); err != nil {
		return dictionary.Word{}, err
	}

	existing, _, err := s.repo.GetWordByID(wordID)
	if err != nil {
		return dictionary.Word{}, err
	}

	// If word text changed, check it doesn't collide with another entry.
	if existing.Word != req.Word {
		dup, err := s.repo.WordExists(req.Word)
		if err != nil {
			return dictionary.Word{}, err
		}
		if dup {
			return dictionary.Word{}, ErrDuplicateWord
		}
	}

	newWord := reqToWord(req)
	if err = s.repo.UpdateWord(wordID, newWord); err != nil {
		return dictionary.Word{}, err
	}
	newWord.ID = wordID
	return newWord, nil
}

// AdminDeleteWord soft-deletes a word by setting its status to 'rejected'.
// The word is hidden from public endpoints but remains in the DB for audit purposes.
func (s *service) AdminDeleteWord(wordID int64) error {
	_, _, err := s.repo.GetWordByID(wordID)
	if err != nil {
		return err
	}
	return s.repo.SetWordStatus(wordID, "rejected", nil, nil)
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────

func validateSubmit(req SubmitRequest) error {
	if strings.TrimSpace(req.Word) == "" {
		return ErrInvalidReq("word is required")
	}
	if len(req.Alphabet) != 1 {
		return ErrInvalidReq("alphabet must be a single character")
	}
	if len(req.Meanings) == 0 {
		return ErrInvalidReq("at least one meaning is required")
	}
	return nil
}

func reqToWord(req SubmitRequest) dictionary.Word {
	return dictionary.Word{
		Word:        strings.TrimSpace(req.Word),
		Syllable:    req.Syllable,
		Alphabet:    strings.ToLower(strings.TrimSpace(req.Alphabet)),
		Meanings:    req.Meanings,
		Derivatives: req.Derivatives,
	}
}
