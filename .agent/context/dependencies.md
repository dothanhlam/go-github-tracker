# Dependencies

## Core Dependencies

### `github.com/google/go-github/v58`
- **Purpose**: Official Go client library for GitHub API v3 and GraphQL API v4
- **Version**: v58 (latest stable)
- **Documentation**: https://pkg.go.dev/github.com/google/go-github/v58
- **Usage**: Fetch Pull Request data, authenticate with GitHub, handle pagination
- **Why**: Official library with comprehensive API coverage and active maintenance

### `github.com/jmoiron/sqlx`
- **Purpose**: Extensions to Go's standard `database/sql` library with struct mapping
- **Version**: v1.3.5+
- **Documentation**: https://pkg.go.dev/github.com/jmoiron/sqlx
- **Usage**: Database operations with automatic struct scanning, named queries
- **Why**: Reduces boilerplate while maintaining compatibility with `database/sql`

### `github.com/lib/pq`
- **Purpose**: Pure Go PostgreSQL driver for database/sql
- **Version**: v1.10.9+
- **Documentation**: https://pkg.go.dev/github.com/lib/pq
- **Usage**: PostgreSQL connection and query execution
- **Why**: Most popular and well-maintained PostgreSQL driver for Go

### `golang.org/x/oauth2`
- **Purpose**: OAuth 2.0 client implementation
- **Version**: Latest from golang.org/x
- **Documentation**: https://pkg.go.dev/golang.org/x/oauth2
- **Usage**: GitHub API authentication with Personal Access Tokens
- **Why**: Official Go OAuth2 library, required by go-github

### `github.com/aws/aws-lambda-go`
- **Purpose**: AWS Lambda runtime for Go
- **Version**: v1.41.0+
- **Documentation**: https://pkg.go.dev/github.com/aws/aws-lambda-go
- **Usage**: Lambda handler implementation, event processing
- **Why**: Official AWS SDK for Lambda functions written in Go

---

## Development Dependencies

### `github.com/stretchr/testify`
- **Purpose**: Testing toolkit with assertions and mocking
- **Version**: v1.8.4+
- **Documentation**: https://pkg.go.dev/github.com/stretchr/testify
- **Usage**: Unit tests, assertions, test suites, mocking
- **Why**: Industry standard for Go testing with rich assertion library

### `github.com/joho/godotenv`
- **Purpose**: Load environment variables from `.env` files
- **Version**: v1.5.1+
- **Documentation**: https://pkg.go.dev/github.com/joho/godotenv
- **Usage**: Local development configuration management
- **Why**: Simplifies local development by loading env vars from files

### `github.com/golang-migrate/migrate/v4`
- **Purpose**: Database migration tool
- **Version**: v4.16.0+
- **Documentation**: https://pkg.go.dev/github.com/golang-migrate/migrate/v4
- **Usage**: Versioned database schema migrations
- **Why**: Robust migration tool with PostgreSQL support (optional, may use raw SQL)

---

## Optional Dependencies (Future Phases)

### `github.com/go-chi/chi/v5`
- **Purpose**: Lightweight HTTP router
- **Version**: v5.0.10+
- **Documentation**: https://pkg.go.dev/github.com/go-chi/chi/v5
- **Usage**: API endpoint routing (Phase 3)
- **Why**: Fast, idiomatic, and compatible with standard library

### `github.com/rs/zerolog`
- **Purpose**: Structured logging library
- **Version**: v1.31.0+
- **Documentation**: https://pkg.go.dev/github.com/rs/zerolog
- **Usage**: JSON structured logging for CloudWatch
- **Why**: Zero-allocation, fast, and produces CloudWatch-friendly JSON logs

### `github.com/aws/aws-sdk-go-v2`
- **Purpose**: AWS SDK for Go v2
- **Version**: v1.21.0+
- **Documentation**: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2
- **Usage**: Secrets Manager integration, CloudWatch metrics
- **Why**: Official AWS SDK with better performance than v1

---

## Dependency Management

### Strategy
- **Use Go modules** (`go.mod` and `go.sum`) for all dependency management
- **Pin major versions** to avoid breaking changes
- **Review dependencies quarterly** for security updates
- **Minimize external dependencies** - prefer standard library when possible

### Adding New Dependencies
1. Evaluate if the standard library can solve the problem
2. Check the dependency's maintenance status and community support
3. Review the dependency's license (must be permissive)
4. Add to `go.mod` via `go get <package>@<version>`
5. Document the dependency in this file with justification

### Security
- Run `go mod tidy` regularly to remove unused dependencies
- Use `go list -m all` to audit all transitive dependencies
- Monitor for security vulnerabilities via GitHub Dependabot
- Keep dependencies up to date with patch releases

### Dependency Tree
```bash
# View dependency tree
go mod graph

# Check for available updates
go list -u -m all

# Update all dependencies to latest patch versions
go get -u=patch ./...
```
