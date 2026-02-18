package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dothanhlam/go-github-tracker/internal/database"
)

// MetricsService handles business logic for metrics queries
type MetricsService struct {
	db *database.DB
}

// NewMetricsService creates a new metrics service
func NewMetricsService(db *database.DB) *MetricsService {
	return &MetricsService{db: db}
}

// Team represents a team
type Team struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"member_count"`
}

// VelocityMetric represents team velocity for a period
type VelocityMetric struct {
	Period           string  `json:"period"`
	PRsMerged        int     `json:"prs_merged"`
	AvgCycleTimeHrs  float64 `json:"avg_cycle_time_hours"`
}

// VelocityResponse represents the API response for velocity
type VelocityResponse struct {
	TeamID   int              `json:"team_id"`
	TeamName string           `json:"team_name"`
	Period   Period           `json:"period"`
	Metrics  []VelocityMetric `json:"metrics"`
}

// LeadTimeMetric represents DORA lead time for a period
type LeadTimeMetric struct {
	Period              string  `json:"period"`
	MedianLeadTimeHrs   float64 `json:"median_lead_time_hours"`
	P95LeadTimeHrs      float64 `json:"p95_lead_time_hours"`
}

// LeadTimeResponse represents the API response for lead time
type LeadTimeResponse struct {
	TeamID   int              `json:"team_id"`
	TeamName string           `json:"team_name"`
	Period   Period           `json:"period"`
	Metrics  []LeadTimeMetric `json:"metrics"`
}

// ReviewTurnaroundMetric represents review turnaround metrics
type ReviewTurnaroundMetric struct {
	Period                 string  `json:"period"`
	AvgTurnaroundHrs       float64 `json:"avg_turnaround_hours"`
	MedianTurnaroundHrs    float64 `json:"median_turnaround_hours"`
}

// ReviewTurnaroundResponse represents the API response for review turnaround
type ReviewTurnaroundResponse struct {
	TeamID   int                      `json:"team_id"`
	TeamName string                   `json:"team_name"`
	Period   Period                   `json:"period"`
	Metrics  []ReviewTurnaroundMetric `json:"metrics"`
}

// ReviewEngagementMetric represents review engagement metrics
type ReviewEngagementMetric struct {
	Period           string  `json:"period"`
	TotalReviews     int     `json:"total_reviews"`
	UniqueReviewers  int     `json:"unique_reviewers"`
	AvgReviewsPerPR  float64 `json:"avg_reviews_per_pr"`
}

// ReviewEngagementResponse represents the API response for review engagement
type ReviewEngagementResponse struct {
	TeamID   int                      `json:"team_id"`
	TeamName string                   `json:"team_name"`
	Period   Period                   `json:"period"`
	Metrics  []ReviewEngagementMetric `json:"metrics"`
}

// KnowledgeSharingMetric represents knowledge sharing metrics
type KnowledgeSharingMetric struct {
	Period               string  `json:"period"`
	CrossTeamReviews     int     `json:"cross_team_reviews"`
	KnowledgeSharingScore float64 `json:"knowledge_sharing_score"`
}

// KnowledgeSharingResponse represents the API response for knowledge sharing
type KnowledgeSharingResponse struct {
	TeamID   int                      `json:"team_id"`
	TeamName string                   `json:"team_name"`
	Period   Period                   `json:"period"`
	Metrics  []KnowledgeSharingMetric `json:"metrics"`
}

// Period represents a time period
type Period struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// ListTeams returns all teams
func (s *MetricsService) ListTeams() ([]Team, error) {
	query := `
		SELECT 
			t.id,
			t.name,
			COUNT(DISTINCT tm.github_username) as member_count
		FROM teams t
		LEFT JOIN team_memberships tm ON t.id = tm.team_id
		GROUP BY t.id, t.name
		ORDER BY t.name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query teams: %w", err)
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var team Team
		if err := rows.Scan(&team.ID, &team.Name, &team.MemberCount); err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}
		teams = append(teams, team)
	}

	return teams, nil
}

