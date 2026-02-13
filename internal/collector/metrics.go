package collector

import (
	"time"

	gh "github.com/google/go-github/v58/github"
)

// calculateCycleTime calculates hours from creation to merge
func calculateCycleTime(createdAt, mergedAt time.Time) int {
	duration := mergedAt.Sub(createdAt)
	return int(duration.Hours())
}

// calculateReviewTurnaround calculates hours from creation to first review
func calculateReviewTurnaround(createdAt, firstReviewAt time.Time) int {
	duration := firstReviewAt.Sub(createdAt)
	return int(duration.Hours())
}

// getFirstReviewTime finds the earliest review timestamp
func getFirstReviewTime(reviews []*gh.PullRequestReview) time.Time {
	if len(reviews) == 0 {
		return time.Time{}
	}

	firstReview := reviews[0].GetSubmittedAt().Time
	for _, review := range reviews[1:] {
		if review.GetSubmittedAt().Time.Before(firstReview) {
			firstReview = review.GetSubmittedAt().Time
		}
	}

	return firstReview
}

// extractReviewers extracts unique reviewer usernames
func extractReviewers(reviews []*gh.PullRequestReview) []string {
	reviewerMap := make(map[string]bool)
	for _, review := range reviews {
		username := review.GetUser().GetLogin()
		reviewerMap[username] = true
	}

	var reviewers []string
	for username := range reviewerMap {
		reviewers = append(reviewers, username)
	}

	return reviewers
}

// countReviewsByState counts reviews with a specific state
func countReviewsByState(reviews []*gh.PullRequestReview, state string) int {
	count := 0
	for _, review := range reviews {
		if review.GetState() == state {
			count++
		}
	}
	return count
}

// countConversations counts unique conversation threads
func countConversations(comments []*gh.PullRequestComment) int {
	// Group by in_reply_to_id to count threads
	threadMap := make(map[int64]bool)
	for _, comment := range comments {
		// If it's a top-level comment, use its own ID
		if comment.GetInReplyTo() == 0 {
			threadMap[comment.GetID()] = true
		} else {
			// Otherwise use the parent comment ID
			threadMap[comment.GetInReplyTo()] = true
		}
	}

	return len(threadMap)
}

// countExternalReviewers counts reviewers outside the team
func (c *Collector) countExternalReviewers(reviewers []string, teamID int) int {
	count := 0
	for _, reviewer := range reviewers {
		if c.teamMgr.IsExternalReviewer(reviewer, teamID) {
			count++
		}
	}
	return count
}
