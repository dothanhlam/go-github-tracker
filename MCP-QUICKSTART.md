# Quick Start: MCP Server for DORA Metrics

## Option 1: Use SQLite MCP Server (Recommended - Works Immediately)

The easiest way to start querying your metrics with AI is using the standard SQLite MCP server.

### Setup

1. **Add to your AI assistant's MCP configuration**:

```json
{
  "mcpServers": {
    "dora-metrics": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-sqlite",
        "./data/dora_metrics.db"
      ]
    }
  }
}
```

2. **Restart your AI assistant** (Gemini Desktop, Claude Desktop, etc.)

3. **Start querying!**

### Example Queries

Ask your AI assistant:

> "Query the view_team_velocity table for team_id 1"

> "Show me all PRs with review_turnaround_hours greater than 48"

> "What's the average cycle time for team 1 over the last 3 months?"

> "Find PRs that took the longest to review"

### Available Views

Your AI can query these views directly:
- `view_team_velocity` - Team velocity metrics
- `view_review_turnaround` - Review turnaround times
- `view_review_engagement` - Review engagement metrics
- `view_knowledge_sharing` - Knowledge sharing metrics
- `pr_metrics` - Raw PR data

## Option 2: Custom Go MCP Server (Needs API Fixes)

The custom Go MCP server with specialized tools needs API compatibility fixes before it can run.

**Status**: Code structure complete, but needs updates to match mcp-go SDK v0.43.2 API

**Location**: `cmd/mcp-server/main.go`

**To fix**: Update resource and tool handlers to match the SDK's expected signatures (see implementation guide in walkthrough.md)

## Recommendation

**Use Option 1 (SQLite MCP Server)** for immediate AI-powered insights. It works out of the box and gives your AI full SQL query access to your metrics database.

The custom Go server (Option 2) can be completed later if you want specialized tools and more structured queries.
