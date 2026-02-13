package mcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mcpsdk "github.com/mark3labs/mcp-go/server"
)

// Tool handlers

// queryMetrics queries metrics with filters
func (h *MetricsHandler) queryMetrics(ctx context.Context, arguments map[string]interface{}) (*mcpsdk.ToolResponse, error) {
	teamID := int(arguments["team_id"].(float64))
	metricType := arguments["metric_type"].(string)
	
	var query string
	var tableName string
	
	switch metricType {
	case "velocity":
		tableName = "view_team_velocity"
		query = `SELECT * FROM view_team_velocity WHERE team_id = ?`
	case "review-turnaround":
		tableName = "view_review_turnaround"
		query = `SELECT * FROM view_review_turnaround WHERE team_id = ?`
	case "engagement":
		tableName = "view_review_engagement"
		query = `SELECT * FROM view_review_engagement WHERE team_id = ?`
	case "knowledge-sharing":
		tableName = "view_knowledge_sharing"
		query = `SELECT * FROM view_knowledge_sharing WHERE team_id = ?`
	default:
		return mcpsdk.NewToolResponse(mcpsdk.NewTextContent(fmt.Sprintf("Invalid metric_type: %s", metricType))), nil
	}

	// Add date filters if provided
	args := []interface{}{teamID}
	if startDate, ok := arguments["start_date"].(string); ok && startDate != "" {
		query += " AND month >= ?"
		args = append(args, startDate)
	}
	if endDate, ok := arguments["end_date"].(string); ok && endDate != "" {
		query += " AND month <= ?"
		args = append(args, endDate)
	}
	
	query += " ORDER BY month DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", tableName, err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcpsdk.NewToolResponse(mcpsdk.NewTextContent(string(data))), nil
}

// getTeamSummary gets comprehensive team summary
func (h *MetricsHandler) getTeamSummary(ctx context.Context, arguments map[string]interface{}) (*mcpsdk.ToolResponse, error) {
	teamID := int(arguments["team_id"].(float64))
	period := "month"
	if p, ok := arguments["period"].(string); ok {
		period = p
	}

	// Calculate date range based on period
	var monthsBack int
	switch period {
	case "week":
		monthsBack = 0 // Current month
	case "month":
		monthsBack = 1
	case "quarter":
		monthsBack = 3
	default:
		monthsBack = 1
	}

	startDate := time.Now().AddDate(0, -monthsBack, 0).Format("2006-01")

	summary := make(map[string]interface{})

	// Get velocity metrics
	velocityQuery := `SELECT * FROM view_team_velocity WHERE team_id = ? AND month >= ? ORDER BY month DESC`
	rows, err := h.db.Query(velocityQuery, teamID, startDate)
	if err == nil {
		defer rows.Close()
		summary["velocity"] = scanRows(rows)
	}

	// Get review turnaround
	turnaroundQuery := `SELECT * FROM view_review_turnaround WHERE team_id = ? AND month >= ? ORDER BY month DESC`
	rows, err = h.db.Query(turnaroundQuery, teamID, startDate)
	if err == nil {
		defer rows.Close()
		summary["review_turnaround"] = scanRows(rows)
	}

	// Get engagement
	engagementQuery := `SELECT * FROM view_review_engagement WHERE team_id = ? AND month >= ? ORDER BY month DESC`
	rows, err = h.db.Query(engagementQuery, teamID, startDate)
	if err == nil {
		defer rows.Close()
		summary["engagement"] = scanRows(rows)
	}

	// Get knowledge sharing
	knowledgeQuery := `SELECT * FROM view_knowledge_sharing WHERE team_id = ? AND month >= ? ORDER BY month DESC`
	rows, err = h.db.Query(knowledgeQuery, teamID, startDate)
	if err == nil {
		defer rows.Close()
		summary["knowledge_sharing"] = scanRows(rows)
	}

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcpsdk.NewToolResponse(mcpsdk.NewTextContent(string(data))), nil
}

