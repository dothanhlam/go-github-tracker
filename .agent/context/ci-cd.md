# CI/CD Workflow

## GitHub Actions

### Test Workflow

**File**: `.github/workflows/test.yml`

**Triggers**:
- Push to `main` branch
- Pull requests to `main`
- Pull request updates

**Steps**:
1. Checkout code
2. Set up Go 1.21 with caching
3. Download dependencies
4. Run tests with coverage (excluding `internal/mcp`)
5. Generate coverage report
6. Upload to Codecov (optional)
7. Add test summary to GitHub UI

**Test Command**:
```bash
go test $(go list ./internal/... | grep -v '/internal/mcp') -v -coverprofile=coverage.out -covermode=atomic
```

**Coverage Report**:
```bash
go tool cover -func=coverage.out
```

## Test Exclusions

The `internal/mcp` package is excluded from CI tests because:
- Server structure is complete
- API compatibility with mcp-go SDK v0.43.2 needs fixes
- Handler signatures need updates

See [`.agent/context/test-exclusions.md`](file:///Users/lamdo/Projects/go-github-tracker/.agent/context/test-exclusions.md)

## Branch Protection (Recommended)

To enforce passing tests before merging:

1. Go to repository Settings → Branches
2. Add rule for `main` branch
3. Enable: "Require status checks to pass before merging"
4. Select: "Run Tests" workflow

## Future CI/CD

- [ ] Build ARM64 binary for Lambda deployment
- [ ] Deploy to AWS Lambda on successful tests
- [ ] Environment-specific deployments (dev/staging/prod)
- [ ] Automated database migrations
