# go-github-tracker

A DORA metrics tracking tool that collects Pull Request data from GitHub to measure team productivity and delivery performance.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![AWS Lambda](https://img.shields.io/badge/AWS-Lambda-FF9900?style=flat&logo=amazon-aws)](https://aws.amazon.com/lambda/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)

## Overview

**go-github-tracker** helps Engineering Managers gain visibility into team performance by automatically collecting and analyzing Pull Request metrics from GitHub. The tool focuses on DORA (DevOps Research and Assessment) metrics, specifically:

- **Team Velocity**: PRs merged per week
- **Lead Time**: Time from PR creation to merge (median & P95)

### Key Features

✅ **Team-Filtered Tracking** - Only track PRs from configured team members, even in shared repositories  
✅ **Multi-Team Support** - Handle team members who belong to multiple teams with weighted allocations  
✅ **Automated Collection** - Runs every 4 hours via AWS EventBridge  
✅ **Idempotent Operations** - Safe to re-run without creating duplicate data  
✅ **SQL Views** - Easy-to-query views for velocity and lead time metrics  
✅ **REST API** - JSON endpoints for programmatic access to metrics  

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
./bin/go-github-tracker collect

# Run the API server (if implemented)
./bin/go-github-tracker serve --port 8080
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

#### API Endpoints (Future)
```bash
# Get team velocity
curl http://localhost:8080/metrics/team/1/velocity

# Get DORA lead time
curl http://localhost:8080/metrics/team/1/lead-time
```

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
│   └── api/                # API server (future)
├── internal/
│   ├── config/             # Configuration management
│   ├── github/             # GitHub API client
│   ├── database/           # Database operations
│   ├── collector/          # PR collection logic
│   └── metrics/            # Metrics calculation
├── pkg/                    # Public libraries (if any)
├── migrations/             # Database migration scripts
├── tests/                  # Integration tests
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

---

## Roadmap

See [.agent/planning/roadmap.md](.agent/planning/roadmap.md) for detailed milestones.

- ✅ **Phase 1**: Foundation & Setup (Week 1)
- ⏳ **Phase 2**: Data Collection (Week 2-3)
- 📅 **Phase 3**: Metrics & API (Week 3-4)
- 📅 **Phase 4**: Deployment & Automation (Week 4-5)
- 📅 **Phase 5**: Production Hardening (Week 6-8)
- 🔮 **Future**: Vue.js Dashboard, Advanced Metrics

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
