# Requirements

## Functional Requirements

### Data Collection
- **REQ-001**: The system shall fetch Pull Request data from configured GitHub repositories using the GitHub API
- **REQ-002**: The system shall filter PRs to only include those authored by team members in the configured team roster
- **REQ-003**: The system shall support tracking multiple teams simultaneously
- **REQ-004**: The system shall handle team members who belong to multiple teams with weighted allocations (0.00 to 1.00)
- **REQ-005**: The system shall respect time-based team membership (joined_at and left_at dates)
- **REQ-006**: The system shall run automatically every 4 hours via AWS EventBridge
- **REQ-007**: The system shall calculate cycle time as the duration between PR creation and merge

### Data Storage
- **REQ-008**: The system shall store team configurations in a PostgreSQL database
- **REQ-009**: The system shall store team membership data with allocation weights
- **REQ-010**: The system shall store PR metrics including cycle time, state, and timestamps
- **REQ-011**: The system shall prevent duplicate PR entries using unique constraints
- **REQ-012**: The system shall support idempotent operations (re-running without creating duplicates)

### Metrics & Reporting
- **REQ-013**: The system shall calculate team velocity (PRs merged per week)
- **REQ-014**: The system shall calculate DORA lead time metrics (median and P95)
- **REQ-015**: The system shall provide SQL views for easy metric consumption
- **REQ-016**: The system shall expose API endpoints for retrieving team velocity
- **REQ-017**: The system shall expose API endpoints for retrieving DORA lead time metrics
- **REQ-018**: The system shall return metrics in JSON format

### Team Management
- **REQ-019**: The system shall support adding/removing team members
- **REQ-020**: The system shall support updating team member allocation weights
- **REQ-021**: The system shall track historical team membership changes

## Non-Functional Requirements

### Performance
- **NFR-001**: API endpoints shall respond within 500ms for typical queries
- **NFR-002**: The system shall process PRs from multiple repositories concurrently
- **NFR-003**: The system shall respect GitHub's primary and secondary rate limits
- **NFR-004**: Database queries shall use connection pooling for efficiency
- **NFR-005**: The system shall complete a full data collection cycle within 15 minutes

### Security
- **NFR-006**: GitHub authentication shall use OAuth2 tokens
- **NFR-007**: GitHub Personal Access Tokens shall be stored in AWS Secrets Manager (production)
- **NFR-008**: API endpoints shall require authentication (API key)
- **NFR-009**: Database credentials shall not be hardcoded in source code
- **NFR-010**: All sensitive configuration shall be passed via environment variables

### Reliability
- **NFR-011**: The system shall implement retry logic for transient GitHub API failures
- **NFR-012**: The system shall use exponential backoff for rate limit errors
- **NFR-013**: Failed processing attempts shall be logged to a dead letter queue
- **NFR-014**: The system shall achieve 99% uptime for scheduled data collection
- **NFR-015**: Database operations shall use transactions to ensure data consistency

### Scalability
- **NFR-016**: The system shall support tracking up to 50 repositories
- **NFR-017**: The system shall support up to 10 teams with 100 members each
- **NFR-018**: The system shall handle up to 10,000 PRs per collection cycle
- **NFR-019**: Database schema shall support horizontal scaling if needed

### Maintainability
- **NFR-020**: Code shall achieve 90%+ test coverage for core logic
- **NFR-021**: All exported functions shall have documentation comments
- **NFR-022**: The system shall use structured logging in JSON format
- **NFR-023**: Database migrations shall be versioned and reversible
- **NFR-024**: The system shall follow Go best practices and conventions

### Observability
- **NFR-025**: The system shall emit CloudWatch metrics for PRs processed per run
- **NFR-026**: The system shall emit CloudWatch metrics for API errors
- **NFR-027**: The system shall emit CloudWatch metrics for database query duration
- **NFR-028**: CloudWatch alarms shall be configured for critical failures
- **NFR-029**: All errors shall be logged with sufficient context for debugging

## Constraints

### Technical Constraints
- **CON-001**: Must use Go as the primary programming language
- **CON-002**: Must use PostgreSQL for data persistence
- **CON-003**: Must deploy to AWS Lambda (ARM64 architecture)
- **CON-004**: Must use GitHub API v3 or v4 (GraphQL)
- **CON-005**: Lambda execution time must not exceed 15 minutes (AWS limit)

### Business Constraints
- **CON-006**: Only track PRs from team members in the configured roster
- **CON-007**: Must handle shared repositories where multiple teams contribute
- **CON-008**: Must support fractional team membership for accurate attribution

### Operational Constraints
- **CON-009**: Data collection frequency is fixed at every 4 hours
- **CON-010**: GitHub API rate limits must be strictly respected
- **CON-011**: Must minimize AWS Lambda costs (hence ARM64 architecture)

## Assumptions

### Technical Assumptions
- **ASM-001**: GitHub API will remain stable and backward compatible
- **ASM-002**: PostgreSQL database will be available and accessible from Lambda
- **ASM-003**: AWS Lambda cold start times are acceptable for scheduled jobs
- **ASM-004**: GitHub Personal Access Token has sufficient permissions to read repository data

### Business Assumptions
- **ASM-005**: Team membership data will be manually maintained initially
- **ASM-006**: Engineering managers will access metrics via API or SQL queries initially
- **ASM-007**: PR cycle time is a sufficient proxy for DORA lead time metric
- **ASM-008**: Merged PRs are the primary indicator of team velocity

### Data Assumptions
- **ASM-009**: PR creation and merge timestamps are accurate in GitHub
- **ASM-010**: Team members use consistent GitHub usernames
- **ASM-011**: Historical PR data is available via GitHub API
- **ASM-012**: Allocation weights sum to 1.00 or less per team member
