# MCP Server Configuration Guide

## Overview

The MCP server now supports flexible database configuration through environment variables. You can easily switch between SQLite (local development) and PostgreSQL (production) by specifying different `.env` files.

## Configuration Options

### Option 1: Default (.env.local)

**mcp-config.json**:
```json
{
  "mcpServers": {
    "dora-metrics": {
      "command": "go",
      "args": ["run", "cmd/mcp-server/main.go"]
    }
  }
}
```

Automatically uses `.env.local` (SQLite)

### Option 2: Specify Environment File

**mcp-config.json**:
```json
{
  "mcpServers": {
    "dora-metrics": {
      "command": "go",
      "args": ["run", "cmd/mcp-server/main.go"],
      "env": {
        "ENV_FILE": ".env.production"
      }
    }
  }
}
```

Uses `.env.production` (PostgreSQL or any other config)

### Option 3: Multiple Environments

**mcp-config.json**:
```json
{
  "mcpServers": {
    "dora-local": {
      "command": "go",
      "args": ["run", "cmd/mcp-server/main.go"],
      "env": { "ENV_FILE": ".env.local" }
    },
    "dora-prod": {
      "command": "go",
      "args": ["run", "cmd/mcp-server/main.go"],
      "env": { "ENV_FILE": ".env.production" }
    }
  }
}
```

Query different environments from AI

## Environment Files

### .env.local (SQLite - Development)
```bash
DB_DRIVER=sqlite3
DB_URL=./data/dora_metrics.db
```

### .env.production (PostgreSQL - Production)
```bash
DB_DRIVER=postgres
DB_URL=postgresql://user:password@host:5432/dora_metrics
```

### .env.staging (PostgreSQL - Staging)
```bash
DB_DRIVER=postgres
DB_URL=postgresql://user:password@staging-host:5432/dora_metrics
```

## Usage with AI

Once configured, ask your AI assistant:

**Local database**:
> "Show me team velocity from the local database"

**Production database** (if configured):
> "Show me team velocity from the production database"

The MCP server automatically connects to the correct database based on your configuration!

## Summary

✅ **Flexible** - Switch databases by changing ENV_FILE  
✅ **Multiple environments** - Configure local, staging, prod  
✅ **No code changes** - Just update mcp-config.json  
✅ **Secure** - Database credentials stay in .env files
