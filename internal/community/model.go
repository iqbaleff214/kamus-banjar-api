package community

import "time"

// VoteSummary holds aggregated vote counts for a word.
type VoteSummary struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}

// VoteRequest is the body for POST /api/v1/entries/:word/votes.
type VoteRequest struct {
	Vote int `json:"vote"`
}

// Comment represents a user comment on a word entry, with optional nested replies.
type Comment struct {
	ID        string    `json:"id"`
	Word      string    `json:"word"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	ParentID  *string   `json:"parent_id,omitempty"`
	Body      string    `json:"body"`
	Replies   []Comment `json:"replies,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentRequest is the body for POST /api/v1/entries/:word/comments.
type CommentRequest struct {
	Body     string  `json:"body"`
	ParentID *string `json:"parent_id"`
}

// WordOfTheDayRequest is the body for PUT /api/v1/admin/word-of-the-day.
type WordOfTheDayRequest struct {
	Word string `json:"word"`
	Date string `json:"date"` // YYYY-MM-DD
}
