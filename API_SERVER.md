# API Server

REST API for accessing DORA metrics programmatically.

## Quick Start

### 1. Build the Server

```bash
go build -o bin/api-server cmd/api-server/main.go
```

### 2. Configure Environment

```bash
# Required
export DB_DRIVER=sqlite3
export DB_URL=./data/metrics.db

# API Keys (comma-separated)
export API_KEYS=dev-key-123,prod-key-456

# Optional
export PORT=8080  # Default: 8080
```

### 3. Run the Server

```bash
./bin/api-server
```

Or with go run:
```bash
API_KEYS=test-key go run cmd/api-server/main.go
```

---

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
All endpoints (except `/health`) require an API key in the `X-API-Key` header:

```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/v1/teams
```

---

## Endpoints

### Health Check
```
GET /api/v1/health
```

**No authentication required**

**Response**:
```json
{
  "status": "healthy",
  "database": "connected",
  "version": "1.0.0"
}
```

**Example**:
```bash
curl http://localhost:8080/api/v1/health
```

---

### List Teams
```
GET /api/v1/teams
```

**Response**:
```json
{
  "teams": [
    {
      "id": 1,
      "name": "Engineering Team",
      "member_count": 8
    }
  ]
}
```

**Example**:
```bash
curl -H "X-API-Key: test-key" http://localhost:8080/api/v1/teams
```

---

### Team Velocity
```
GET /api/v1/teams/{id}/velocity
```

**Query Parameters**:
- `start_date` (optional): ISO 8601 date, default: 30 days ago
- `end_date` (optional): ISO 8601 date, default: today
- `granularity` (optional): `day`, `week`, `month`, default: `week`

**Response**:
```json
{
  "team_id": 1,
  "team_name": "Engineering Team",
  "period": {
    "start": "2026-01-15T00:00:00Z",
    "end": "2026-02-15T00:00:00Z"
  },
  "metrics": [
    {
      "period": "2026-02-09",
      "prs_merged": 12,
      "avg_cycle_time_hours": 18.5
    }
  ]
}
```

**Example**:
```bash
curl -H "X-API-Key: test-key" \
  "http://localhost:8080/api/v1/teams/1/velocity?start_date=2026-01-01&end_date=2026-02-15"
```

---

### DORA Lead Time
```
GET /api/v1/teams/{id}/lead-time
```

**Query Parameters**: Same as velocity

**Response**:
```json
{
  "team_id": 1,
  "team_name": "Engineering Team",
  "period": {
    "start": "2026-01-15T00:00:00Z",
    "end": "2026-02-15T00:00:00Z"
  },
  "metrics": [
    {
      "period": "2026-02",
      "median_lead_time_hours": 24.0,
      "p95_lead_time_hours": 72.0
    }
  ]
}
```

**Example**:
```bash
curl -H "X-API-Key: test-key" \
  http://localhost:8080/api/v1/teams/1/lead-time
```

---

### Review Turnaround
```
GET /api/v1/teams/{id}/review-turnaround
```

**Query Parameters**: `start_date`, `end_date`

**Response**:
```json
{
  "team_id": 1,
  "team_name": "Engineering Team",
  "period": {
    "start": "2026-01-15T00:00:00Z",
    "end": "2026-02-15T00:00:00Z"
  },
  "metrics": [
    {
      "period": "2026-02-09",
      "avg_turnaround_hours": 4.2,
      "median_turnaround_hours": 3.0
    }
  ]
}
```

**Example**:
```bash
curl -H "X-API-Key: test-key" \
  http://localhost:8080/api/v1/teams/1/review-turnaround
```

---

### Review Engagement
```
GET /api/v1/teams/{id}/review-engagement
```

**Response**:
```json
{
  "team_id": 1,
  "team_name": "Engineering Team",
  "period": {...},
  "metrics": [
    {
      "period": "2026-02-09",
      "total_reviews": 45,
      "unique_reviewers": 8,
      "avg_reviews_per_pr": 3.2
    }
  ]
}
```

**Example**:
```bash
curl -H "X-API-Key: test-key" \
  http://localhost:8080/api/v1/teams/1/review-engagement
```

---

### Knowledge Sharing
```
GET /api/v1/teams/{id}/knowledge-sharing
```

**Response**:
```json
{
  "team_id": 1,
  "team_name": "Engineering Team",
  "period": {...},
  "metrics": [
    {
      "period": "2026-02-09",
      "cross_team_reviews": 12,
      "knowledge_sharing_score": 0.75
    }
  ]
}
```

**Example**:
```bash
curl -H "X-API-Key: test-key" \
  http://localhost:8080/api/v1/teams/1/knowledge-sharing
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Team not found",
    "status": 404
  }
}
```

**Error Codes**:
- `BAD_REQUEST` (400): Invalid request parameters
- `UNAUTHORIZED` (401): Missing or invalid API key
- `NOT_FOUND` (404): Resource not found
- `INTERNAL_ERROR` (500): Server error

---

## Testing

### Manual Testing

```bash
# Start server
API_KEYS=test-key go run cmd/api-server/main.go

# In another terminal, test endpoints
curl http://localhost:8080/api/v1/health

curl -H "X-API-Key: test-key" \
  http://localhost:8080/api/v1/teams

curl -H "X-API-Key: test-key" \
  http://localhost:8080/api/v1/teams/1/velocity
```

### Load Testing

```bash
# Install hey
brew install hey

# Test endpoint
hey -n 1000 -c 10 \
  -H "X-API-Key: test-key" \
  http://localhost:8080/api/v1/teams
```

---

## Deployment

### Standalone Server

Run as a systemd service or Docker container:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-server cmd/api-server/main.go

FROM alpine:latest
COPY --from=builder /app/api-server /api-server
EXPOSE 8080
CMD ["/api-server"]
```

### Lambda + API Gateway

See `terraform/` directory for Lambda deployment configuration.

---

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_DRIVER` | Database driver (`sqlite3`, `postgres`) | - | Yes |
| `DB_URL` | Database connection URL | - | Yes |
| `API_KEYS` | Comma-separated API keys | - | Yes |
| `PORT` | HTTP server port | `8080` | No |

### CORS

CORS is configured in `internal/api/server.go`. For production, update:

```go
AllowedOrigins: []string{"https://your-dashboard.com"},
```

---

## Architecture

```
cmd/api-server/
  main.go                 # Entry point

internal/api/
  server.go              # Server setup, routing
  handlers/
    health.go            # Health check
    teams.go             # Team endpoints
  middleware/
    auth.go              # API key authentication
    logging.go           # Request logging
  response/
    response.go          # Response helpers

internal/service/
  metrics.go             # Business logic
```

---

## Next Steps

1. ✅ API server implemented
2. Add unit tests for handlers
3. Add integration tests
4. Deploy to production
5. Build dashboard (Phase 5)
