# Architecture

## Overview

The **go-github-tracker** is a serverless DORA metrics tracking system that collects Pull Request data from GitHub and stores it in PostgreSQL for analysis. The system is designed for automated, scheduled execution with a focus on data integrity, idempotency, and cost efficiency.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      AWS Cloud                              │
│                                                             │
│  ┌──────────────┐                                          │
│  │ EventBridge  │  Triggers every 4 hours                  │
│  │  Scheduler   │                                          │
│  └──────┬───────┘                                          │
│         │                                                   │
│         ▼                                                   │
│  ┌──────────────────────────────────────────┐             │
│  │      Lambda Function (ARM64)             │             │
│  │  ┌────────────────────────────────────┐  │             │
│  │  │  Go Collector Application          │  │             │
│  │  │  • Config Loader                   │  │             │
│  │  │  • GitHub Client (OAuth2)          │  │             │
│  │  │  • Team Membership Sync            │  │             │
│  │  │  • PR Fetcher (Concurrent)         │  │             │
│  │  │  • Filter & Attribution Engine     │  │             │
│  │  │  • Database Persistence            │  │             │
│  │  └────────────────────────────────────┘  │             │
│  └──────────────┬───────────────────────────┘             │
│                 │                                           │
│                 │ Reads                                     │
│                 ▼                                           │
│  ┌──────────────────────────┐                              │
│  │   Secrets Manager        │                              │
│  │   • GITHUB_PAT           │                              │
│  │   • DB_URL               │                              │
│  └──────────────────────────┘                              │
│                                                             │
│  ┌──────────────────────────┐                              │
│  │   CloudWatch             │                              │
│  │   • Logs (JSON)          │                              │
│  │   • Metrics              │                              │
│  │   • Alarms               │                              │
│  └──────────────────────────┘                              │
└─────────────────────────────────────────────────────────────┘
                 │
                 │ Writes PR Metrics
                 ▼
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL Database (RDS/External)             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Tables:                                            │   │
│  │  • teams                                            │   │
│  │  • team_memberships                                 │   │
│  │  • pr_metrics                                       │   │
│  │                                                      │   │
│  │  Views:                                             │   │
│  │  • view_team_velocity                               │   │
│  │  • view_dora_lead_time                              │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                 │
                 │ Queries
                 ▼
┌─────────────────────────────────────────────────────────────┐
│              API Layer (Future - Phase 3)                   │
│  ┌──────────────────────────────────────────┐              │
│  │  Lambda Function (API Handler)           │              │
│  │  • GET /metrics/team/{id}/velocity       │              │
│  │  • GET /metrics/team/{id}/lead-time      │              │
│  └──────────────────────────────────────────┘              │
└─────────────────────────────────────────────────────────────┘
                 │
                 ▼
           External Clients
        (Engineering Managers)
```

---

## System Components

### 1. Collector (Lambda Function)

The core component that orchestrates data collection.

#### Responsibilities
- Load configuration from environment variables
- Authenticate with GitHub API using OAuth2
- Fetch team membership data from database
- Poll GitHub for Pull Requests from configured repositories
- Filter PRs by team membership
- Calculate metrics (cycle time)
- Persist PR data to PostgreSQL with idempotent operations

#### Key Modules

**Config Module** (`internal/config`)
- Load and validate environment variables
- Parse team configuration JSON
- Provide configuration to other modules

**GitHub Client** (`internal/github`)
- Wrap `google/go-github` library
- Handle OAuth2 authentication
- Implement rate limit handling and backoff
- Fetch PRs with pagination support
- Concurrent repository fetching

**Team Membership Module** (`internal/membership`)
- Load active team members from database
- Build in-memory lookup map: `username -> []TeamMembership`
- Handle time-based membership (joined_at/left_at)
- Support weighted allocations

**Collector Module** (`internal/collector`)
- Main orchestration logic
- Coordinate GitHub fetching and database persistence
- Implement filtering and attribution logic
- Error handling and retry logic

**Database Module** (`internal/database`)
- Database connection management (connection pooling)
- CRUD operations for teams, memberships, and PR metrics
- Idempotent upsert operations using `ON CONFLICT`
- Transaction management

---

### 2. Database (PostgreSQL)

Persistent storage for team configurations and PR metrics.

#### Schema Design

**teams** - Team registry
```sql
id, name, description, created_at, updated_at
```

**team_memberships** - Many-to-many team membership with time ranges
```sql
id, team_id, github_username, allocation_weight, joined_at, left_at, created_at
UNIQUE(team_id, github_username, joined_at)
```

**pr_metrics** - Flattened PR data for fast querying
```sql
id, team_id, pr_number, repository, author, title, 
created_at, merged_at, closed_at, cycle_time_hours, state, created_date
UNIQUE(team_id, repository, pr_number)
```

#### Views

**view_team_velocity** - Weekly PR merge velocity
```sql
team_id, week, prs_merged, avg_cycle_time_hours
```

**view_dora_lead_time** - Monthly lead time metrics
```sql
team_id, month, median_lead_time_hours, p95_lead_time_hours
```

---

### 3. API Layer (Future - Phase 3)

REST API for programmatic access to metrics.

#### Endpoints
- `GET /metrics/team/{team_id}/velocity` - Team velocity over time
- `GET /metrics/team/{team_id}/lead-time` - DORA lead time metrics

#### Authentication
- API key-based authentication
- Keys stored in database or Secrets Manager

---

## Data Flow

### Collection Flow (Every 4 Hours)

```
1. EventBridge triggers Lambda
   ↓
2. Lambda starts, loads configuration
   ↓
