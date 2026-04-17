package community

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service defines business logic for community features.
type Service interface {
	// Votes
	CastVote(word, userID string, vote int) error
	GetVoteSummary(word string) (up, down int, err error)

	// Comments
	PostComment(word, userID, userName, body string, parentID *string) (Comment, error)
	ListComments(word string) ([]Comment, error)
	DeleteComment(commentID, callerID string, isAdmin bool) error

	// Bookmarks
	AddBookmark(word, userID string) error
	RemoveBookmark(word, userID string) error
	ListBookmarks(userID string, page, limit int) ([]string, int, error)

	// Word of the day
	GetWordOfTheDay() (string, error)
	SetWordOfTheDay(word, adminID, date string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ─────────────────────────────────────────────────────────────
// Votes
// ─────────────────────────────────────────────────────────────

func (s *service) CastVote(word, userID string, vote int) error {
	if vote != 1 && vote != -1 {
		return ErrInvalidVote
	}

	exists, err := s.repo.ActiveWordExists(word)
	if err != nil {
		return err
	}
	if !exists {
		return ErrWordNotFound
	}

	current, err := s.repo.GetUserVote(word, userID)
	if err != nil {
		return err
	}

	if current == vote {
		// Same direction → toggle off.
		return s.repo.DeleteVote(word, userID)
	}
	return s.repo.UpsertVote(word, userID, vote)
}

func (s *service) GetVoteSummary(word string) (up, down int, err error) {
	return s.repo.GetVoteSummary(word)
}

// ─────────────────────────────────────────────────────────────
// Comments
// ─────────────────────────────────────────────────────────────

func (s *service) PostComment(word, userID, userName, body string, parentID *string) (Comment, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return Comment{}, ErrEmptyBody
	}

	exists, err := s.repo.ActiveWordExists(word)
	if err != nil {
		return Comment{}, err
	}
	if !exists {
		return Comment{}, ErrWordNotFound
	}

	if parentID != nil && *parentID != "" {
		parent, err := s.repo.GetComment(*parentID)
		if err != nil {
			if errors.Is(err, ErrCommentNotFound) {
				return Comment{}, ErrCommentNotFound
			}
			return Comment{}, err
		}
		// Replies must be on the same word.
		if parent.Word != word {
			return Comment{}, ErrCommentNotFound
		}
	} else {
		parentID = nil
	}

	c := Comment{
		ID:        uuid.New().String(),
		Word:      word,
		UserID:    userID,
		UserName:  userName,
		ParentID:  parentID,
		Body:      body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err = s.repo.CreateComment(c); err != nil {
		return Comment{}, err
	}
	return c, nil
}

func (s *service) ListComments(word string) ([]Comment, error) {
	return s.repo.ListComments(word)
}

func (s *service) DeleteComment(commentID, callerID string, isAdmin bool) error {
	c, err := s.repo.GetComment(commentID)
	if err != nil {
		return err
	}
	if !isAdmin && c.UserID != callerID {
		return ErrForbidden
	}
	return s.repo.DeleteComment(commentID)
}

// ─────────────────────────────────────────────────────────────
// Bookmarks
// ─────────────────────────────────────────────────────────────

func (s *service) AddBookmark(word, userID string) error {
	exists, err := s.repo.ActiveWordExists(word)
	if err != nil {
		return err
	}
	if !exists {
		return ErrWordNotFound
	}
	return s.repo.AddBookmark(word, userID)
}

func (s *service) RemoveBookmark(word, userID string) error {
	return s.repo.RemoveBookmark(word, userID)
}

func (s *service) ListBookmarks(userID string, page, limit int) ([]string, int, error) {
	return s.repo.ListBookmarks(userID, page, limit)
}

// ─────────────────────────────────────────────────────────────
// Word of the day
// ─────────────────────────────────────────────────────────────

func (s *service) GetWordOfTheDay() (string, error) {
	today := time.Now().Format("2006-01-02")

	word, err := s.repo.GetWordOfTheDay(today)
	if err == nil {
		return word, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return "", err
	}

	// No entry for today — pick a random active word and persist it.
	word, err = s.repo.GetRandomActiveWord()
	if err != nil {
		return "", err
	}

	if err = s.repo.SetWordOfTheDay(word, "", today); err != nil {
		// Non-fatal: still return the word even if we couldn't persist.
		return word, nil
	}
	return word, nil
}

func (s *service) SetWordOfTheDay(word, adminID, date string) error {
	if strings.TrimSpace(word) == "" {
		return ErrWordNotFound
	}
	// Validate date format.
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
	}
	exists, err := s.repo.ActiveWordExists(word)
	if err != nil {
		return err
	}
	if !exists {
		return ErrWordNotFound
	}
	return s.repo.SetWordOfTheDay(word, adminID, date)
}
