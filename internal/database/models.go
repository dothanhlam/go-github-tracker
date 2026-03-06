package database

import (
	"time"
)

// Team represents a team in the system
type Team struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"` // Nullable
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
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

	// Review metrics
	FirstReviewAt          *time.Time `db:"first_review_at"`
	ReviewTurnaroundHours  *int       `db:"review_turnaround_hours"`
	ReviewCommentsCount    int        `db:"review_comments_count"`
	ConversationCount      int        `db:"conversation_count"`
	ChangesRequestedCount  int        `db:"changes_requested_count"`
	ApprovedCount          int        `db:"approved_count"`
	ReviewersCount         int        `db:"reviewers_count"`
	ExternalReviewersCount int        `db:"external_reviewers_count"`
	ReviewersList          string     `db:"reviewers_list"` // JSON array
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

// ReviewTurnaround represents the view_review_turnaround view
type ReviewTurnaround struct {
	TeamID             int     `db:"team_id"`
	Month              string  `db:"month"`
	AvgTurnaroundHours float64 `db:"avg_turnaround_hours"`
	MinTurnaroundHours int     `db:"min_turnaround_hours"`
	MaxTurnaroundHours int     `db:"max_turnaround_hours"`
	Within24hCount     int     `db:"within_24h_count"`
	Over24hCount       int     `db:"over_24h_count"`
	PRCount            int     `db:"pr_count"`
}

// ReviewEngagement represents the view_review_engagement view
type ReviewEngagement struct {
	TeamID                int     `db:"team_id"`
	Month                 string  `db:"month"`
	AvgCommentsPerPR      float64 `db:"avg_comments_per_pr"`
	AvgConversationsPerPR float64 `db:"avg_conversations_per_pr"`
	AvgReviewersPerPR     float64 `db:"avg_reviewers_per_pr"`
	ChangesRequestedRate  float64 `db:"changes_requested_rate"`
	ApprovalRate          float64 `db:"approval_rate"`
	PRCount               int     `db:"pr_count"`
}

// KnowledgeSharing represents the view_knowledge_sharing view
type KnowledgeSharing struct {
	TeamID                  int     `db:"team_id"`
	Month                   string  `db:"month"`
	AvgReviewers            float64 `db:"avg_reviewers"`
	AvgExternalReviewers    float64 `db:"avg_external_reviewers"`
	ExternalReviewerRate    float64 `db:"external_reviewer_rate"`
	TotalExternalReviews    int     `db:"total_external_reviews"`
	PRsWithExternalReviews  int     `db:"prs_with_external_reviews"`
	PRCount                 int     `db:"pr_count"`
}

// CommitMetric represents Commit metrics
type CommitMetric struct {
	ID          int        `db:"id"`
	TeamID      int        `db:"team_id"`
	Repository  string     `db:"repository"`
	CommitHash  string     `db:"commit_hash"`
	Author      string     `db:"author"`
	Message     string     `db:"message"`
	CreatedAt   time.Time  `db:"created_at"`
	CreatedDate *time.Time `db:"created_date"`
}

// CommentMetric represents Comment metrics
type CommentMetric struct {
	ID          int        `db:"id"`
	TeamID      int        `db:"team_id"`
	Repository  string     `db:"repository"`
	CommentID   int64      `db:"comment_id"`
	Author      string     `db:"author"`
	Body        string     `db:"body"`
	CreatedAt   time.Time  `db:"created_at"`
	CreatedDate *time.Time `db:"created_date"`
	CommentType string     `db:"comment_type"`
}

// TeamCommitVelocity represents the view_team_commit_velocity view
type TeamCommitVelocity struct {
	TeamID       int    `db:"team_id"`
	Week         string `db:"week"`
	CommitsCount int    `db:"commits_count"`
}

// TeamCommentActivity represents the view_team_comment_activity view
type TeamCommentActivity struct {
	TeamID        int    `db:"team_id"`
	Week          string `db:"week"`
	CommentsCount int    `db:"comments_count"`
	CommentType   string `db:"comment_type"`
}

// MemberCommitVelocity represents the view_member_commit_velocity view
type MemberCommitVelocity struct {
	TeamID         int    `db:"team_id"`
	GitHubUsername string `db:"github_username"`
	Week           string `db:"week"`
	CommitsCount   int    `db:"commits_count"`
}

// MemberCommentActivity represents the view_member_comment_activity view
type MemberCommentActivity struct {
	TeamID         int    `db:"team_id"`
	GitHubUsername string `db:"github_username"`
	Week           string `db:"week"`
	CommentsCount  int    `db:"comments_count"`
	CommentType    string `db:"comment_type"`
}
