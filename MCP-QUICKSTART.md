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

## Option 2: Custom Go MCP Server (Optional - Not Required)

The custom Go MCP server with specialized tools is **optional** and not required since the SQLite MCP server (Option 1) is more flexible.

**Status**: 
- ✅ REST API is now complete (`API_SERVER.md`)
- ✅ Code structure complete for custom MCP server
- ⚠️ Needs API compatibility fixes to match mcp-go SDK v0.43.2

**Location**: `cmd/mcp-server/main.go`, `internal/mcp/`

**Note**: The SQLite MCP server (Option 1) is recommended because:
- Works immediately without fixes
- AI can write any SQL query (more flexible than predefined tools)
- Direct access to all views and tables
- No maintenance needed

## Recommendation

**Use Option 1 (SQLite MCP Server)** for immediate AI-powered insights. It works out of the box and gives your AI full SQL query access to your metrics database.

The custom Go server (Option 2) is optional and can be completed later if you want specialized tools, but it's not necessary since the SQLite version is more powerful and flexible.

