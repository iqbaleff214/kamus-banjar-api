package contribution

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/iqbaleff214/kamus-banjar-api/internal/dictionary"
)

type mysqlRepository struct {
	db *sql.DB
}

// ─────────────────────────────────────────────────────────────
// Word writes
// ─────────────────────────────────────────────────────────────

func (r *mysqlRepository) WordExists(word string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM words WHERE word = ?`, word).Scan(&count)
	return count > 0, err
}

func (r *mysqlRepository) CreateWord(w dictionary.Word, contributorID, source, status string) (int64, error) {
	data, err := serializeWordData(w)
	if err != nil {
		return 0, err
	}

	// Ensure the letter exists in the letters table.
	if _, err = r.db.Exec(`INSERT IGNORE INTO letters (letter) VALUES (?)`, w.Alphabet); err != nil {
		return 0, err
	}

	var contribID any
	if contributorID != "" {
		contribID = contributorID
	}

	res, err := r.db.Exec(
		`INSERT INTO words (word, letter, source, status, contributor_id, data) VALUES (?, ?, ?, ?, ?, ?)`,
		w.Word, w.Alphabet, source, status, contribID, data,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *mysqlRepository) UpdateWord(wordID int64, w dictionary.Word) error {
	data, err := serializeWordData(w)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		`UPDATE words SET word = ?, letter = ?, data = ? WHERE id = ?`,
		w.Word, w.Alphabet, data, wordID,
	)
	return err
}

func (r *mysqlRepository) SetWordStatus(wordID int64, status string, reviewerID *string, approvedAt *time.Time) error {
	_, err := r.db.Exec(
		`UPDATE words SET status = ?, approved_by = ?, approved_at = ? WHERE id = ?`,
		status, reviewerID, approvedAt, wordID,
	)
	return err
}

func (r *mysqlRepository) GetWordByID(wordID int64) (dictionary.Word, string, error) {
	var (
		data          string
		status        string
		source        string
		contributorID sql.NullString
		approvedBy    sql.NullString
		approvedAt    sql.NullTime
		createdAt     sql.NullTime
		updatedAt     sql.NullTime
	)

	err := r.db.QueryRow(
		`SELECT data, status, source, contributor_id, approved_by, approved_at, created_at, updated_at
		 FROM words WHERE id = ?`,
		wordID,
	).Scan(&data, &status, &source, &contributorID, &approvedBy, &approvedAt, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return dictionary.Word{}, "", ErrWordNotFound
	}
	if err != nil {
		return dictionary.Word{}, "", err
	}

	var w dictionary.Word
	if err = json.Unmarshal([]byte(data), &w); err != nil {
		return dictionary.Word{}, "", err
	}
	w.ID = wordID
	w.Status = status
	w.Source = source
	if contributorID.Valid {
		w.ContributorID = &contributorID.String
	}
	if approvedBy.Valid {
		w.ApprovedBy = &approvedBy.String
	}
	if approvedAt.Valid {
		w.ApprovedAt = &approvedAt.Time
	}
	if createdAt.Valid {
		w.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		w.UpdatedAt = &updatedAt.Time
	}
	return w, status, nil
}

// ─────────────────────────────────────────────────────────────
// Contribution log
// ─────────────────────────────────────────────────────────────

func (r *mysqlRepository) LogAction(c Contribution) error {
	_, err := r.db.Exec(
		`INSERT INTO contributions (id, word_id, contributor_id, reviewer_id, action, notes) VALUES (?, ?, ?, ?, ?, ?)`,
		c.ID, c.WordID, c.ContributorID, c.ReviewerID, c.Action, c.Notes,
	)
	return err
}

func (r *mysqlRepository) GetByID(id string) (Contribution, error) {
	var (
		c           Contribution
		reviewerID  sql.NullString
		notes       sql.NullString
		wordData    string
		wordStatus  string
		wordSource  string
	)
	err := r.db.QueryRow(
		`SELECT c.id, c.word_id, c.contributor_id, c.reviewer_id, c.action, c.notes, c.created_at,
		        w.data, w.status, w.source
		 FROM contributions c
		 JOIN words w ON c.word_id = w.id
		 WHERE c.id = ?`,
		id,
	).Scan(
		&c.ID, &c.WordID, &c.ContributorID, &reviewerID, &c.Action, &notes, &c.CreatedAt,
		&wordData, &wordStatus, &wordSource,
	)
	if err == sql.ErrNoRows {
		return Contribution{}, ErrNotFound
	}
	if err != nil {
		return Contribution{}, err
	}
	if reviewerID.Valid {
		c.ReviewerID = &reviewerID.String
	}
	if notes.Valid {
		c.Notes = &notes.String
	}

	var w dictionary.Word
	if err = json.Unmarshal([]byte(wordData), &w); err == nil {
		w.ID = c.WordID
		w.Status = wordStatus
		w.Source = wordSource
		c.Word = &w
	}
	return c, nil
}

func (r *mysqlRepository) ListByContributor(contributorID string, page, limit int) ([]Contribution, int, error) {
	offset := (page - 1) * limit

	var total int
	if err := r.db.QueryRow(
		`SELECT COUNT(*) FROM contributions WHERE contributor_id = ?`, contributorID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		`SELECT c.id, c.word_id, c.contributor_id, c.reviewer_id, c.action, c.notes, c.created_at, w.word
		 FROM contributions c
		 JOIN words w ON c.word_id = w.id
		 WHERE c.contributor_id = ?
		 ORDER BY c.created_at DESC
		 LIMIT ? OFFSET ?`,
		contributorID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []Contribution
	for rows.Next() {
		var (
			c          Contribution
			reviewerID sql.NullString
			notes      sql.NullString
			wordText   string
		)
		if err = rows.Scan(&c.ID, &c.WordID, &c.ContributorID, &reviewerID, &c.Action, &notes, &c.CreatedAt, &wordText); err != nil {
			return nil, 0, err
		}
		if reviewerID.Valid {
			c.ReviewerID = &reviewerID.String
		}
		if notes.Valid {
			c.Notes = &notes.String
		}
		c.Word = &dictionary.Word{Word: wordText}
		list = append(list, c)
	}
	if list == nil {
		list = []Contribution{}
	}
	return list, total, rows.Err()
}

func (r *mysqlRepository) ListAll(status string, page, limit int) ([]Contribution, int, error) {
	offset := (page - 1) * limit

	whereArgs := []any{}
	where := `WHERE 1=1`
	if status != "" {
		where += ` AND c.action = ?`
		whereArgs = append(whereArgs, status)
	}

	var total int
	if err := r.db.QueryRow(
		`SELECT COUNT(*) FROM contributions c `+where, whereArgs...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := append(whereArgs, limit, offset)
	rows, err := r.db.Query(
		`SELECT c.id, c.word_id, c.contributor_id, c.reviewer_id, c.action, c.notes, c.created_at, w.word
		 FROM contributions c
		 JOIN words w ON c.word_id = w.id `+where+`
		 ORDER BY c.created_at DESC
		 LIMIT ? OFFSET ?`,
		queryArgs...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []Contribution
	for rows.Next() {
		var (
			c          Contribution
			reviewerID sql.NullString
			notes      sql.NullString
			wordText   string
		)
		if err = rows.Scan(&c.ID, &c.WordID, &c.ContributorID, &reviewerID, &c.Action, &notes, &c.CreatedAt, &wordText); err != nil {
			return nil, 0, err
		}
		if reviewerID.Valid {
			c.ReviewerID = &reviewerID.String
		}
		if notes.Valid {
			c.Notes = &notes.String
		}
		c.Word = &dictionary.Word{Word: wordText}
		list = append(list, c)
	}
	if list == nil {
		list = []Contribution{}
	}
	return list, total, rows.Err()
}

// ─────────────────────────────────────────────────────────────
// Admin word queries
// ─────────────────────────────────────────────────────────────

func (r *mysqlRepository) ListWords(status, source string, page, limit int) ([]dictionary.Word, int, error) {
	offset := (page - 1) * limit

	where := `WHERE 1=1`
	args := []any{}
	if status != "" {
		where += ` AND status = ?`
		args = append(args, status)
	}
	if source != "" {
		where += ` AND source = ?`
		args = append(args, source)
	}

	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM words `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := append(args, limit, offset)
	rows, err := r.db.Query(
		`SELECT id, data, source, status, contributor_id, approved_by, approved_at, created_at, updated_at
		 FROM words `+where+` ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		queryArgs...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []dictionary.Word
	for rows.Next() {
		var (
			wordID        int64
			data          string
			source        string
			status        string
			contributorID sql.NullString
			approvedBy    sql.NullString
			approvedAt    sql.NullTime
			createdAt     sql.NullTime
			updatedAt     sql.NullTime
		)
		if err = rows.Scan(&wordID, &data, &source, &status, &contributorID, &approvedBy, &approvedAt, &createdAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		var w dictionary.Word
		if err = json.Unmarshal([]byte(data), &w); err != nil {
			return nil, 0, err
		}
		w.ID = wordID
		w.Status = status
		w.Source = source
		if contributorID.Valid {
			w.ContributorID = &contributorID.String
		}
		if approvedBy.Valid {
			w.ApprovedBy = &approvedBy.String
		}
		if approvedAt.Valid {
			w.ApprovedAt = &approvedAt.Time
		}
		if createdAt.Valid {
			w.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			w.UpdatedAt = &updatedAt.Time
		}
		words = append(words, w)
	}
	if words == nil {
		words = []dictionary.Word{}
	}
	return words, total, rows.Err()
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────

func (r *mysqlRepository) DeleteWord(wordID int64) error {
	_, err := r.db.Exec(`DELETE FROM words WHERE id = ?`, wordID)
	return err
}

// serializeWordData marshals only the core linguistic fields into the data column JSON.
// ID, Status, Source, and contributor fields are stored as separate columns.
func serializeWordData(w dictionary.Word) (string, error) {
	core := dictionary.Word{
		Word:        w.Word,
		Syllable:    w.Syllable,
		Alphabet:    w.Alphabet,
		Meanings:    w.Meanings,
		Derivatives: w.Derivatives,
	}
	b, err := json.Marshal(core)
	return string(b), err
}
