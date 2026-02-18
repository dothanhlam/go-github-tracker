# go-github-tracker

A DORA metrics tracking tool that collects Pull Request data from GitHub to measure team productivity and delivery performance.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![AWS Lambda](https://img.shields.io/badge/AWS-Lambda-FF9900?style=flat&logo=amazon-aws)](https://aws.amazon.com/lambda/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)

## Overview

**go-github-tracker** helps Engineering Managers gain visibility into team performance by automatically collecting and analyzing Pull Request metrics from GitHub. The tool focuses on DORA (DevOps Research and Assessment) metrics, specifically:

- **Team Velocity**: PRs merged per week
- **Lead Time**: Time from PR creation to merge (median & P95)
- **Review Metrics**: Turnaround time, comment ratios, approval rates, knowledge sharing

**📖 [View Complete Features Documentation](docs/FEATURES.md)**

## Features

✅ **Data Collection**
- Automated GitHub PR metrics collection
- Team-based metrics tracking
- Review turnaround time analysis
- Knowledge sharing metrics
- DORA metrics (Lead Time, Deployment Frequency)

✅ **Storage & Analysis**
- SQLite/PostgreSQL support
- Pre-built SQL views for common queries
- Migration tracking system
- Efficient data aggregation

✅ **REST API** ✨ NEW
- HTTP endpoints for all metrics
- API key authentication
- Query parameters for date ranges
- Health check endpoint
- See [API_SERVER.md](API_SERVER.md) for details

✅ **AI-Powered Insights**
- MCP server integration for AI assistants
- Natural language queries via SQLite MCP
- Direct SQL access to all metrics
- See [MCP-QUICKSTART.md](MCP-QUICKSTART.md) for setup

✅ **AWS Deployment Ready**
- Complete Terraform infrastructure
- Lambda function for scheduled collection
- RDS PostgreSQL support
- EventBridge scheduling
- CloudWatch monitoring
- See [terraform/README.md](terraform/README.md) for deployment

## Current Status

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1: Foundation | ✅ Complete | Database schema, migrations, configuration |
| Phase 2: Collector | ✅ Complete | GitHub API integration, data collection |
| Phase 2.5: Testing | ✅ Complete | Unit tests, CI/CD with GitHub Actions |
| Phase 3: SQL Views | ✅ Complete | Pre-built views for metrics analysis |
| **Phase 3: API** | **✅ Complete** | **REST API with 7 endpoints** |
| Phase 3.5: MCP | ✅ Working | SQLite MCP server for AI queries |
| Phase 4: AWS Deploy | 🚧 Ready | Terraform complete, awaiting deployment |
| Phase 5: Dashboard | 📋 Planned | Vue.js frontend (future) |

---

## Architecture

```
┌─────────────────┐
│  EventBridge    │  Triggers every 4 hours
│  (Scheduler)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐      ┌──────────────┐
│  AWS Lambda     │─────▶│   GitHub     │
│  (Go Collector) │◀─────│     API      │
└────────┬────────┘      └──────────────┘
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │
│   (Metrics DB)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   API Gateway   │  Exposes metrics via REST
│   (Lambda)      │
└─────────────────┘
```

### Technology Stack

- **Backend**: Go 1.21+
- **Database**: PostgreSQL 13+
- **GitHub API**: `google/go-github`
- **Database ORM**: `jmoiron/sqlx`
- **Authentication**: OAuth2
- **Deployment**: AWS Lambda (ARM64)
- **Scheduler**: AWS EventBridge

---

## Getting Started

### Prerequisites

- **Go**: 1.21 or higher ([install](https://go.dev/doc/install))
- **PostgreSQL**: 13+ ([install](https://www.postgresql.org/download/))
- **GitHub Personal Access Token**: With `repo` scope ([create](https://github.com/settings/tokens))
- **AWS Account**: For Lambda deployment (optional for local development)

### Installation

```bash
# Clone the repository
git clone https://github.com/dothanhlam/go-github-tracker.git
cd go-github-tracker

# Download dependencies
go mod download

# Build the project
go build -o bin/go-github-tracker ./cmd/collector
```

### Database Setup

```bash
# Create the database
createdb dora_metrics

# Run migrations
psql -d dora_metrics -f migrations/001_create_teams.sql
psql -d dora_metrics -f migrations/002_create_team_memberships.sql
psql -d dora_metrics -f migrations/003_create_pr_metrics.sql
psql -d dora_metrics -f migrations/004_create_views.sql
```

### Configuration

Create a `.env` file in the project root:

```bash
# GitHub Configuration
GITHUB_PAT=ghp_your_personal_access_token_here

# Database Configuration
DB_URL=postgres://user:password@localhost:5432/dora_metrics?sslmode=disable

# Collection Configuration
# Limit PR collection to last N days (default: 90)
# Prevents performance issues with large repositories
COLLECTION_LOOKBACK_DAYS=90

# Team Configuration (JSON array)
TEAM_CONFIG_JSON='[
  {
    "team_id": 1,
    "members": [
      {"username": "alice", "allocation": 1.0},
      {"username": "bob", "allocation": 0.5}
    ]
  }
]'

# Repositories to track (comma-separated)
REPOSITORIES=owner/repo1,owner/repo2
```

### Running Locally

```bash
# Run the collector
./bin/go-github-tracker

# Or build and run the API server
go build -o bin/api-server cmd/api-server/main.go
API_KEYS=test-key ./bin/api-server
```

---

## Usage

### Querying Metrics

#### Team Velocity (SQL)
```sql
SELECT 
    team_id,
    week,
    prs_merged,
    avg_cycle_time_hours
FROM view_team_velocity
WHERE team_id = 1
ORDER BY week DESC
LIMIT 10;
```

#### DORA Lead Time (SQL)
```sql
SELECT 
    team_id,
    month,
    median_lead_time_hours,
    p95_lead_time_hours
FROM view_dora_lead_time
WHERE team_id = 1
ORDER BY month DESC
LIMIT 6;
```

#### REST API (Available Now)
```bash
# Health check
curl http://localhost:8080/api/v1/health

# List teams
curl -H "X-API-Key: your-key" \
  http://localhost:8080/api/v1/teams

# Get team velocity
curl -H "X-API-Key: your-key" \
  "http://localhost:8080/api/v1/teams/1/velocity?start_date=2026-01-01&end_date=2026-02-15"

# Get DORA lead time
curl -H "X-API-Key: your-key" \
  http://localhost:8080/api/v1/teams/1/lead-time

# Get review metrics
curl -H "X-API-Key: your-key" \
  http://localhost:8080/api/v1/teams/1/review-turnaround

curl -H "X-API-Key: your-key" \
  http://localhost:8080/api/v1/teams/1/review-engagement

curl -H "X-API-Key: your-key" \
  http://localhost:8080/api/v1/teams/1/knowledge-sharing
```

**See [API_SERVER.md](API_SERVER.md) for complete API documentation.**

---

## Development

### Project Structure

```
.
├── .agent/                  # AI-specific documentation
│   ├── context/            # Architecture, conventions, dependencies
│   ├── planning/           # Plan, requirements, roadmap
│   └── workflows/          # Build, test, deploy workflows
├── cmd/
│   ├── collector/          # Main collector application
│   ├── api-server/         # REST API server ✨ NEW
│   └── mcp-server/         # MCP server (optional)
├── internal/
│   ├── api/                # API handlers, middleware ✨ NEW
│   ├── service/            # Business logic layer ✨ NEW
│   ├── config/             # Configuration management
│   ├── github/             # GitHub API client
│   ├── database/           # Database operations
│   ├── collector/          # PR collection logic
│   ├── metrics/            # Metrics calculation
│   └── mcp/                # MCP resources and tools
├── pkg/                    # Public libraries (if any)
├── migrations/             # Database migration scripts
├── terraform/              # AWS infrastructure ✨ NEW
├── tests/                  # Integration tests
├── API_SERVER.md          # API documentation ✨ NEW
├── MCP-QUICKSTART.md      # MCP setup guide
└── README.md              # This file
```

### Development Workflows

- **Build**: See [.agent/workflows/build.md](.agent/workflows/build.md)
- **Test**: See [.agent/workflows/test.md](.agent/workflows/test.md)
- **Deploy**: See [.agent/workflows/deploy.md](.agent/workflows/deploy.md) *(coming soon)*

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

---

## Deployment

### AWS Lambda Deployment

```bash
# Build for ARM64
GOOS=linux GOARCH=arm64 go build -o bootstrap ./cmd/collector

# Package for Lambda
zip function.zip bootstrap

# Deploy using AWS CLI
aws lambda update-function-code \
  --function-name dora-metrics-collector \
  --zip-file fileb://function.zip
```

### Environment Variables (Lambda)

Configure these in the Lambda function settings:

- `GITHUB_PAT` - GitHub Personal Access Token (use Secrets Manager in production)
- `DB_URL` - PostgreSQL connection string
- `TEAM_CONFIG_JSON` - Team configuration JSON
- `REPOSITORIES` - Comma-separated list of repositories
- `COLLECTION_LOOKBACK_DAYS` - Number of days to look back (default: 90, prevents performance issues)

---

## Roadmap

See [.agent/planning/plan.md](.agent/planning/plan.md) for detailed milestones.

- ✅ **Phase 1**: Foundation & Setup - Database, migrations, configuration
- ✅ **Phase 2**: Data Collection - GitHub API integration, PR metrics
- ✅ **Phase 2.5**: Testing & QA - Unit tests, CI/CD with GitHub Actions
- ✅ **Phase 3**: SQL Views & API - Pre-built views, REST API with 7 endpoints
- ✅ **Phase 3.5**: MCP Integration - SQLite MCP server for AI queries
- 🚧 **Phase 4**: AWS Deployment - Terraform complete, ready to deploy
- 📅 **Phase 5**: Dashboard - Vue.js frontend (future)

---

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Read the documentation** in `.agent/planning/` before starting
2. **Follow Go conventions** outlined in `.agent/context/conventions.md`
3. **Write tests** for all new functionality (90%+ coverage)
4. **Update documentation** as needed
5. **Submit a pull request** with a clear description

---

## License

MIT License - See [LICENSE](LICENSE) for details

---

## Support

For questions or issues:
- **Documentation**: See `.agent/` directory for detailed planning and architecture
- **Issues**: [GitHub Issues](https://github.com/dothanhlam/go-github-tracker/issues)
- **Discussions**: [GitHub Discussions](https://github.com/dothanhlam/go-github-tracker/discussions)