// GetTeamVelocity returns velocity metrics for a team
func (s *MetricsService) GetTeamVelocity(teamID int, startDate, endDate time.Time, granularity string) (*VelocityResponse, error) {
	// Get team name
	var teamName string
	err := s.db.QueryRow("SELECT name FROM teams WHERE id = ?", teamID).Scan(&teamName)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Query velocity metrics
	query := `
		SELECT 
			week_start,
			prs_merged,
			avg_cycle_time_hours
		FROM view_team_velocity
		WHERE team_id = ?
			AND week_start >= ?
			AND week_start <= ?
		ORDER BY week_start
	`

	rows, err := s.db.Query(query, teamID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query velocity: %w", err)
	}
	defer rows.Close()

	var metrics []VelocityMetric
	for rows.Next() {
		var metric VelocityMetric
		var cycleTime sql.NullFloat64
		if err := rows.Scan(&metric.Period, &metric.PRsMerged, &cycleTime); err != nil {
			return nil, fmt.Errorf("failed to scan velocity: %w", err)
		}
		if cycleTime.Valid {
			metric.AvgCycleTimeHrs = cycleTime.Float64
		}
		metrics = append(metrics, metric)
	}

	return &VelocityResponse{
		TeamID:   teamID,
		TeamName: teamName,
		Period: Period{
			Start: startDate.Format(time.RFC3339),
			End:   endDate.Format(time.RFC3339),
		},
		Metrics: metrics,
	}, nil
}

// GetTeamLeadTime returns DORA lead time metrics for a team
func (s *MetricsService) GetTeamLeadTime(teamID int, startDate, endDate time.Time, granularity string) (*LeadTimeResponse, error) {
	// Get team name
	var teamName string
	err := s.db.QueryRow("SELECT name FROM teams WHERE id = ?", teamID).Scan(&teamName)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Query lead time metrics
	query := `
		SELECT 
			month,
			median_lead_time_hours,
			p95_lead_time_hours
		FROM view_dora_lead_time
		WHERE team_id = ?
			AND month >= ?
			AND month <= ?
		ORDER BY month
	`

	rows, err := s.db.Query(query, teamID, startDate.Format("2006-01"), endDate.Format("2006-01"))
	if err != nil {
		return nil, fmt.Errorf("failed to query lead time: %w", err)
	}
	defer rows.Close()

	var metrics []LeadTimeMetric
	for rows.Next() {
		var metric LeadTimeMetric
		var median, p95 sql.NullFloat64
		if err := rows.Scan(&metric.Period, &median, &p95); err != nil {
			return nil, fmt.Errorf("failed to scan lead time: %w", err)
		}
		if median.Valid {
			metric.MedianLeadTimeHrs = median.Float64
		}
		if p95.Valid {
			metric.P95LeadTimeHrs = p95.Float64
		}
		metrics = append(metrics, metric)
	}

	return &LeadTimeResponse{
		TeamID:   teamID,
		TeamName: teamName,
		Period: Period{
			Start: startDate.Format(time.RFC3339),
			End:   endDate.Format(time.RFC3339),
		},
		Metrics: metrics,
	}, nil
}