// analyzeTrends analyzes metric trends over time
func (h *MetricsHandler) analyzeTrends(ctx context.Context, arguments map[string]interface{}) (*mcpsdk.ToolResponse, error) {
	teamID := int(arguments["team_id"].(float64))
	metricName := arguments["metric_name"].(string)
	months := 6
	if m, ok := arguments["months"].(float64); ok {
		months = int(m)
	}

	startDate := time.Now().AddDate(0, -months, 0).Format("2006-01")

	// Determine which view to query based on metric name
	var query string
	if contains(metricName, []string{"cycle_time", "prs_merged"}) {
		query = fmt.Sprintf(`SELECT month, %s FROM view_team_velocity WHERE team_id = ? AND month >= ? ORDER BY month ASC`, metricName)
	} else if contains(metricName, []string{"turnaround"}) {
		query = fmt.Sprintf(`SELECT month, %s FROM view_review_turnaround WHERE team_id = ? AND month >= ? ORDER BY month ASC`, metricName)
	} else if contains(metricName, []string{"comments", "conversations", "approval", "changes_requested"}) {
		query = fmt.Sprintf(`SELECT month, %s FROM view_review_engagement WHERE team_id = ? AND month >= ? ORDER BY month ASC`, metricName)
	} else {
		query = fmt.Sprintf(`SELECT month, %s FROM view_knowledge_sharing WHERE team_id = ? AND month >= ? ORDER BY month ASC`, metricName)
	}

	rows, err := h.db.Query(query, teamID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query trends: %w", err)
	}
	defer rows.Close()

	var dataPoints []map[string]interface{}
	var values []float64
	
	for rows.Next() {
		var month string
		var value interface{}
		
		if err := rows.Scan(&month, &value); err != nil {
			return nil, err
		}

		dataPoints = append(dataPoints, map[string]interface{}{
			"month": month,
			"value": value,
		})

		if v, ok := value.(float64); ok {
			values = append(values, v)
		} else if v, ok := value.(int64); ok {
			values = append(values, float64(v))
		}
	}

	// Calculate trend
	var trend string
	var change float64
	if len(values) >= 2 {
		first := values[0]
		last := values[len(values)-1]
		change = ((last - first) / first) * 100
		
		if change > 5 {
			trend = "improving"
		} else if change < -5 {
			trend = "declining"
		} else {
			trend = "stable"
		}
	}

	result := map[string]interface{}{
		"metric":      metricName,
		"data_points": dataPoints,
		"trend":       trend,
		"change_pct":  change,
		"analysis":    fmt.Sprintf("The %s has %s by %.1f%% over the last %d months", metricName, trend, change, months),
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcpsdk.NewToolResponse(mcpsdk.NewTextContent(string(data))), nil
}

// findBottlenecks finds PRs with slow review times
func (h *MetricsHandler) findBottlenecks(ctx context.Context, arguments map[string]interface{}) (*mcpsdk.ToolResponse, error) {
	teamID := int(arguments["team_id"].(float64))
	thresholdHours := 48
	if t, ok := arguments["threshold_hours"].(float64); ok {
		thresholdHours = int(t)
	}
	limit := 10
	if l, ok := arguments["limit"].(float64); ok {
		limit = int(l)
	}

	query := `
		SELECT 
			pr_number,
			repository,
			title,
			author,
			review_turnaround_hours,
			cycle_time_hours,
			created_at
		FROM pr_metrics
		WHERE team_id = ? 
		  AND review_turnaround_hours > ?
		ORDER BY review_turnaround_hours DESC
		LIMIT ?
	`

	rows, err := h.db.Query(query, teamID, thresholdHours, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query bottlenecks: %w", err)
	}
	defer rows.Close()

	var bottlenecks []map[string]interface{}
	for rows.Next() {
		var prNumber int
		var repository, title, author string
		var turnaround, cycleTime interface{}
		var createdAt string
		
		if err := rows.Scan(&prNumber, &repository, &title, &author, &turnaround, &cycleTime, &createdAt); err != nil {
			return nil, err
		}

		bottlenecks = append(bottlenecks, map[string]interface{}{
			"pr_number":                prNumber,
			"repository":               repository,
			"title":                    title,
			"author":                   author,
			"review_turnaround_hours":  turnaround,
			"cycle_time_hours":         cycleTime,
			"created_at":               createdAt,
		})
	}

	result := map[string]interface{}{
		"threshold_hours": thresholdHours,
		"count":           len(bottlenecks),
		"bottlenecks":     bottlenecks,
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcpsdk.NewToolResponse(mcpsdk.NewTextContent(string(data))), nil
}

// Helper functions

func scanRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	var results []map[string]interface{}
	
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}
	
	return results
}

func contains(str string, list []string) bool {
	for _, item := range list {
		if strings.Contains(str, item) {
			return true
		}
	}
	return false
}
