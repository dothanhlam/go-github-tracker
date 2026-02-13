---
description: How to run tests
---

# Test Workflow

## Running Tests

1. **Run all tests**
   ```bash
# Run all tests except MCP (which is incomplete)
go test $(go list ./internal/... | grep -v '/internal/mcp') -v
   ```

2. **Run tests with coverage**
   ```bash
   go test -cover ./...
   ```

3. **Generate coverage report**
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out -o coverage.html
   ```

4. **Run tests with verbose output**
   ```bash
   go test -v ./...
   ```

## Test Organization

- Unit tests: `*_test.go` files alongside source code
- Integration tests: `tests/integration/` directory
- Test helpers: `tests/testutil/` directory

## Best Practices

- Write table-driven tests for multiple test cases
- Use `t.Parallel()` for independent tests
- Mock external dependencies
- Aim for >80% code coverage
- Test edge cases and error conditions