// GetReviewTurnaround returns review turnaround metrics for a team
func (s *MetricsService) GetReviewTurnaround(teamID int, startDate, endDate time.Time) (*ReviewTurnaroundResponse, error) {
	// Get team name
	var teamName string
	err := s.db.QueryRow("SELECT name FROM teams WHERE id = ?", teamID).Scan(&teamName)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Query review turnaround metrics
	query := `
		SELECT 
			week_start,
			avg_turnaround_hours,
			median_turnaround_hours
		FROM view_review_turnaround
		WHERE team_id = ?
			AND week_start >= ?
			AND week_start <= ?
		ORDER BY week_start
	`

	rows, err := s.db.Query(query, teamID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query review turnaround: %w", err)
	}
	defer rows.Close()

	var metrics []ReviewTurnaroundMetric
	for rows.Next() {
		var metric ReviewTurnaroundMetric
		var avg, median sql.NullFloat64
		if err := rows.Scan(&metric.Period, &avg, &median); err != nil {
			return nil, fmt.Errorf("failed to scan review turnaround: %w", err)
		}
		if avg.Valid {
			metric.AvgTurnaroundHrs = avg.Float64
		}
		if median.Valid {
			metric.MedianTurnaroundHrs = median.Float64
		}
		metrics = append(metrics, metric)
	}

	return &ReviewTurnaroundResponse{
		TeamID:   teamID,
		TeamName: teamName,
		Period: Period{
			Start: startDate.Format(time.RFC3339),
			End:   endDate.Format(time.RFC3339),
		},
		Metrics: metrics,
	}, nil
}

// GetReviewEngagement returns review engagement metrics for a team
func (s *MetricsService) GetReviewEngagement(teamID int, startDate, endDate time.Time) (*ReviewEngagementResponse, error) {
	// Get team name
	var teamName string
	err := s.db.QueryRow("SELECT name FROM teams WHERE id = ?", teamID).Scan(&teamName)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Query review engagement metrics
	query := `
		SELECT 
			week_start,
			total_reviews,
			unique_reviewers,
			avg_reviews_per_pr
		FROM view_review_engagement
		WHERE team_id = ?
			AND week_start >= ?
			AND week_start <= ?
		ORDER BY week_start
	`

	rows, err := s.db.Query(query, teamID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query review engagement: %w", err)
	}
	defer rows.Close()

	var metrics []ReviewEngagementMetric
	for rows.Next() {
		var metric ReviewEngagementMetric
		var avgReviews sql.NullFloat64
		if err := rows.Scan(&metric.Period, &metric.TotalReviews, &metric.UniqueReviewers, &avgReviews); err != nil {
			return nil, fmt.Errorf("failed to scan review engagement: %w", err)
		}
		if avgReviews.Valid {
			metric.AvgReviewsPerPR = avgReviews.Float64
		}
		metrics = append(metrics, metric)
	}

	return &ReviewEngagementResponse{
		TeamID:   teamID,
		TeamName: teamName,
		Period: Period{
			Start: startDate.Format(time.RFC3339),
			End:   endDate.Format(time.RFC3339),
		},
		Metrics: metrics,
	}, nil
}

// GetKnowledgeSharing returns knowledge sharing metrics for a team
func (s *MetricsService) GetKnowledgeSharing(teamID int, startDate, endDate time.Time) (*KnowledgeSharingResponse, error) {
	// Get team name
	var teamName string
	err := s.db.QueryRow("SELECT name FROM teams WHERE id = ?", teamID).Scan(&teamName)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Query knowledge sharing metrics
	query := `
		SELECT 
			week_start,
			cross_team_reviews,
			knowledge_sharing_score
		FROM view_knowledge_sharing
		WHERE team_id = ?
			AND week_start >= ?
			AND week_start <= ?
		ORDER BY week_start
	`

	rows, err := s.db.Query(query, teamID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query knowledge sharing: %w", err)
	}
	defer rows.Close()

	var metrics []KnowledgeSharingMetric
	for rows.Next() {
		var metric KnowledgeSharingMetric
		var score sql.NullFloat64
		if err := rows.Scan(&metric.Period, &metric.CrossTeamReviews, &score); err != nil {
			return nil, fmt.Errorf("failed to scan knowledge sharing: %w", err)
		}
		if score.Valid {
			metric.KnowledgeSharingScore = score.Float64
		}
		metrics = append(metrics, metric)
	}

	return &KnowledgeSharingResponse{
		TeamID:   teamID,
		TeamName: teamName,
		Period: Period{
			Start: startDate.Format(time.RFC3339),
			End:   endDate.Format(time.RFC3339),
		},
		Metrics: metrics,
	}, nil
}
