package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Resource handlers

// getTeamVelocity returns team velocity metrics
func (h *MetricsHandler) getTeamVelocity(ctx context.Context, uri string) (string, error) {
	teamID, err := extractTeamID(uri)
	if err != nil {
		return "", err
	}

	query := `
		SELECT 
			month,
			prs_merged,
			avg_cycle_time_hours,
			median_cycle_time_hours,
			p95_cycle_time_hours
		FROM view_team_velocity
		WHERE team_id = ?
		ORDER BY month DESC
		LIMIT 12
	`

	var results []map[string]interface{}
	rows, err := h.db.Query(query, teamID)
	if err != nil {
		return "", fmt.Errorf("failed to query velocity: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var month string
		var prsMerged, avgCycle, medianCycle, p95Cycle interface{}
		
		if err := rows.Scan(&month, &prsMerged, &avgCycle, &medianCycle, &p95Cycle); err != nil {
			return "", err
		}

		results = append(results, map[string]interface{}{
			"month":                   month,
			"prs_merged":              prsMerged,
			"avg_cycle_time_hours":    avgCycle,
			"median_cycle_time_hours": medianCycle,
			"p95_cycle_time_hours":    p95Cycle,
		})
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// getReviewTurnaround returns review turnaround metrics
func (h *MetricsHandler) getReviewTurnaround(ctx context.Context, uri string) (string, error) {
	teamID, err := extractTeamID(uri)
	if err != nil {
		return "", err
	}

	query := `
		SELECT 
			month,
			avg_turnaround_hours,
			min_turnaround_hours,
			max_turnaround_hours,
			within_24h_count,
			over_24h_count
		FROM view_review_turnaround
		WHERE team_id = ?
		ORDER BY month DESC
		LIMIT 12
	`

	var results []map[string]interface{}
	rows, err := h.db.Query(query, teamID)
	if err != nil {
		return "", fmt.Errorf("failed to query review turnaround: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var month string
		var avgTurnaround, minTurnaround, maxTurnaround, within24h, over24h interface{}
		
		if err := rows.Scan(&month, &avgTurnaround, &minTurnaround, &maxTurnaround, &within24h, &over24h); err != nil {
			return "", err
		}

		results = append(results, map[string]interface{}{
			"month":                  month,
			"avg_turnaround_hours":   avgTurnaround,
			"min_turnaround_hours":   minTurnaround,
			"max_turnaround_hours":   maxTurnaround,
			"within_24h_count":       within24h,
			"over_24h_count":         over24h,
		})
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// getReviewEngagement returns review engagement metrics
func (h *MetricsHandler) getReviewEngagement(ctx context.Context, uri string) (string, error) {
	teamID, err := extractTeamID(uri)
	if err != nil {
		return "", err
	}

	query := `
		SELECT 
			month,
			avg_comments_per_pr,
			avg_conversations_per_pr,
			avg_reviewers_per_pr,
			changes_requested_rate,
			approval_rate
		FROM view_review_engagement
		WHERE team_id = ?
		ORDER BY month DESC
		LIMIT 12
	`

	var results []map[string]interface{}
	rows, err := h.db.Query(query, teamID)
	if err != nil {
		return "", fmt.Errorf("failed to query review engagement: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var month string
		var avgComments, avgConversations, avgReviewers, changesRate, approvalRate interface{}
		
		if err := rows.Scan(&month, &avgComments, &avgConversations, &avgReviewers, &changesRate, &approvalRate); err != nil {
			return "", err
		}

		results = append(results, map[string]interface{}{
			"month":                      month,
			"avg_comments_per_pr":        avgComments,
			"avg_conversations_per_pr":   avgConversations,
			"avg_reviewers_per_pr":       avgReviewers,
			"changes_requested_rate":     changesRate,
			"approval_rate":              approvalRate,
		})
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// getKnowledgeSharing returns knowledge sharing metrics
func (h *MetricsHandler) getKnowledgeSharing(ctx context.Context, uri string) (string, error) {
	teamID, err := extractTeamID(uri)
	if err != nil {
		return "", err
	}

	query := `
		SELECT 
			month,
			avg_reviewers_per_pr,
			avg_external_reviewers,
			external_reviewer_rate,
			total_external_reviews,
			prs_with_external_reviews
		FROM view_knowledge_sharing
		WHERE team_id = ?
		ORDER BY month DESC
		LIMIT 12
	`

	var results []map[string]interface{}
	rows, err := h.db.Query(query, teamID)
	if err != nil {
		return "", fmt.Errorf("failed to query knowledge sharing: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var month string
		var avgReviewers, avgExternal, externalRate, totalExternal, prsWithExternal interface{}
		
		if err := rows.Scan(&month, &avgReviewers, &avgExternal, &externalRate, &totalExternal, &prsWithExternal); err != nil {
			return "", err
		}

		results = append(results, map[string]interface{}{
			"month":                       month,
			"avg_reviewers_per_pr":        avgReviewers,
			"avg_external_reviewers":      avgExternal,
			"external_reviewer_rate":      externalRate,
			"total_external_reviews":      totalExternal,
			"prs_with_external_reviews":   prsWithExternal,
		})
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// getRecentPRs returns recent PRs with full metrics
func (h *MetricsHandler) getRecentPRs(ctx context.Context, uri string) (string, error) {
	teamID, err := extractTeamID(uri)
	if err != nil {
		return "", err
	}

	query := `
		SELECT 
			pr_number,
			repository,
			title,
			author,
			state,
			created_at,
			merged_at,
			cycle_time_hours,
			review_turnaround_hours,
			reviewers_count,
			review_comments_count
		FROM pr_metrics
		WHERE team_id = ?
		ORDER BY created_at DESC
		LIMIT 50
	`

	var results []map[string]interface{}
	rows, err := h.db.Query(query, teamID)
	if err != nil {
		return "", fmt.Errorf("failed to query recent PRs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var prNumber int
		var repository, title, author, state string
		var createdAt, mergedAt, cycleTime, turnaround, reviewers, comments interface{}
		
		if err := rows.Scan(&prNumber, &repository, &title, &author, &state, &createdAt, &mergedAt, &cycleTime, &turnaround, &reviewers, &comments); err != nil {
			return "", err
		}

		results = append(results, map[string]interface{}{
			"pr_number":                prNumber,
			"repository":               repository,
			"title":                    title,
			"author":                   author,
			"state":                    state,
			"created_at":               createdAt,
			"merged_at":                mergedAt,
			"cycle_time_hours":         cycleTime,
			"review_turnaround_hours":  turnaround,
			"reviewers_count":          reviewers,
			"review_comments_count":    comments,
		})
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Helper function to extract team ID from URI
func extractTeamID(uri string) (int, error) {
	parts := strings.Split(uri, "/")
	for i, part := range parts {
		if part == "team" && i+1 < len(parts) {
			return strconv.Atoi(parts[i+1])
		}
	}
	return 0, fmt.Errorf("team_id not found in URI: %s", uri)
}
