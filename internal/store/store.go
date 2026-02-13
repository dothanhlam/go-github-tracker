package store

import (
	"fmt"

	"github.com/dothanhlam/go-github-tracker/internal/database"
)

// Store handles database persistence
type Store struct {
	db *database.DB
}

// New creates a new store
func New(db *database.DB) *Store {
	return &Store{db: db}
}

// UpsertPRMetric inserts or updates a PR metric (idempotent)
func (s *Store) UpsertPRMetric(metric *database.PRMetric) error {
	query := `
		INSERT INTO pr_metrics (
			team_id, pr_number, repository, author, title,
			created_at, merged_at, closed_at, cycle_time_hours, state, created_date,
			first_review_at, review_turnaround_hours,
			review_comments_count, conversation_count,
			changes_requested_count, approved_count,
			reviewers_count, external_reviewers_count, reviewers_list
		) VALUES (
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?,
			?, ?,
			?, ?,
			?, ?, ?
		)
		ON CONFLICT(team_id, repository, pr_number) DO UPDATE SET
			title = excluded.title,
			merged_at = excluded.merged_at,
			closed_at = excluded.closed_at,
			cycle_time_hours = excluded.cycle_time_hours,
			state = excluded.state,
			first_review_at = excluded.first_review_at,
			review_turnaround_hours = excluded.review_turnaround_hours,
			review_comments_count = excluded.review_comments_count,
			conversation_count = excluded.conversation_count,
			changes_requested_count = excluded.changes_requested_count,
			approved_count = excluded.approved_count,
			reviewers_count = excluded.reviewers_count,
			external_reviewers_count = excluded.external_reviewers_count,
			reviewers_list = excluded.reviewers_list
	`

	_, err := s.db.Exec(query,
		metric.TeamID, metric.PRNumber, metric.Repository, metric.Author, metric.Title,
		metric.CreatedAt, metric.MergedAt, metric.ClosedAt, metric.CycleTimeHours, metric.State, metric.CreatedAt,
		metric.FirstReviewAt, metric.ReviewTurnaroundHours,
		metric.ReviewCommentsCount, metric.ConversationCount,
		metric.ChangesRequestedCount, metric.ApprovedCount,
		metric.ReviewersCount, metric.ExternalReviewersCount, metric.ReviewersList,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert PR metric: %w", err)
	}

	return nil
}
