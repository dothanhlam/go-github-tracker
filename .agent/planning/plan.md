# DORA Metrics Tool: Development Plan

## 1. Vision & Context

### Goal
Build a custom tracking tool for team productivity using GitHub APIs to measure DORA (DevOps Research and Assessment) metrics.

### Primary Users
Engineering Managers who need visibility into team performance and delivery metrics.

### Key Constraints
- **Team Filtering**: Only track a specific set of team members (the "Team Filter") even if they share repositories with other squads
- **Multi-team Support**: Handle team members who belong to multiple teams
- **Weighted Allocations**: Support fractional team membership (e.g., 50% on Team A, 50% on Team B)

### Technology Stack
- **Backend**: Go
- **Database**: PostgreSQL
- **Frontend**: Vue.js *(future phase)*
- **Deployment**: AWS Lambda (ARM64)
- **Scheduler**: AWS EventBridge (every 4 hours)

---

## 2. Technical Specifications

### Core Dependencies
- `google/go-github` - GitHub API interaction
- `jmoiron/sqlx` - Database operations with struct mapping
- `lib/pq` - PostgreSQL driver
- `golang.org/x/oauth2` - GitHub authentication

### Database Schema

#### `teams` Table
```sql
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### `team_memberships` Table
```sql
CREATE TABLE team_memberships (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    github_username VARCHAR(255) NOT NULL,
    allocation_weight DECIMAL(3,2) DEFAULT 1.00, -- 0.00 to 1.00
    joined_at TIMESTAMP NOT NULL,
    left_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(team_id, github_username, joined_at)
);
```

#### `pr_metrics` Table
```sql
CREATE TABLE pr_metrics (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    pr_number INTEGER NOT NULL,
    repository VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    title TEXT,
    created_at TIMESTAMP NOT NULL,
    merged_at TIMESTAMP,
    closed_at TIMESTAMP,
    cycle_time_hours INTEGER, -- Merged - Created
    state VARCHAR(50), -- open, closed, merged
    created_date DATE GENERATED ALWAYS AS (DATE(created_at)) STORED,
    UNIQUE(team_id, repository, pr_number)
);
```

#### SQL Views for Reporting
```sql
-- Team velocity (PRs merged per week)
CREATE VIEW view_team_velocity AS
SELECT 
    team_id,
    DATE_TRUNC('week', merged_at) as week,
    COUNT(*) as prs_merged,
    AVG(cycle_time_hours) as avg_cycle_time_hours
FROM pr_metrics
WHERE merged_at IS NOT NULL
GROUP BY team_id, week;

-- DORA Lead Time
CREATE VIEW view_dora_lead_time AS
SELECT 
    team_id,
    DATE_TRUNC('month', merged_at) as month,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY cycle_time_hours) as median_lead_time_hours,
    PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY cycle_time_hours) as p95_lead_time_hours
