package database

import (
	"time"
)

// Team represents a team in the system
type Team struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TeamMembership represents a user's membership in a team
type TeamMembership struct {
	ID               int        `db:"id"`
	TeamID           int        `db:"team_id"`
	GitHubUsername   string     `db:"github_username"`
	AllocationWeight float64    `db:"allocation_weight"`
	JoinedAt         time.Time  `db:"joined_at"`
	LeftAt           *time.Time `db:"left_at"`
	CreatedAt        time.Time  `db:"created_at"`
}

// PRMetric represents Pull Request metrics
type PRMetric struct {
	ID              int        `db:"id"`
	TeamID          int        `db:"team_id"`
	PRNumber        int        `db:"pr_number"`
	Repository      string     `db:"repository"`
	Author          string     `db:"author"`
	Title           string     `db:"title"`
	CreatedAt       time.Time  `db:"created_at"`
	MergedAt        *time.Time `db:"merged_at"`
	ClosedAt        *time.Time `db:"closed_at"`
	CycleTimeHours  *int       `db:"cycle_time_hours"`
	State           string     `db:"state"`
	CreatedDate     *time.Time `db:"created_date"`
}

// TeamVelocity represents the view_team_velocity view
type TeamVelocity struct {
	TeamID             int     `db:"team_id"`
	Week               string  `db:"week"`
	PRsMerged          int     `db:"prs_merged"`
	AvgCycleTimeHours  float64 `db:"avg_cycle_time_hours"`
}

// DORALeadTime represents the view_dora_lead_time view
type DORALeadTime struct {
	TeamID              int     `db:"team_id"`
	Month               string  `db:"month"`
	MedianLeadTimeHours float64 `db:"median_lead_time_hours"`
	P95LeadTimeHours    float64 `db:"p95_lead_time_hours"`
	PRCount             int     `db:"pr_count"`
}
