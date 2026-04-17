package admin

// WordStats holds word counts broken down by source and status.
type WordStats struct {
	Official  int `json:"official"`
	Community int `json:"community"`
	Pending   int `json:"pending"`
	Rejected  int `json:"rejected"`
}

// UserStats holds user counts broken down by active status.
type UserStats struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
}

// ContributionStats holds contribution activity counts.
type ContributionStats struct {
	Pending   int `json:"pending"`
	ThisWeek  int `json:"this_week"`
	ThisMonth int `json:"this_month"`
}

// TopContributor represents a single entry in the top-contributors list.
type TopContributor struct {
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	ApprovedCount int    `json:"approved_count"`
}

// Stats is the full admin dashboard snapshot.
type Stats struct {
	Words            WordStats          `json:"words"`
	Users            UserStats          `json:"users"`
	Contributions    ContributionStats  `json:"contributions"`
	TopContributors  []TopContributor   `json:"top_contributors"`
}