FROM pr_metrics
WHERE merged_at IS NOT NULL
GROUP BY team_id, month;
```

### Deployment Architecture
- **Runtime**: AWS Lambda (ARM64 for cost efficiency)
- **Trigger**: EventBridge scheduled rule (every 4 hours)
- **Environment Variables**:
  - `GITHUB_PAT` - GitHub Personal Access Token
  - `DB_URL` - PostgreSQL connection string
  - `TEAM_CONFIG_JSON` - JSON array of team configurations
  - `REPOSITORIES` - Comma-separated list of repos to track

---

## 3. Implementation Roadmap

### Phase 1: Foundation (Database & Setup) ✅ COMPLETE
- [x] Initialize Go module
- [x] Install core dependencies
  - [x] `google/go-github`
  - [x] `jmoiron/sqlx`
  - [x] `lib/pq`
  - [x] `golang.org/x/oauth2`
  - [x] `mattn/go-sqlite3` (for local dev)
  - [x] `joho/godotenv`
- [x] Create database migration scripts
  - [x] `001_create_teams.sql` (SQLite & PostgreSQL)
  - [x] `002_create_team_memberships.sql` (SQLite & PostgreSQL)
  - [x] `003_create_pr_metrics.sql` (SQLite & PostgreSQL)
  - [x] `004_create_views.sql` (SQLite & PostgreSQL)
- [x] Implement configuration management
  - [x] Environment variable loader
  - [x] Team configuration parser (JSON)
  - [x] Database connection pool setup
  - [x] SQLite/PostgreSQL driver switching
  - [x] Auto-generated team IDs (removed manual team_id from config)

### Phase 2: The Collector (Go Agent Logic) ✅ COMPLETE
- [x] **Member Sync Module**
  - [x] Load active team members from database
  - [x] Build in-memory lookup map: `github_username -> []TeamMembership`
  - [x] Handle time-based membership (joined_at/left_at)
- [x] **GitHub Polling Module**
  - [x] Authenticate with GitHub using OAuth2
  - [x] Fetch closed/merged PRs from configured repositories
  - [x] Implement pagination for large result sets
  - [x] Respect GitHub rate limits (primary & secondary)
- [x] **Filter & Attribution Engine**
  - [x] Iterate through fetched PRs
  - [x] Check if PR author is in membership map
  - [x] Calculate cycle time (merged_at - created_at)
  - [x] Attribute PR to all teams the author belongs to (weighted)
  - [x] Batch upsert into `pr_metrics` (use `ON CONFLICT` for idempotency)
- [x] **Error Handling & Logging**
  - [x] Structured logging (JSON format for CloudWatch)
  - [x] Retry logic for transient failures

### Phase 2.5: Testing & Quality Assurance ✅ COMPLETE
- [x] **Unit Tests**
  - [x] Collector metrics tests (cycle time, review turnaround, reviewers)
  - [x] Config parsing and validation tests
  - [x] Database migration tests
  - [x] 15 test functions, 42 sub-tests, 100% pass rate
- [x] **CI/CD Pipeline**
  - [x] GitHub Actions workflow for automated testing
  - [x] Run tests on push to main and pull requests
  - [x] Coverage reporting
  - [x] Test exclusions for incomplete packages (MCP)

### Phase 3: Reporting & API ✅ COMPLETE
- [x] **SQL Views** (already defined in schema)
  - [x] `view_team_velocity` - Working
  - [x] `view_dora_lead_time` - Working
  - [x] `view_review_turnaround` - Working
  - [x] `view_review_engagement` - Working
  - [x] `view_knowledge_sharing` - Working
- [x] **API Handler** ✅ COMPLETE
  - [x] Create standalone HTTP server with chi router
  - [x] Endpoint: `GET /api/v1/teams` - List teams
  - [x] Endpoint: `GET /api/v1/teams/{id}/velocity` - Team velocity
  - [x] Endpoint: `GET /api/v1/teams/{id}/lead-time` - DORA lead time
  - [x] Endpoint: `GET /api/v1/teams/{id}/review-turnaround` - Review turnaround
  - [x] Endpoint: `GET /api/v1/teams/{id}/review-engagement` - Review engagement
  - [x] Endpoint: `GET /api/v1/teams/{id}/knowledge-sharing` - Knowledge sharing
  - [x] Endpoint: `GET /api/v1/health` - Health check
  - [x] API key authentication via X-API-Key header
  - [x] Request logging middleware
  - [x] CORS support
  - [x] Graceful shutdown
  - [x] Documentation (API_SERVER.md)

### Phase 3.5: MCP Integration (AI-Powered Insights) ✅ WORKING (SQLite)
- [x] **MCP Server Setup**
  - [x] SQLite MCP Server (Recommended - Working)
    - [x] Configuration in mcp-config.json
    - [x] Direct SQL query access for AI
    - [x] All views accessible
    - [x] Documentation (MCP-QUICKSTART.md)
  - [ ] Custom Go MCP Server (Optional - Not Required)
    - [x] Server structure created
    - [x] Resources defined (5 views)
    - [x] Tools defined (4 analytical tools)
    - [ ] API compatibility fixes (optional enhancement)
    Note: SQLite MCP server is more flexible and already working

**Current Status**: SQLite MCP server is fully functional and recommended. Custom Go MCP server is optional and not required since SQLite version provides more flexibility.
### Phase 4: AWS Lambda Deployment 🚧 READY TO DEPLOY
- [x] **Terraform Infrastructure** ✅ COMPLETE
  - [x] VPC module (subnets, NAT gateway, routing)
  - [x] RDS PostgreSQL module (db.t4g.micro for dev)
  - [x] Lambda module (ARM64, IAM roles, VPC integration)
  - [x] EventBridge module (scheduled triggers every 4 hours)
  - [x] Secrets Manager module (GitHub PAT, DB credentials)
  - [x] Monitoring module (CloudWatch alarms)
  - [x] Dev environment configuration
  - [x] Build scripts (build-lambda.sh)
  - [x] Documentation (terraform/README.md, DEPLOYMENT_COMMANDS.md)
  - [x] Verification script (verify-terraform.sh)
- [ ] **AWS Deployment** (Not Yet Executed)
  - [ ] Create S3 backend for Terraform state
  - [ ] Create DynamoDB table for state locking
  - [ ] Store GitHub PAT in Secrets Manager
  - [ ] Build Lambda deployment package
  - [ ] Run terraform init & plan
  - [ ] Deploy infrastructure (terraform apply)
  - [ ] Verify scheduled execution
  - [ ] Monitor for 1 week

**Current Status**: All Terraform code complete and ready. Waiting for AWS deployment execution.

### Phase 5: Frontend (Future)
- [ ] Vue.js dashboard for visualizing metrics
- [ ] Team selection and filtering
- [ ] Date range selection
- [ ] Charts for velocity and lead time trends

---

## 4. Key Guardrails & Best Practices

### Authentication
- ✅ Always use `oauth2` for GitHub API calls
- ✅ Store GitHub PAT in AWS Secrets Manager (not environment variables in production)
- ✅ Rotate tokens regularly

### Performance
- ✅ Use Go concurrency (goroutines) for fetching multiple repositories
- ✅ Strictly follow GitHub's secondary rate limits (avoid triggering abuse detection)
- ✅ Implement exponential backoff for rate limit errors
- ✅ Use connection pooling for database operations

### Data Integrity
- ✅ All database operations must use `ON CONFLICT` for idempotency
- ✅ Allow re-runs without duplicate data
- ✅ Use database transactions for multi-table operations
- ✅ Validate data before insertion (non-null checks, date ranges)

### Observability
- ✅ Structured logging (JSON format)
- ✅ Emit CloudWatch metrics for:
  - PRs processed per run
  - API errors
  - Database query duration
- ✅ Set up CloudWatch alarms for failures

---

## 5. Success Criteria

### MVP (Minimum Viable Product)
- [ ] Successfully fetch PRs from configured repositories every 4 hours
- [ ] Correctly attribute PRs to team members based on membership
- [ ] Store metrics in PostgreSQL with no duplicates
- [ ] Provide API endpoints that return accurate velocity and lead time data

### Quality Metrics
- [ ] 90%+ test coverage for core logic
- [ ] Zero data loss or duplication over 1 week of operation
- [ ] API response time < 500ms for typical queries
- [ ] Handle GitHub rate limits gracefully without failures

### User Validation
- [ ] Engineering managers can view their team's metrics
- [ ] Metrics match manual calculations (spot-check validation)
- [ ] Dashboard is intuitive and requires no training