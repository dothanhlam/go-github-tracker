# DORA Metrics Tracker - Features

## Overview

The DORA Metrics Tracker is a comprehensive tool for tracking Pull Request metrics and team performance indicators using GitHub APIs. It provides insights into development velocity, code review quality, and team collaboration patterns.

---

## Core Features

### 1. Multi-Team PR Tracking

**Description**: Track Pull Requests across multiple teams with support for weighted team membership allocation.

**Capabilities**:
- ✅ Team-based PR filtering and attribution
- ✅ Weighted allocation for members across multiple teams
- ✅ Time-based team membership (joined_at/left_at)
- ✅ Automatic team membership management

**Use Cases**:
- Track PRs for platform teams, frontend teams, backend teams separately
- Handle developers who split time across multiple teams (e.g., 50% Platform, 50% Frontend)
- Historical team composition tracking

**Configuration**:
```json
{
  "team_id": 1,
  "name": "Platform",
  "members": [
    {"username": "alice", "allocation": 1.0},
    {"username": "bob", "allocation": 0.5}
  ]
}
```

---

### 2. DORA Metrics

**Description**: Standard DORA (DevOps Research and Assessment) metrics for measuring software delivery performance.

#### 2.1 Lead Time for Changes
**Definition**: Time from PR creation to merge  
**Calculation**: `merged_at - created_at`  
**View**: `view_dora_lead_time`

**Metrics Provided**:
- Median lead time (PostgreSQL)
- P95 lead time (PostgreSQL)
- Average lead time (SQLite)
- Min/max lead time
- PR count per month

**Query Example**:
```sql
SELECT * FROM view_dora_lead_time 
WHERE team_id = 1 
ORDER BY month DESC 
LIMIT 6;
```

#### 2.2 Deployment Frequency
**Definition**: PRs merged per week  
**View**: `view_team_velocity`

**Metrics Provided**:
- PRs merged per week
- Average cycle time
- Weekly trends

**Query Example**:
```sql
SELECT * FROM view_team_velocity 
WHERE team_id = 1 
ORDER BY week DESC 
LIMIT 12;
```

---

### 3. PR Review Metrics

**Description**: Advanced metrics for measuring code review quality, engagement, and collaboration.

#### 3.1 Review Turnaround Time
**Definition**: Time from PR creation to first review  
**Fields**: `first_review_at`, `review_turnaround_hours`  
**View**: `view_review_turnaround`

**Metrics Provided**:
- Average turnaround time
- Min/max turnaround time
- Median and P95 (PostgreSQL)
- PRs reviewed within 24 hours
- PRs taking over 24 hours

**Use Cases**:
- Identify review bottlenecks
- Measure team responsiveness
- Track SLA compliance

**Query Example**:
```sql
SELECT 
    month,
    avg_turnaround_hours,
    within_24h_count,
    over_24h_count
FROM view_review_turnaround
WHERE team_id = 1;
```

#### 3.2 Comment-to-PR Ratio
**Definition**: Number of review comments and conversations per PR  
**Fields**: `review_comments_count`, `conversation_count`  
**View**: `view_review_engagement`

**Metrics Provided**:
- Average comments per PR
- Average conversations per PR
- Average reviewers per PR

**Use Cases**:
- Measure review thoroughness
- Identify overly complex PRs
- Track review engagement quality

#### 3.3 PR Rejection/Request Changes Rate
**Definition**: Percentage of PRs requiring changes vs immediate approval  
**Fields**: `changes_requested_count`, `approved_count`  
**View**: `view_review_engagement`

**Metrics Provided**:
- Changes requested rate (%)
- Approval rate (%)
- Review outcome distribution

**Use Cases**:
- Measure code quality
- Track review rigor
- Identify training opportunities

#### 3.4 Knowledge Sharing Depth
**Definition**: Unique reviewers and cross-team participation  
**Fields**: `reviewers_count`, `external_reviewers_count`, `reviewers_list`  
**View**: `view_knowledge_sharing`

**Metrics Provided**:
- Average reviewers per PR
- Average external reviewers (outside team)
- External reviewer rate (%)
- Total external reviews
- PRs with external reviews

**Use Cases**:
- Measure knowledge distribution
- Track cross-team collaboration
- Identify knowledge silos
- Encourage pair programming and mentorship

**Query Example**:
```sql
SELECT 
    month,
    avg_external_reviewers,
    external_reviewer_rate,
    prs_with_external_reviews
FROM view_knowledge_sharing
WHERE team_id = 1;
```

---

### 4. Automated Data Collection

**Description**: Scheduled data collection from GitHub APIs (Phase 2 - Planned).

**Capabilities**:
- ⏳ Automated PR data fetching every 4 hours
- ⏳ GitHub API authentication with Personal Access Tokens
- ⏳ Rate limit handling
- ⏳ Incremental updates (fetch only new/updated PRs)
- ⏳ Idempotent database operations

**Deployment**:
- AWS Lambda function
- EventBridge scheduler (every 4 hours)
- Secrets Manager for GitHub PAT

---

### 5. Multi-Environment Configuration

**Description**: Support for different environment configurations.

**Capabilities**:
- ✅ Environment-specific `.env` files
- ✅ Local development (`.env.local`)
- ✅ Testing (`.env.test`)
- ✅ Staging (`.env.staging`)
- ✅ Production (`.env.production`)

**Usage**:
```bash
# Local development (default)
./bin/collector

# Test environment
APP_ENV=test ./bin/collector

# Production
APP_ENV=production ./bin/collector
```

