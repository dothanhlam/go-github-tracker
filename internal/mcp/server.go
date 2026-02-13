package mcp

import (
	"github.com/dothanhlam/go-github-tracker/internal/database"
	mcpsdk "github.com/mark3labs/mcp-go/server"
)

// MetricsHandler handles MCP requests for metrics data
type MetricsHandler struct {
	db *database.DB
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(db *database.DB) *MetricsHandler {
	return &MetricsHandler{
		db: db,
	}
}

// RegisterResources registers all MCP resources
func (h *MetricsHandler) RegisterResources(server *mcpsdk.MCPServer) {
	// Resource 1: Team Velocity
	server.AddResource(
		"team/{team_id}/velocity",
		"Monthly team velocity metrics (cycle time, throughput)",
		"application/json",
		h.getTeamVelocity,
	)

	// Resource 2: Review Turnaround
	server.AddResource(
		"team/{team_id}/review-turnaround",
		"Review turnaround time metrics",
		"application/json",
		h.getReviewTurnaround,
	)

	// Resource 3: Review Engagement
	server.AddResource(
		"team/{team_id}/review-engagement",
		"Review engagement metrics (comments, approvals)",
		"application/json",
		h.getReviewEngagement,
	)

	// Resource 4: Knowledge Sharing
	server.AddResource(
		"team/{team_id}/knowledge-sharing",
		"Knowledge sharing metrics (external reviewers)",
		"application/json",
		h.getKnowledgeSharing,
	)

	// Resource 5: Recent PRs
	server.AddResource(
		"team/{team_id}/recent-prs",
		"Recent pull requests with full metrics",
		"application/json",
		h.getRecentPRs,
	)
}

// RegisterTools registers all MCP tools
func (h *MetricsHandler) RegisterTools(server *mcpsdk.MCPServer) {
	// Tool 1: Query Metrics
	server.AddTool(
		mcpsdk.Tool{
			Name:        "query_metrics",
			Description: "Query metrics with filters (date range, metric type)",
			InputSchema: mcpsdk.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"team_id": map[string]interface{}{
						"type":        "integer",
						"description": "Team ID to query",
					},
					"metric_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of metric: velocity, review-turnaround, engagement, knowledge-sharing",
						"enum":        []string{"velocity", "review-turnaround", "engagement", "knowledge-sharing"},
					},
					"start_date": map[string]interface{}{
						"type":        "string",
						"description": "Start date (YYYY-MM-DD)",
					},
					"end_date": map[string]interface{}{
						"type":        "string",
						"description": "End date (YYYY-MM-DD)",
					},
				},
				Required: []string{"team_id", "metric_type"},
			},
		},
		h.queryMetrics,
	)

	// Tool 2: Get Team Summary
	server.AddTool(
		mcpsdk.Tool{
			Name:        "get_team_summary",
			Description: "Get comprehensive team summary for a period",
			InputSchema: mcpsdk.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"team_id": map[string]interface{}{
						"type":        "integer",
						"description": "Team ID",
					},
					"period": map[string]interface{}{
						"type":        "string",
						"description": "Time period: week, month, quarter",
						"enum":        []string{"week", "month", "quarter"},
					},
				},
				Required: []string{"team_id"},
			},
		},
		h.getTeamSummary,
	)

	// Tool 3: Analyze Trends
	server.AddTool(
		mcpsdk.Tool{
			Name:        "analyze_trends",
			Description: "Analyze metric trends over time",
			InputSchema: mcpsdk.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"team_id": map[string]interface{}{
						"type":        "integer",
						"description": "Team ID",
					},
					"metric_name": map[string]interface{}{
						"type":        "string",
						"description": "Metric to analyze (e.g., avg_cycle_time_hours, avg_turnaround_hours)",
					},
					"months": map[string]interface{}{
						"type":        "integer",
						"description": "Number of months to analyze (default: 6)",
					},
				},
				Required: []string{"team_id", "metric_name"},
			},
		},
		h.analyzeTrends,
	)

	// Tool 4: Find Bottlenecks
	server.AddTool(
		mcpsdk.Tool{
			Name:        "find_bottlenecks",
			Description: "Find PRs with slow review times or long cycle times",
			InputSchema: mcpsdk.ToolInputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"team_id": map[string]interface{}{
						"type":        "integer",
						"description": "Team ID",
					},
					"threshold_hours": map[string]interface{}{
						"type":        "integer",
						"description": "Review time threshold in hours (default: 48)",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results (default: 10)",
					},
				},
				Required: []string{"team_id"},
			},
		},
		h.findBottlenecks,
	)
}
