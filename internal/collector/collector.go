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
	repoFullName := fmt.Sprintf("%s/%s", owner, repo)
	fmt.Printf("🔄 Processing repository: %s\n", repoFullName)

	// Determine collection window: incremental or initial
	since, collectionType, err := c.getCollectionSince(repoFullName)
	if err != nil {
		return 0, err
	}
	fmt.Printf("  📅 Collection type: %s (since %s)\n",
		collectionType, since.Format("2006-01-02 15:04:05"))

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
			metric := c.processPR(pr, reviews, comments, teamID, repoFullName)
			if err := c.store.UpsertPRMetric(metric); err != nil {
				fmt.Printf("  ⚠️  Failed to store PR #%d: %v\n", pr.GetNumber(), err)
				continue
			}
			processedCount++
		}
	}

	// Also process commits
	if err := c.processRepositoryCommits(owner, repo, since); err != nil {
		fmt.Printf("  ⚠️  Failed to process commits: %v\n", err)
	}

	// Also process comments
	if err := c.processRepositoryComments(owner, repo, since); err != nil {
		fmt.Printf("  ⚠️  Failed to process comments: %v\n", err)
	}

	// Update last collection timestamp after successful run
	if err := c.store.UpdateLastCollectionTime(repoFullName, time.Now()); err != nil {
		fmt.Printf("  ⚠️  Failed to update collection timestamp: %v\n", err)
	} else {
		fmt.Printf("  ✅ Collection timestamp updated\n")
	}

	fmt.Printf("  ✓ Processed %d PRs for team members\n", processedCount)
	return processedCount, nil
}

// getCollectionSince determines the start time for PR collection.
// On first run: uses COLLECTION_LOOKBACK_DAYS (default 7).
// On subsequent runs: uses the last recorded collection timestamp.
func (c *Collector) getCollectionSince(repoFullName string) (time.Time, string, error) {
	lastCollected, err := c.store.GetLastCollectionTime(repoFullName)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("failed to get last collection time: %w", err)
	}

	if lastCollected.IsZero() {
		// First run: fall back to lookback days
		since := time.Now().AddDate(0, 0, -c.config.LookbackDays)
		return since, fmt.Sprintf("initial (%d-day lookback)", c.config.LookbackDays), nil
	}

	// Incremental: collect since last run
	return lastCollected, "incremental", nil
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

// processRepositoryCommits collects and stores commits for a repository
func (c *Collector) processRepositoryCommits(owner, repo string, since time.Time) error {
	repoFullName := fmt.Sprintf("%s/%s", owner, repo)
	commits, err := c.github.FetchCommits(owner, repo, since)
	if err != nil {
		return err
	}

	processedCount := 0
	for _, commit := range commits {
		// author can be nil sometimes
		var author string
		if commit.Author != nil && commit.Author.Login != nil {
			author = *commit.Author.Login
		} else if commit.Commit != nil && commit.Commit.Author != nil && commit.Commit.Author.Name != nil {
			author = *commit.Commit.Author.Name
		} else {
			continue // skip if we can't identify author
		}

		if !c.teamMgr.IsMember(author) {
			continue
		}

		teams := c.teamMgr.GetTeamsForUser(author)
		for _, teamID := range teams {
			var createdAt time.Time
			if commit.Commit != nil && commit.Commit.Author != nil {
				createdAt = commit.Commit.Author.GetDate().Time
			}

			metric := &database.CommitMetric{
				TeamID:      teamID,
				Repository:  repoFullName,
				CommitHash:  commit.GetSHA(),
				Author:      author,
				Message:     commit.Commit.GetMessage(),
				CreatedAt:   createdAt,
				CreatedDate: &createdAt,
			}
			if err := c.store.UpsertCommitMetric(metric); err != nil {
				fmt.Printf("  ⚠️  Failed to store commit %s: %v\n", commit.GetSHA(), err)
				continue
			}
			processedCount++
		}
	}

	fmt.Printf("  ✓ Processed %d commits for team members\n", processedCount)
	return nil
}

// processRepositoryComments collects and stores comments for a repository
func (c *Collector) processRepositoryComments(owner, repo string, since time.Time) error {
	repoFullName := fmt.Sprintf("%s/%s", owner, repo)

	// Issue comments
	issueComments, err := c.github.FetchIssueComments(owner, repo, since)
	if err != nil {
		return err
	}

	processedCount := 0
	for _, comment := range issueComments {
		if comment.User == nil || comment.User.Login == nil {
			continue
		}
		author := *comment.User.Login
		if !c.teamMgr.IsMember(author) {
			continue
		}

		teams := c.teamMgr.GetTeamsForUser(author)
		for _, teamID := range teams {
			createdAt := comment.GetCreatedAt().Time
			metric := &database.CommentMetric{
				TeamID:      teamID,
				Repository:  repoFullName,
				CommentID:   comment.GetID(),
				Author:      author,
				Body:        comment.GetBody(),
				CreatedAt:   createdAt,
				CreatedDate: &createdAt,
				CommentType: "issue",
			}
			if err := c.store.UpsertCommentMetric(metric); err != nil {
				fmt.Printf("  ⚠️  Failed to store issue comment %d: %v\n", comment.GetID(), err)
				continue
			}
			processedCount++
		}
	}

	// Commit comments
	commitComments, err := c.github.FetchCommitComments(owner, repo, since)
	if err != nil {
		return err
	}

	for _, comment := range commitComments {
		if comment.User == nil || comment.User.Login == nil {
			continue
		}
		author := *comment.User.Login
		if !c.teamMgr.IsMember(author) {
			continue
		}

		teams := c.teamMgr.GetTeamsForUser(author)
		for _, teamID := range teams {
			createdAt := comment.GetCreatedAt().Time
			metric := &database.CommentMetric{
				TeamID:      teamID,
				Repository:  repoFullName,
				CommentID:   comment.GetID(),
				Author:      author,
				Body:        comment.GetBody(),
				CreatedAt:   createdAt,
				CreatedDate: &createdAt,
				CommentType: "commit",
			}
			if err := c.store.UpsertCommentMetric(metric); err != nil {
				fmt.Printf("  ⚠️  Failed to store commit comment %d: %v\n", comment.GetID(), err)
				continue
			}
			processedCount++
		}
	}

	fmt.Printf("  ✓ Processed %d overall comments for team members\n", processedCount)
	return nil
}
