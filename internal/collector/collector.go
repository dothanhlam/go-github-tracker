package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dothanhlam/go-github-tracker/internal/config"
	"github.com/dothanhlam/go-github-tracker/internal/database"
	"github.com/dothanhlam/go-github-tracker/internal/github"
	"github.com/dothanhlam/go-github-tracker/internal/store"
	"github.com/dothanhlam/go-github-tracker/internal/team"
	gh "github.com/google/go-github/v58/github"
)

// Collector orchestrates PR data collection
type Collector struct {
	github  *github.Client
	teamMgr *team.Manager
	store   *store.Store
	config  *config.Config
}

// New creates a new collector
func New(cfg *config.Config, db *database.DB) (*Collector, error) {
	// Create GitHub client
	ghClient := github.NewClient(cfg.GitHubPAT)

	// Create team manager
	teamMgr, err := team.NewManager(db, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create team manager: %w", err)
	}

	// Create store
	st := store.New(db)

	return &Collector{
		github:  ghClient,
		teamMgr: teamMgr,
		store:   st,
		config:  cfg,
	}, nil
}

// Run executes the collection process
func (c *Collector) Run() error {
	fmt.Printf("📊 Teams configured: %d\n", len(c.teamMgr.GetAllTeamIDs()))
	fmt.Printf("📦 Repositories to track: %d\n\n", len(c.config.Repositories))

	totalPRs := 0
	for _, repoFullName := range c.config.Repositories {
		parts := strings.Split(repoFullName, "/")
		if len(parts) != 2 {
			fmt.Printf("⚠️  Invalid repository format: %s (expected owner/repo)\n", repoFullName)
			continue
		}

		owner, repo := parts[0], parts[1]
		count, err := c.collectRepository(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to collect %s/%s: %w", owner, repo, err)
		}
		totalPRs += count
	}

	fmt.Printf("\n✅ Collection complete! Processed %d PRs\n", totalPRs)
	return nil
}

// collectRepository collects PRs from a single repository
func (c *Collector) collectRepository(owner, repo string) (int, error) {
	fmt.Printf("🔄 Processing repository: %s/%s\n", owner, repo)

	// Calculate lookback date
	since := time.Now().AddDate(0, 0, -c.config.LookbackDays)
	fmt.Printf("  📅 Collecting PRs updated since: %s (%d days lookback)\n", 
		since.Format("2006-01-02"), c.config.LookbackDays)

	// Fetch PRs with date filter
	prs, err := c.github.FetchPRs(owner, repo, since)
	if err != nil {
		return 0, err
	}

	processedCount := 0
	for i, pr := range prs {
		if (i+1)%10 == 0 {
			fmt.Printf("  ⏳ Processing PR %d/%d...\n", i+1, len(prs))
		}

		// Fetch reviews and comments
		reviews, err := c.github.FetchReviews(owner, repo, pr.GetNumber())
		if err != nil {
			fmt.Printf("  ⚠️  Failed to fetch reviews for PR #%d: %v\n", pr.GetNumber(), err)
			continue
		}

		comments, err := c.github.FetchComments(owner, repo, pr.GetNumber())
		if err != nil {
			fmt.Printf("  ⚠️  Failed to fetch comments for PR #%d: %v\n", pr.GetNumber(), err)
			continue
		}

		// Check if PR involves team members
		if !c.shouldIncludePR(pr, reviews) {
			continue
		}

		// Process PR for each relevant team
		teams := c.getRelevantTeams(pr, reviews)
		for _, teamID := range teams {
			metric := c.processPR(pr, reviews, comments, teamID, fmt.Sprintf("%s/%s", owner, repo))
			if err := c.store.UpsertPRMetric(metric); err != nil {
				fmt.Printf("  ⚠️  Failed to store PR #%d: %v\n", pr.GetNumber(), err)
				continue
			}
			processedCount++
		}
	}

	fmt.Printf("  ✓ Processed %d PRs for team members\n", processedCount)
	return processedCount, nil
}

// shouldIncludePR checks if PR involves any team member
func (c *Collector) shouldIncludePR(pr *gh.PullRequest, reviews []*gh.PullRequestReview) bool {
	// Check if author is team member
	if c.teamMgr.IsMember(pr.GetUser().GetLogin()) {
		return true
	}

	// Check if any reviewer is team member
	for _, review := range reviews {
		if c.teamMgr.IsMember(review.GetUser().GetLogin()) {
			return true
		}
	}

	return false
}

// getRelevantTeams returns team IDs that should track this PR
func (c *Collector) getRelevantTeams(pr *gh.PullRequest, reviews []*gh.PullRequestReview) []int {
	teamsMap := make(map[int]bool)

	// Add teams for author
	for _, teamID := range c.teamMgr.GetTeamsForUser(pr.GetUser().GetLogin()) {
		teamsMap[teamID] = true
	}

	// Add teams for reviewers
	for _, review := range reviews {
		for _, teamID := range c.teamMgr.GetTeamsForUser(review.GetUser().GetLogin()) {
			teamsMap[teamID] = true
		}
	}

	// Convert map to slice
	var teams []int
	for teamID := range teamsMap {
		teams = append(teams, teamID)
	}

	return teams
}

// processPR converts GitHub PR to database metric
func (c *Collector) processPR(
	pr *gh.PullRequest,
	reviews []*gh.PullRequestReview,
	comments []*gh.PullRequestComment,
	teamID int,
	repository string,
) *database.PRMetric {
	metric := &database.PRMetric{
		TeamID:     teamID,
		PRNumber:   pr.GetNumber(),
		Repository: repository,
		Author:     pr.GetUser().GetLogin(),
		Title:      pr.GetTitle(),
		CreatedAt:  pr.GetCreatedAt().Time,
		State:      pr.GetState(),
	}

	// Set merged/closed timestamps
	if pr.MergedAt != nil {
		mergedAt := pr.GetMergedAt().Time
		metric.MergedAt = &mergedAt
		metric.State = "merged"

		// Calculate cycle time
		cycleTime := calculateCycleTime(metric.CreatedAt, mergedAt)
		metric.CycleTimeHours = &cycleTime
	}

	if pr.ClosedAt != nil {
		closedAt := pr.GetClosedAt().Time
		metric.ClosedAt = &closedAt
	}

	// Process reviews
	if len(reviews) > 0 {
		firstReviewAt := getFirstReviewTime(reviews)
		metric.FirstReviewAt = &firstReviewAt

		turnaround := calculateReviewTurnaround(metric.CreatedAt, firstReviewAt)
		metric.ReviewTurnaroundHours = &turnaround

		metric.ChangesRequestedCount = countReviewsByState(reviews, "CHANGES_REQUESTED")
		metric.ApprovedCount = countReviewsByState(reviews, "APPROVED")

		reviewers := extractReviewers(reviews)
		metric.ReviewersCount = len(reviewers)
		metric.ExternalReviewersCount = c.countExternalReviewers(reviewers, teamID)

		reviewersJSON, _ := json.Marshal(reviewers)
		metric.ReviewersList = string(reviewersJSON)
	}

	// Process comments
	metric.ReviewCommentsCount = len(comments)
	metric.ConversationCount = countConversations(comments)

	return metric
}
