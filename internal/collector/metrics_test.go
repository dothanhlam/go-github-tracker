package collector

import (
	"testing"
	"time"

	gh "github.com/google/go-github/v58/github"
)

// TestCalculateCycleTime tests the cycle time calculation
func TestCalculateCycleTime(t *testing.T) {
	tests := []struct {
		name      string
		createdAt time.Time
		mergedAt  time.Time
		want      int
	}{
		{
			name:      "1 hour difference",
			createdAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			mergedAt:  time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
			want:      1,
		},
		{
			name:      "24 hours difference",
			createdAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			mergedAt:  time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
			want:      24,
		},
		{
			name:      "30 minutes difference (rounds down)",
			createdAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			mergedAt:  time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
			want:      0,
		},
		{
			name:      "90 minutes difference",
			createdAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			mergedAt:  time.Date(2024, 1, 1, 11, 30, 0, 0, time.UTC),
			want:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateCycleTime(tt.createdAt, tt.mergedAt)
			if got != tt.want {
				t.Errorf("calculateCycleTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCalculateReviewTurnaround tests the review turnaround calculation
func TestCalculateReviewTurnaround(t *testing.T) {
	tests := []struct {
		name          string
		createdAt     time.Time
		firstReviewAt time.Time
		want          int
	}{
		{
			name:          "2 hours to first review",
			createdAt:     time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			firstReviewAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			want:          2,
		},
		{
			name:          "same day review (8 hours)",
			createdAt:     time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			firstReviewAt: time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC),
			want:          8,
		},
		{
			name:          "next day review",
			createdAt:     time.Date(2024, 1, 1, 16, 0, 0, 0, time.UTC),
			firstReviewAt: time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
			want:          18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateReviewTurnaround(tt.createdAt, tt.firstReviewAt)
			if got != tt.want {
				t.Errorf("calculateReviewTurnaround() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetFirstReviewTime tests finding the earliest review timestamp
func TestGetFirstReviewTime(t *testing.T) {
	time1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	time2 := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	time3 := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		reviews []*gh.PullRequestReview
		want    time.Time
	}{
		{
			name:    "empty reviews",
			reviews: []*gh.PullRequestReview{},
			want:    time.Time{},
		},
		{
			name: "single review",
			reviews: []*gh.PullRequestReview{
				{SubmittedAt: &gh.Timestamp{Time: time1}},
			},
			want: time1,
		},
		{
			name: "multiple reviews - earliest is first",
			reviews: []*gh.PullRequestReview{
				{SubmittedAt: &gh.Timestamp{Time: time3}},
				{SubmittedAt: &gh.Timestamp{Time: time1}},
				{SubmittedAt: &gh.Timestamp{Time: time2}},
			},
			want: time3,
		},
		{
			name: "multiple reviews - earliest is last",
			reviews: []*gh.PullRequestReview{
				{SubmittedAt: &gh.Timestamp{Time: time1}},
				{SubmittedAt: &gh.Timestamp{Time: time2}},
				{SubmittedAt: &gh.Timestamp{Time: time3}},
			},
			want: time3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFirstReviewTime(tt.reviews)
			if !got.Equal(tt.want) {
				t.Errorf("getFirstReviewTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExtractReviewers tests extracting unique reviewer usernames
func TestExtractReviewers(t *testing.T) {
	tests := []struct {
		name    string
		reviews []*gh.PullRequestReview
		want    int // number of unique reviewers
	}{
		{
			name:    "no reviews",
			reviews: []*gh.PullRequestReview{},
			want:    0,
		},
		{
			name: "single reviewer",
			reviews: []*gh.PullRequestReview{
				{User: &gh.User{Login: gh.String("reviewer1")}},
			},
			want: 1,
		},
		{
			name: "multiple reviews from same reviewer",
			reviews: []*gh.PullRequestReview{
				{User: &gh.User{Login: gh.String("reviewer1")}},
				{User: &gh.User{Login: gh.String("reviewer1")}},
				{User: &gh.User{Login: gh.String("reviewer1")}},
			},
			want: 1,
		},
		{
			name: "multiple unique reviewers",
			reviews: []*gh.PullRequestReview{
				{User: &gh.User{Login: gh.String("reviewer1")}},
				{User: &gh.User{Login: gh.String("reviewer2")}},
				{User: &gh.User{Login: gh.String("reviewer3")}},
			},
			want: 3,
		},
		{
			name: "mixed - some duplicate reviewers",
			reviews: []*gh.PullRequestReview{
				{User: &gh.User{Login: gh.String("reviewer1")}},
				{User: &gh.User{Login: gh.String("reviewer2")}},
				{User: &gh.User{Login: gh.String("reviewer1")}},
				{User: &gh.User{Login: gh.String("reviewer3")}},
				{User: &gh.User{Login: gh.String("reviewer2")}},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractReviewers(tt.reviews)
			if len(got) != tt.want {
				t.Errorf("extractReviewers() returned %d reviewers, want %d", len(got), tt.want)
			}
		})
	}
}

// TestCountReviewsByState tests counting reviews by state
func TestCountReviewsByState(t *testing.T) {
	tests := []struct {
		name    string
		reviews []*gh.PullRequestReview
		state   string
		want    int
	}{
		{
			name:    "no reviews",
			reviews: []*gh.PullRequestReview{},
			state:   "APPROVED",
			want:    0,
		},
		{
			name: "all approved",
			reviews: []*gh.PullRequestReview{
				{State: gh.String("APPROVED")},
				{State: gh.String("APPROVED")},
			},
			state: "APPROVED",
			want:  2,
		},
		{
			name: "mixed states - count approved",
			reviews: []*gh.PullRequestReview{
				{State: gh.String("APPROVED")},
				{State: gh.String("CHANGES_REQUESTED")},
				{State: gh.String("APPROVED")},
				{State: gh.String("COMMENTED")},
			},
			state: "APPROVED",
			want:  2,
		},
		{
			name: "mixed states - count changes requested",
			reviews: []*gh.PullRequestReview{
				{State: gh.String("APPROVED")},
				{State: gh.String("CHANGES_REQUESTED")},
				{State: gh.String("APPROVED")},
				{State: gh.String("CHANGES_REQUESTED")},
			},
			state: "CHANGES_REQUESTED",
			want:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countReviewsByState(tt.reviews, tt.state)
			if got != tt.want {
				t.Errorf("countReviewsByState() = %v, want %v", got, tt.want)
			}
		})
	}
}
