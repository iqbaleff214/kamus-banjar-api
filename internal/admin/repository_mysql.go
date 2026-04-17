package admin

import "database/sql"

type mysqlRepository struct {
	db *sql.DB
}

// NewRepository returns a MySQL-backed Repository.
func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetStats() (Stats, error) {
	var s Stats

	if err := r.scanWordStats(&s.Words); err != nil {
		return Stats{}, err
	}
	if err := r.scanUserStats(&s.Users); err != nil {
		return Stats{}, err
	}
	if err := r.scanContributionStats(&s.Contributions); err != nil {
		return Stats{}, err
	}

	contributors, err := r.topContributors()
	if err != nil {
		return Stats{}, err
	}
	s.TopContributors = contributors

	return s, nil
}

func (r *mysqlRepository) scanWordStats(ws *WordStats) error {
	const q = `
		SELECT
			SUM(CASE WHEN source = 'official'  AND status = 'active'   THEN 1 ELSE 0 END),
			SUM(CASE WHEN source = 'community' AND status = 'active'   THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = 'pending'                           THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = 'rejected'                          THEN 1 ELSE 0 END)
		FROM words`
	return r.db.QueryRow(q).Scan(
		&ws.Official, &ws.Community, &ws.Pending, &ws.Rejected,
	)
}

func (r *mysqlRepository) scanUserStats(us *UserStats) error {
	const q = `
		SELECT
			COUNT(*),
			SUM(CASE WHEN is_active = 1 THEN 1 ELSE 0 END),
			SUM(CASE WHEN is_active = 0 THEN 1 ELSE 0 END)
		FROM users`
	return r.db.QueryRow(q).Scan(&us.Total, &us.Active, &us.Inactive)
}

func (r *mysqlRepository) scanContributionStats(cs *ContributionStats) error {
	const q = `
		SELECT
			SUM(CASE WHEN action = 'submitted'                                          THEN 1 ELSE 0 END),
			SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)                THEN 1 ELSE 0 END),
			SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)              THEN 1 ELSE 0 END)
		FROM contributions`
	return r.db.QueryRow(q).Scan(&cs.Pending, &cs.ThisWeek, &cs.ThisMonth)
}

func (r *mysqlRepository) topContributors() ([]TopContributor, error) {
	const q = `
		SELECT u.id, u.name, COUNT(c.id) AS approved_count
		FROM contributions c
		JOIN users u ON c.contributor_id = u.id
		WHERE c.action = 'approved'
		GROUP BY u.id, u.name
		ORDER BY approved_count DESC
		LIMIT 5`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []TopContributor
	for rows.Next() {
		var tc TopContributor
		if err := rows.Scan(&tc.UserID, &tc.Name, &tc.ApprovedCount); err != nil {
			return nil, err
		}
		list = append(list, tc)
	}
	if list == nil {
		list = []TopContributor{}
	}
	return list, rows.Err()
}
