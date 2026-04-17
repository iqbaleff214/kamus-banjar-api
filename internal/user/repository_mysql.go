package user

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type mysqlRepository struct {
	db *sql.DB
}

func (r *mysqlRepository) Create(u User, passwordHash string) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, name, email, password, role, is_active) VALUES (?, ?, ?, ?, ?, ?)`,
		u.ID, u.Name, u.Email, passwordHash, u.Role, u.IsActive,
	)
	return err
}

func (r *mysqlRepository) FindByEmail(email string) (User, string, error) {
	var u User
	var hash string
	err := r.db.QueryRow(
		`SELECT id, name, email, role, is_active, created_at, updated_at, password
		 FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &hash)
	if err == sql.ErrNoRows {
		return User{}, "", ErrNotFound
	}
	return u, hash, err
}

func (r *mysqlRepository) FindByID(id string) (User, error) {
	var u User
	err := r.db.QueryRow(
		`SELECT id, name, email, role, is_active, created_at, updated_at
		 FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return User{}, ErrNotFound
	}
	return u, err
}

func (r *mysqlRepository) Update(id, name string) error {
	_, err := r.db.Exec(`UPDATE users SET name = ? WHERE id = ?`, name, id)
	return err
}

func (r *mysqlRepository) UpdatePassword(id, hash string) error {
	_, err := r.db.Exec(`UPDATE users SET password = ? WHERE id = ?`, hash, id)
	return err
}

func (r *mysqlRepository) SaveRefreshToken(userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at) VALUES (?, ?, ?, ?)`,
		uuid.New().String(), userID, tokenHash, expiresAt,
	)
	return err
}

func (r *mysqlRepository) RevokeRefreshToken(tokenHash string) error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE token_hash = ?`, tokenHash)
	return err
}

func (r *mysqlRepository) FindRefreshToken(tokenHash string) (string, error) {
	var userID string
	var expiresAt time.Time
	err := r.db.QueryRow(
		`SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = ?`,
		tokenHash,
	).Scan(&userID, &expiresAt)
	if err == sql.ErrNoRows {
		return "", ErrTokenInvalid
	}
	if err != nil {
		return "", err
	}
	if time.Now().After(expiresAt) {
		r.db.Exec(`DELETE FROM refresh_tokens WHERE token_hash = ?`, tokenHash)
		return "", ErrTokenInvalid
	}
	return userID, nil
}

func (r *mysqlRepository) ListUsers(page, limit int, role string, active *bool) ([]User, int, error) {
	offset := (page - 1) * limit

	where := `WHERE 1=1`
	args := []any{}

	if role != "" {
		where += ` AND role = ?`
		args = append(args, role)
	}
	if active != nil {
		where += ` AND is_active = ?`
		args = append(args, *active)
	}

	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM users `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	rows, err := r.db.Query(
		`SELECT id, name, email, role, is_active, created_at, updated_at
		 FROM users `+where+` ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, rows.Err()
}

func (r *mysqlRepository) SetActive(id string, active bool) error {
	_, err := r.db.Exec(`UPDATE users SET is_active = ? WHERE id = ?`, active, id)
	return err
}

func (r *mysqlRepository) SetRole(id, role string) error {
	_, err := r.db.Exec(`UPDATE users SET role = ? WHERE id = ?`, role, id)
	return err
}
