package dictionary

import "time"

// Alphabet model to contain alphabets
type Alphabet struct {
	Letter string `json:"letter"`
	Total  int    `json:"total"`
}

// Word model to contain word definition.
// Source and contributor fields are populated only by the MySQL repository.
// ID and Status are admin-only fields: omitempty keeps them out of public responses.
type Word struct {
	ID            int64            `json:"id,omitempty"`
	Word          string           `json:"word"`
	Syllable      string           `json:"syllables,omitempty"`
	Alphabet      string           `json:"alphabet"`
	Meanings      []WordMeaning    `json:"meanings"`
	Derivatives   []WordDerivative `json:"derivatives,omitempty"`
	Source        string           `json:"source,omitempty"`
	Status        string           `json:"status,omitempty"`
	ContributorID *string          `json:"contributor_id,omitempty"`
	ApprovedBy    *string          `json:"approved_by,omitempty"`
	ApprovedAt    *time.Time       `json:"approved_at,omitempty"`
	CreatedAt     *time.Time       `json:"created_at,omitempty"`
	UpdatedAt     *time.Time       `json:"updated_at,omitempty"`
}

type WordMeaning struct {
	Definitions []WordDefinition `json:"definitions"`
}

type WordDefinition struct {
	Definition   string        `json:"definition,omitempty"`
	Refer        string        `json:"refer,omitempty"`
	PartOfSpeech string        `json:"partOfSpeech"`
	Examples     []WordExample `json:"examples,omitempty"`
}

type WordExample struct {
	Bjn string `json:"bjn,omitempty"`
	Id  string `json:"id,omitempty"`
}

type WordDerivative struct {
	Word        string           `json:"word"`
	Syllable    string           `json:"syllables,omitempty"`
	Definitions []WordDefinition `json:"definitions"`
}

type AlphabeticWord map[string][]Word

type SearchResult struct {
	Search string   `json:"search"`
	Total  int      `json:"total"`
	Words  []string `json:"words"`
}
