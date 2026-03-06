package store

import (
	"fmt"
	"time"

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

// GetLastCollectionTime returns the last time a repository was collected.
// Returns zero time if the repository has never been collected.
func (s *Store) GetLastCollectionTime(repository string) (time.Time, error) {
	var lastCollectedAt time.Time
	query := `SELECT last_collected_at FROM collection_metadata WHERE repository = ?`
	err := s.db.Get(&lastCollectedAt, query, repository)
	if err != nil {
		// No row found means this is the first collection
		if err.Error() == "sql: no rows in result set" {
			return time.Time{}, nil
		}
		return time.Time{}, fmt.Errorf("failed to get last collection time: %w", err)
	}
	return lastCollectedAt, nil
}

// UpdateLastCollectionTime upserts the last collection timestamp for a repository.
func (s *Store) UpdateLastCollectionTime(repository string, timestamp time.Time) error {
	query := `
		INSERT INTO collection_metadata (repository, last_collected_at)
		VALUES (?, ?)
		ON CONFLICT(repository) DO UPDATE SET
			last_collected_at = excluded.last_collected_at
	`
	_, err := s.db.Exec(query, repository, timestamp)
	if err != nil {
		return fmt.Errorf("failed to update last collection time: %w", err)
	}
	return nil
}

// UpsertCommitMetric inserts or updates a Commit metric (idempotent)
func (s *Store) UpsertCommitMetric(metric *database.CommitMetric) error {
	query := `
		INSERT INTO commit_metrics (
			team_id, repository, commit_hash, author, message, created_at, created_date
		) VALUES (
			?, ?, ?, ?, ?, ?, ?
		)
		ON CONFLICT(team_id, repository, commit_hash) DO UPDATE SET
			message = excluded.message,
			author = excluded.author
	`

	_, err := s.db.Exec(query,
		metric.TeamID, metric.Repository, metric.CommitHash, metric.Author,
		metric.Message, metric.CreatedAt, metric.CreatedDate,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert commit metric: %w", err)
	}

	return nil
}

// UpsertCommentMetric inserts or updates a Comment metric (idempotent)
func (s *Store) UpsertCommentMetric(metric *database.CommentMetric) error {
	query := `
		INSERT INTO comment_metrics (
			team_id, repository, comment_id, author, body, created_at, created_date, comment_type
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)
		ON CONFLICT(team_id, repository, comment_id, comment_type) DO UPDATE SET
			body = excluded.body,
			author = excluded.author
	`

	_, err := s.db.Exec(query,
		metric.TeamID, metric.Repository, metric.CommentID, metric.Author,
		metric.Body, metric.CreatedAt, metric.CreatedDate, metric.CommentType,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert comment metric: %w", err)
	}

	return nil
}
