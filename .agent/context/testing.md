# Testing Strategy

## Overview

The project uses Go's built-in testing framework with comprehensive unit tests covering core functionality.

## Test Structure

### Unit Tests (by Package)

#### `internal/collector/metrics_test.go`
**Coverage**: Metric calculation functions
- `TestCalculateCycleTime` - PR cycle time calculation
- `TestCalculateReviewTurnaround` - Review turnaround time
- `TestGetFirstReviewTime` - First review timestamp extraction
- `TestExtractReviewers` - Unique reviewer extraction
- `TestCountReviewsByState` - Review state counting

**Total**: 5 test functions, 22 sub-tests

#### `internal/config/config_test.go`
**Coverage**: Configuration parsing and validation
- `TestTeamConfigJSONParsing` - JSON parsing with various inputs
- `TestTeamConfigStructure` - Struct field validation
- `TestConfigValidation` - Database config validation
- `TestTeamMemberAllocation` - Allocation weight parsing

**Total**: 4 test functions, 13 sub-tests

#### `internal/database/migrations_test.go`
**Coverage**: Database migration system
- `TestSchemaMigrationsTable` - Migration tracking table creation
- `TestMigrationTracking` - Migration applied/recorded logic
- `TestMultipleMigrations` - Multiple migration handling
- `TestDatabaseConnection` - Connection validation
- `TestMigrationIdempotency` - Idempotent migration execution
- `TestRunMigrationsWithNoFiles` - Empty migration directory handling

**Total**: 6 test functions, 7 sub-tests

## Test Statistics

- **Total Test Functions**: 15
- **Total Sub-tests**: 42
- **Pass Rate**: 100%
- **Excluded Packages**: `internal/mcp` (incomplete, API compatibility issues)

## Running Tests

### All Tests (Excluding MCP)
```bash
go test $(go list ./internal/... | grep -v '/internal/mcp') -v
```

### With Coverage
```bash
go test $(go list ./internal/... | grep -v '/internal/mcp') -v -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out
```

### Specific Package
```bash
go test ./internal/collector -v
go test ./internal/config -v
go test ./internal/database -v
```

## CI/CD Integration

Tests run automatically via GitHub Actions on:
- Push to `main` branch
- Pull request creation/updates to `main`

See [`.github/workflows/test.yml`](file:///Users/lamdo/Projects/go-github-tracker/.github/workflows/test.yml)

## Test Patterns

### Table-Driven Tests
Most tests use table-driven approach for comprehensive coverage:
```go
tests := []struct {
    name    string
    input   interface{}
    want    interface{}
    wantErr bool
}{
    // test cases...
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test logic
    })
}
```

### Temporary Test Databases
Database tests use `t.TempDir()` for isolated test databases:
```go
tempDir := t.TempDir()
dbPath := filepath.Join(tempDir, "test.db")
```

## Future Testing

- [ ] Integration tests in `tests/` directory
- [ ] MCP server handler tests (after API fixes)
- [ ] End-to-end collector tests with mock GitHub API
- [ ] Performance benchmarks for large datasets
