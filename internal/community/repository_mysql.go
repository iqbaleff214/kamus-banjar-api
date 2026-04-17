package community

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type mysqlRepository struct {
	db *sql.DB
}

// ─────────────────────────────────────────────────────────────
// Shared
// ─────────────────────────────────────────────────────────────

func (r *mysqlRepository) ActiveWordExists(word string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM words WHERE word = ? AND status = 'active'`, word).Scan(&count)
	return count > 0, err
}

// ─────────────────────────────────────────────────────────────
// Votes
// ─────────────────────────────────────────────────────────────

func (r *mysqlRepository) GetUserVote(word, userID string) (int, error) {
	var vote int
	err := r.db.QueryRow(`SELECT vote FROM word_votes WHERE word = ? AND user_id = ?`, word, userID).Scan(&vote)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return vote, err
}

func (r *mysqlRepository) UpsertVote(word, userID string, vote int) error {
	_, err := r.db.Exec(
		`INSERT INTO word_votes (id, word, user_id, vote) VALUES (?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE vote = VALUES(vote)`,
		uuid.New().String(), word, userID, vote,
	)
	return err
}

func (r *mysqlRepository) DeleteVote(word, userID string) error {
	_, err := r.db.Exec(`DELETE FROM word_votes WHERE word = ? AND user_id = ?`, word, userID)
	return err
}

func (r *mysqlRepository) GetVoteSummary(word string) (up, down int, err error) {
	rows, err := r.db.Query(
		`SELECT vote, COUNT(*) FROM word_votes WHERE word = ? GROUP BY vote`,
		word,
	)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var v, count int
		if err = rows.Scan(&v, &count); err != nil {
			return 0, 0, err
		}
		if v == 1 {
			up = count
		} else if v == -1 {
			down = count
		}
	}
	return up, down, rows.Err()
}

// ─────────────────────────────────────────────────────────────
// Comments
// ─────────────────────────────────────────────────────────────

func (r *mysqlRepository) CreateComment(c Comment) error {
	_, err := r.db.Exec(
		`INSERT INTO word_comments (id, word, user_id, parent_id, body) VALUES (?, ?, ?, ?, ?)`,
		c.ID, c.Word, c.UserID, c.ParentID, c.Body,
	)
	return err
}

func (r *mysqlRepository) GetComment(id string) (Comment, error) {
	var (
		c        Comment
		parentID sql.NullString
	)
	err := r.db.QueryRow(
		`SELECT id, word, user_id, parent_id, body, created_at, updated_at FROM word_comments WHERE id = ?`,
		id,
	).Scan(&c.ID, &c.Word, &c.UserID, &parentID, &c.Body, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Comment{}, ErrCommentNotFound
	}
	if err != nil {
		return Comment{}, err
	}
	if parentID.Valid {
		c.ParentID = &parentID.String
	}
	return c, nil
}

func (r *mysqlRepository) ListComments(word string) ([]Comment, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.word, c.user_id, u.name, c.parent_id, c.body, c.created_at, c.updated_at
		 FROM word_comments c
		 JOIN users u ON c.user_id = u.id
		 WHERE c.word = ?
		 ORDER BY c.created_at ASC`,
		word,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Comment
	for rows.Next() {
		var (
			c        Comment
			parentID sql.NullString
		)
		if err = rows.Scan(&c.ID, &c.Word, &c.UserID, &c.UserName, &parentID, &c.Body, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			c.ParentID = &parentID.String
		}
		all = append(all, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return buildCommentTree(all), nil
}

func (r *mysqlRepository) DeleteComment(id string) error {
	_, err := r.db.Exec(`DELETE FROM word_comments WHERE id = ?`, id)
	return err
}

// buildCommentTree assembles a flat list into a top-level + replies tree (one level deep).
func buildCommentTree(flat []Comment) []Comment {
	byID := make(map[string]*Comment, len(flat))
	var roots []Comment

	// Build id→comment map using a temporary slice with stable pointers.
	tmp := make([]Comment, len(flat))
	for i, c := range flat {
		tmp[i] = c
		byID[c.ID] = &tmp[i]
	}

	for i := range tmp {
		c := &tmp[i]
		if c.ParentID == nil {
			roots = append(roots, *c)
		} else if parent, ok := byID[*c.ParentID]; ok {
			parent.Replies = append(parent.Replies, *c)
		}
	}

	// Sync parent mutations back into roots slice.
	rootIdx := 0
	for i := range tmp {
		if tmp[i].ParentID == nil {
			if rootIdx < len(roots) {
				roots[rootIdx] = tmp[i]
				rootIdx++
			}
		}
	}

	return roots
}

// ─────────────────────────────────────────────────────────────
// Bookmarks
// ─────────────────────────────────────────────────────────────

func (r *mysqlRepository) AddBookmark(word, userID string) error {
	_, err := r.db.Exec(
		`INSERT IGNORE INTO word_bookmarks (id, word, user_id) VALUES (?, ?, ?)`,
		uuid.New().String(), word, userID,
	)
	return err
}

func (r *mysqlRepository) RemoveBookmark(word, userID string) error {
	_, err := r.db.Exec(`DELETE FROM word_bookmarks WHERE word = ? AND user_id = ?`, word, userID)
	return err
}

func (r *mysqlRepository) ListBookmarks(userID string, page, limit int) ([]string, int, error) {
	offset := (page - 1) * limit

	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM word_bookmarks WHERE user_id = ?`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		`SELECT word FROM word_bookmarks WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var words []string
	for rows.Next() {
		var w string
		if err = rows.Scan(&w); err != nil {
			return nil, 0, err
		}
		words = append(words, w)
	}
	if words == nil {
		words = []string{}
	}
	return words, total, rows.Err()
}

// ─────────────────────────────────────────────────────────────
// Word of the day
// ─────────────────────────────────────────────────────────────

func (r *mysqlRepository) GetWordOfTheDay(date string) (string, error) {
	var word string
	err := r.db.QueryRow(`SELECT word FROM word_of_the_day WHERE date = ?`, date).Scan(&word)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return word, err
}

func (r *mysqlRepository) SetWordOfTheDay(word, setBy, date string) error {
	var setByVal any
	if setBy != "" {
		setByVal = setBy
	}
	_, err := r.db.Exec(
		`INSERT INTO word_of_the_day (id, word, date, set_by) VALUES (?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE word = VALUES(word), set_by = VALUES(set_by)`,
		uuid.New().String(), word, date, setByVal,
	)
	return err
}

func (r *mysqlRepository) GetRandomActiveWord() (string, error) {
	var word string
	err := r.db.QueryRow(`SELECT word FROM words WHERE status = 'active' ORDER BY RAND() LIMIT 1`).Scan(&word)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoWordToday
	}
	return word, err
}
