package contribution

import (
	"time"

	"github.com/iqbaleff214/kamus-banjar-api/internal/dictionary"
)

// Contribution is an audit-log entry for every action taken on a community word.
type Contribution struct {
	ID            string           `json:"id"`
	WordID        int64            `json:"word_id"`
	Word          *dictionary.Word `json:"word,omitempty"`
	ContributorID string           `json:"contributor_id"`
	ReviewerID    *string          `json:"reviewer_id,omitempty"`
	Action        string           `json:"action"`
	Notes         *string          `json:"notes,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
}

// SubmitRequest is the body for POST /api/v1/contributions
// and for PUT /api/v1/contributions/:id (edit).
type SubmitRequest struct {
	Word        string                     `json:"word"`
	Syllable    string                     `json:"syllables"`
	Alphabet    string                     `json:"alphabet"`
	Meanings    []dictionary.WordMeaning   `json:"meanings"`
	Derivatives []dictionary.WordDerivative `json:"derivatives"`
}

// RejectRequest is the body for PATCH /api/v1/admin/contributions/:id/reject.
type RejectRequest struct {
	Notes string `json:"notes"`
}