3. Authenticate with GitHub (OAuth2)
   ↓
4. Load team memberships from database
   ↓
5. For each configured repository (concurrent):
   a. Fetch closed/merged PRs (paginated)
   b. Filter PRs by team membership
   c. Calculate cycle time
   d. Attribute to teams (weighted)
   ↓
6. Batch upsert PR metrics to database
   ↓
7. Log results to CloudWatch
   ↓
8. Lambda terminates
```

### Query Flow (Future - API)

```
1. Client sends GET request to API Gateway
   ↓
2. API Gateway triggers Lambda (API handler)
   ↓
3. Lambda validates API key
   ↓
4. Lambda queries PostgreSQL view
   ↓
5. Lambda returns JSON response
   ↓
6. Client receives metrics
```

---

## Key Design Decisions

### 1. Serverless Architecture (AWS Lambda)
**Decision**: Use AWS Lambda instead of long-running servers  
**Rationale**:
- Cost-efficient for scheduled jobs (pay per execution)
- No server management overhead
- Automatic scaling (not needed, but future-proof)
- ARM64 architecture for 20% cost savings

**Trade-offs**:
- Cold start latency (acceptable for scheduled jobs)
- 15-minute execution time limit (sufficient for our use case)

---

### 2. Idempotent Database Operations
**Decision**: Use `ON CONFLICT` for all inserts/updates  
**Rationale**:
- Safe to re-run collector without data corruption
- Handles Lambda retries gracefully
- Simplifies error recovery

**Implementation**:
```sql
INSERT INTO pr_metrics (team_id, repository, pr_number, ...)
VALUES (?, ?, ?, ...)
ON CONFLICT (team_id, repository, pr_number)
DO UPDATE SET 
    merged_at = EXCLUDED.merged_at,
    cycle_time_hours = EXCLUDED.cycle_time_hours,
    ...;
```

---

### 3. Flattened PR Metrics Table
**Decision**: Store PR metrics in a denormalized table instead of normalized relational model  
**Rationale**:
- Fast queries for analytics (no joins needed)
- Supports multi-team attribution (same PR can appear for multiple teams)
- Optimized for read-heavy workload

**Trade-offs**:
- Data duplication (same PR stored multiple times if multi-team)
- Larger storage footprint (acceptable for our scale)

---

### 4. In-Memory Team Membership Map
**Decision**: Load all team memberships into memory at startup  
**Rationale**:
- Fast lookup during PR filtering (O(1) instead of database query per PR)
- Team membership data is small (< 1000 members expected)
- Reduces database load

**Trade-offs**:
- Stale data if memberships change mid-execution (acceptable, next run will pick up changes)

---

### 5. Concurrent Repository Fetching
**Decision**: Use goroutines to fetch PRs from multiple repositories concurrently  
**Rationale**:
- Reduces total execution time
- GitHub API allows concurrent requests
- Maximizes Lambda efficiency

**Implementation**:
- Use `sync.WaitGroup` for coordination
- Respect GitHub rate limits with semaphore pattern
- Collect errors and handle gracefully

---

### 6. SQL Views for Metrics
**Decision**: Use PostgreSQL views instead of application-level aggregation  
**Rationale**:
- Leverage database's query optimization
- Reusable across different clients (API, SQL clients, BI tools)
- Easier to test and validate

**Trade-offs**:
- Less flexible than application code (but sufficient for our needs)

---

## Technology Stack Rationale

### Go
- **Why**: Fast, compiled, excellent concurrency support, small binary size
- **Alternatives considered**: Python (slower, larger Lambda package), Node.js (less type-safe)

### PostgreSQL
- **Why**: Robust, supports advanced SQL features (percentiles, window functions), widely used
- **Alternatives considered**: DynamoDB (less flexible for analytics), MySQL (weaker analytics features)

### AWS Lambda (ARM64)
- **Why**: Cost-efficient, serverless, ARM64 for 20% savings
- **Alternatives considered**: ECS (overkill for scheduled jobs), EC2 (requires management)

### google/go-github
- **Why**: Official library, comprehensive API coverage, active maintenance
- **Alternatives considered**: Direct REST API calls (more boilerplate)

---

## Scalability Considerations

### Current Design Supports
- Up to 50 repositories
- Up to 10 teams with 100 members each
- Up to 10,000 PRs per collection cycle

### Future Scaling Options
- **Horizontal scaling**: Split repositories across multiple Lambda invocations
- **Database scaling**: PostgreSQL read replicas for API queries
- **Caching**: Add Redis for frequently accessed metrics
- **Partitioning**: Partition `pr_metrics` table by date for faster queries

---

## Security Considerations

### Authentication & Authorization
- GitHub PAT stored in AWS Secrets Manager (production)
- Database credentials passed via environment variables (encrypted at rest)
- API endpoints require API key authentication

### Data Privacy
- Only collect public PR metadata (no code content)
- Team membership data is internal only
- No PII stored beyond GitHub usernames

### Network Security
- Lambda runs in VPC (if database is in private subnet)
- Database accessible only from Lambda security group
- API Gateway with throttling and WAF (future)

---

## Monitoring & Observability

### Logging
- Structured JSON logs to CloudWatch
- Log levels: DEBUG, INFO, WARN, ERROR
- Include correlation IDs for tracing

### Metrics
- PRs processed per run
- API errors (rate limits, auth failures)
- Database query duration
- Lambda execution time

### Alarms
- Lambda execution failures
- GitHub API rate limit exceeded
- Database connection failures
- Unexpected data anomalies (e.g., zero PRs fetched)