**Configuration Variables**:
- `DB_DRIVER` - Database driver (sqlite3 or postgres)
- `DB_URL` - Database connection string
- `GITHUB_PAT` - GitHub Personal Access Token
- `TEAM_CONFIG_JSON` - Team configuration (JSON)
- `REPOSITORIES` - Repositories to track (comma-separated)

---

### 6. Dual Database Support

**Description**: Support for both SQLite (local development) and PostgreSQL (production).

**SQLite Features**:
- ✅ In-memory database for testing (`:memory:`)
- ✅ File-based database for local development
- ✅ No external dependencies
- ✅ Fast local iteration

**PostgreSQL Features**:
- ✅ Production-grade performance
- ✅ JSONB support for complex data
- ✅ Advanced analytics (PERCENTILE_CONT)
- ✅ Scalability for large datasets

**Automatic Driver Selection**:
```bash
# SQLite (local)
DB_DRIVER=sqlite3
DB_URL=./data/dora_metrics.db

# PostgreSQL (production)
DB_DRIVER=postgres
DB_URL=postgres://user:pass@host:5432/dbname
```

---

### 7. Database Schema & Migrations

**Description**: Robust database schema with automated migrations.

**Tables**:
1. **teams** - Team definitions
2. **team_memberships** - Team member allocations with time-based tracking
3. **pr_metrics** - Pull Request metrics (21 columns)

**Views**:
1. **view_team_velocity** - Weekly PR velocity
2. **view_dora_lead_time** - DORA lead time metrics
3. **view_review_turnaround** - Review response time
4. **view_review_engagement** - Review quality metrics
5. **view_knowledge_sharing** - Collaboration metrics

**Migration System**:
- ✅ Automatic migration execution on startup
- ✅ Driver-specific migrations (SQLite vs PostgreSQL)
- ✅ Sequential execution (001, 002, 003...)
- ✅ Idempotent operations

---

### 8. Analytics & Reporting

**Description**: Pre-built SQL views for common analytics queries.

**Available Reports**:

#### Team Velocity Dashboard
```sql
SELECT * FROM view_team_velocity 
WHERE team_id = 1 
ORDER BY week DESC;
```

#### DORA Metrics Summary
```sql
SELECT * FROM view_dora_lead_time 
WHERE team_id = 1 
ORDER BY month DESC;
```

#### Review Quality Report
```sql
SELECT 
    e.month,
    e.avg_comments_per_pr,
    e.changes_requested_rate,
    t.avg_turnaround_hours
FROM view_review_engagement e
JOIN view_review_turnaround t USING (team_id, month)
WHERE e.team_id = 1;
```

#### Knowledge Sharing Report
```sql
SELECT * FROM view_knowledge_sharing 
WHERE team_id = 1 
ORDER BY month DESC;
```

---

## Technical Features

### Data Models

**Comprehensive Go Structs**:
- `Team` - Team information
- `TeamMembership` - Member allocations
- `PRMetric` - PR metrics (21 fields)
- `TeamVelocity` - Velocity view model
- `DORALeadTime` - DORA metrics view model
- `ReviewTurnaround` - Review turnaround view model
- `ReviewEngagement` - Review engagement view model
- `KnowledgeSharing` - Knowledge sharing view model

### Database Indexes

**Optimized for Common Queries**:
- Team ID lookups
- Date range queries
- Author searches
- Review metrics filtering

**Total Indexes**: 10+

---

## Planned Features (Phase 2+)

### Phase 2: Data Collection
- ⏳ GitHub API integration
- ⏳ PR data fetcher
- ⏳ Review data fetcher
- ⏳ Comment data fetcher
- ⏳ Automated scheduling

### Phase 3: REST API
- ⏳ RESTful API endpoints
- ⏳ API key authentication
- ⏳ JSON response formatting
- ⏳ Filtering and pagination

### Phase 4: Deployment
- ⏳ AWS Lambda deployment
- ⏳ EventBridge scheduling
- ⏳ Secrets Manager integration
- ⏳ CloudWatch monitoring

### Future Enhancements
- ⏳ Change Failure Rate tracking
- ⏳ Mean Time to Recovery (MTTR)
- ⏳ Deployment tracking
- ⏳ Incident correlation
- ⏳ Web dashboard UI
- ⏳ Slack/email notifications
- ⏳ Custom metric definitions

---

## Use Cases

### Engineering Managers
- Track team velocity and productivity
- Identify review bottlenecks
- Measure code quality trends
- Monitor cross-team collaboration

### Team Leads
- Optimize review processes
- Balance workload across team members
- Identify knowledge silos
- Improve onboarding effectiveness

### Individual Contributors
- Understand personal metrics
- Compare against team averages
- Track improvement over time
- Identify areas for growth

### Organizations
- Benchmark across teams
- Identify best practices
- Data-driven process improvements
- Executive reporting

---

## Benefits

### Visibility
- 📊 Real-time metrics on team performance
- 📈 Historical trend analysis
- 🎯 Identify bottlenecks and inefficiencies

### Quality
- ✅ Measure review thoroughness
- 🔍 Track code quality indicators
- 📝 Monitor review engagement

### Collaboration
- 🤝 Track cross-team participation
- 💡 Identify knowledge sharing patterns
- 🌟 Encourage mentorship

### Efficiency
- ⚡ Automated data collection
- 🔄 Idempotent operations
- 📦 Easy deployment (AWS Lambda)

---

## Getting Started

See [README.md](file:///Users/lamdo/Projects/go-github-tracker/README.md) for setup instructions.

For detailed metric definitions, see [.agent/context/review_metrics.md](file:///Users/lamdo/Projects/go-github-tracker/.agent/context/review_metrics.md).

For architecture details, see [.agent/context/architecture.md](file:///Users/lamdo/Projects/go-github-tracker/.agent/context/architecture.md).
