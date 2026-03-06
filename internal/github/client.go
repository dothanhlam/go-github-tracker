package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API client
type Client struct {
	client *github.Client
	ctx    context.Context
}

// NewClient creates a new GitHub API client with authentication
func NewClient(pat string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
		ctx:    ctx,
	}
}

// FetchPRs fetches pull requests from a repository since a given date
func (c *Client) FetchPRs(owner, repo string, since time.Time) ([]*github.PullRequest, error) {
	fmt.Printf("  📥 Fetching PRs from %s/%s (since %s)...\n", owner, repo, since.Format("2006-01-02"))

	var allPRs []*github.PullRequest
	opts := &github.PullRequestListOptions{
		State:     "all", // Get open, closed, and merged PRs
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		prs, resp, err := c.client.PullRequests.List(c.ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch PRs: %w", err)
		}

		// Filter PRs by date - stop early if PRs are too old
		for _, pr := range prs {
			// Check if PR was updated before the lookback date
			if pr.UpdatedAt != nil && pr.UpdatedAt.Before(since) {
				fmt.Printf("  ⏹️  Stopped at PR #%d (updated %s, before lookback date)\n", 
					pr.GetNumber(), pr.UpdatedAt.Format("2006-01-02"))
				fmt.Printf("  ✓ Fetched %d PRs within lookback window\n", len(allPRs))
				return allPRs, nil // Early exit
			}
			allPRs = append(allPRs, pr)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage

		// Check rate limit
		if err := c.checkRateLimit(); err != nil {
			return nil, err
		}
	}

	fmt.Printf("  ✓ Fetched %d PRs\n", len(allPRs))
	return allPRs, nil
}

// FetchReviews fetches all reviews for a pull request
func (c *Client) FetchReviews(owner, repo string, prNumber int) ([]*github.PullRequestReview, error) {
	opts := &github.ListOptions{
		PerPage: 100,
	}

	var allReviews []*github.PullRequestReview
	for {
		reviews, resp, err := c.client.PullRequests.ListReviews(c.ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch reviews: %w", err)
		}

		allReviews = append(allReviews, reviews...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReviews, nil
}

// FetchComments fetches all review comments for a pull request
func (c *Client) FetchComments(owner, repo string, prNumber int) ([]*github.PullRequestComment, error) {
	opts := &github.PullRequestListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allComments []*github.PullRequestComment
	for {
		comments, resp, err := c.client.PullRequests.ListComments(c.ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch comments: %w", err)
		}

		allComments = append(allComments, comments...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allComments, nil
}

// checkRateLimit checks GitHub API rate limit and waits if necessary
func (c *Client) checkRateLimit() error {
	rate, _, err := c.client.RateLimits(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	if rate.Core.Remaining < 100 {
		waitTime := time.Until(rate.Core.Reset.Time)
		fmt.Printf("  ⏳ Rate limit low (%d remaining), waiting %v...\n", rate.Core.Remaining, waitTime)
		time.Sleep(waitTime)
	}

	return nil
}

// FetchCommits fetches commits from a repository since a given date
func (c *Client) FetchCommits(owner, repo string, since time.Time) ([]*github.RepositoryCommit, error) {
	fmt.Printf("  📥 Fetching Commits from %s/%s (since %s)...\n", owner, repo, since.Format("2006-01-02"))

	var allCommits []*github.RepositoryCommit
	opts := &github.CommitsListOptions{
		Since: since,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		commits, resp, err := c.client.Repositories.ListCommits(c.ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch commits: %w", err)
		}

		allCommits = append(allCommits, commits...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage

		if err := c.checkRateLimit(); err != nil {
			return nil, err
		}
	}

	fmt.Printf("  ✓ Fetched %d commits\n", len(allCommits))
	return allCommits, nil
}

// FetchIssueComments fetches issue comments from a repository since a given date
func (c *Client) FetchIssueComments(owner, repo string, since time.Time) ([]*github.IssueComment, error) {
	fmt.Printf("  📥 Fetching Issue Comments from %s/%s (since %s)...\n", owner, repo, since.Format("2006-01-02"))

	var allComments []*github.IssueComment
	opts := &github.IssueListCommentsOptions{
		Since: &since,
		Sort:  github.String("updated"),
		Direction: github.String("desc"),
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		// Pass 0 to fetch comments for the entire repository, rather than a specific issue
		comments, resp, err := c.client.Issues.ListComments(c.ctx, owner, repo, 0, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch issue comments: %w", err)
		}

		// Filter comments by date
		for _, comment := range comments {
			if comment.UpdatedAt != nil && comment.UpdatedAt.Before(since) {
				fmt.Printf("  ✓ Fetched %d issue comments within lookback window\n", len(allComments))
				return allComments, nil
			}
			allComments = append(allComments, comment)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage

		if err := c.checkRateLimit(); err != nil {
			return nil, err
		}
	}

	fmt.Printf("  ✓ Fetched %d issue comments\n", len(allComments))
	return allComments, nil
}

// FetchCommitComments fetches commit comments from a repository
func (c *Client) FetchCommitComments(owner, repo string, since time.Time) ([]*github.RepositoryComment, error) {
	fmt.Printf("  📥 Fetching Commit Comments from %s/%s (since %s)...\n", owner, repo, since.Format("2006-01-02"))

	var allComments []*github.RepositoryComment
	opts := &github.ListOptions{
		PerPage: 100,
	}

	for {
		comments, resp, err := c.client.Repositories.ListComments(c.ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch commit comments: %w", err)
		}

		for _, comment := range comments {
			if comment.UpdatedAt != nil && !comment.UpdatedAt.Before(since) {
				allComments = append(allComments, comment)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage

		if err := c.checkRateLimit(); err != nil {
			return nil, err
		}
	}

	fmt.Printf("  ✓ Fetched %d commit comments\n", len(allComments))
	return allComments, nil
}
