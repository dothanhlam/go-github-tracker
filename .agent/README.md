# AI Agent Directory - DORA Metrics Tool

This directory contains AI-specific resources to help AI assistants understand and work with the DORA Metrics tracking tool project effectively.

## Project Overview

**go-github-tracker** is a DORA metrics tracking tool that collects Pull Request data from GitHub to measure team productivity and delivery performance. The system is designed for Engineering Managers who need visibility into team velocity and lead time metrics.

### Key Characteristics
- **Backend**: Go-based Lambda function
- **Data Source**: GitHub API (Pull Requests)
- **Storage**: PostgreSQL
- **Deployment**: AWS Lambda (ARM64) with EventBridge scheduling
- **Metrics**: Team velocity, DORA lead time (median & P95)

---

## Directory Structure

### `context/` - Project Context
Essential reading for understanding the project architecture and conventions.

- **`architecture.md`** - System design, component interactions, data flow
- **`conventions.md`** - Go coding standards and project-specific conventions
- **`dependencies.md`** - Core dependencies (go-github, sqlx, oauth2, mcp-go) and their purposes
- **`environments.md`** - Environment configuration (.env files, database switching)
- **`testing.md`** - Testing strategy, test structure, coverage statistics
- **`ci-cd.md`** - GitHub Actions workflow, test automation, deployment pipeline
- **`mcp-configuration.md`** - MCP server configuration for AI-powered insights
- **`test-exclusions.md`** - Packages excluded from CI tests and why
- **`review_metrics.md`** - PR review metrics schema and queries

### `planning/` - Planning Documents
Comprehensive planning and requirements documentation.

- **`plan.md`** - **START HERE** - Complete development plan with phases and tasks
- **`requirements.md`** - Functional and non-functional requirements (REQ-001 to REQ-021, NFR-001 to NFR-029)
- **`roadmap.md`** - Milestones, timeline, and future enhancements

### `workflows/` - Reusable Workflows
Step-by-step instructions for common development tasks.

- **`build.md`** - How to build the Go binary (local, production, cross-compilation)
- **`test.md`** - How to run tests, generate coverage, and organize test files

---

## AI Assistant Guidelines

### Before Making Changes
1. **Read `planning/plan.md`** - Understand the overall project structure and phases
2. **Check `planning/requirements.md`** - Ensure changes align with requirements
3. **Review `context/architecture.md`** - Understand system design and component boundaries
4. **Follow `context/conventions.md`** - Adhere to Go best practices and project standards

### When Implementing Features
1. **Refer to the current phase** in `planning/plan.md` - Focus on the active phase
2. **Update checkboxes** in `plan.md` as tasks are completed
3. **Document design decisions** in `context/architecture.md`
4. **Add new dependencies** to `context/dependencies.md` with justification
5. **Write tests first** - Follow TDD approach where applicable

### When Debugging
1. **Check requirements** - Ensure the expected behavior is documented
2. **Review architecture** - Understand component interactions
3. **Check workflows** - Follow established debugging procedures
4. **Update documentation** - Document any discovered gotchas or edge cases

### Documentation Updates
As the project evolves, keep these files current:
- ✅ Mark completed tasks in `plan.md`
- ✅ Document architectural decisions in `architecture.md`
- ✅ Add new dependencies to `dependencies.md`
- ✅ Update workflows as processes change

---

## Key Project Constraints

### Team Filtering
- Only track PRs from team members in the configured roster
- Support multi-team membership with weighted allocations
- Respect time-based membership (joined_at/left_at)

### Data Integrity
- All database operations must be idempotent (use `ON CONFLICT`)
- No duplicate PR entries allowed
- Support re-running without data corruption

### Performance
- Respect GitHub rate limits (primary & secondary)
- Use Go concurrency for multi-repository fetching
- API response time < 500ms

### Security
- Use OAuth2 for GitHub authentication
- Store secrets in AWS Secrets Manager (production)
- Never hardcode credentials

---

## Quick Reference

### Current Phase
**Phase 1: Foundation (Database & Setup)** - See `planning/plan.md` for details

### Primary Technologies
- Go 1.21+
- PostgreSQL (sqlx, lib/pq)
- GitHub API (google/go-github)
- AWS Lambda (ARM64)

### Key Files to Update
- `planning/plan.md` - Mark tasks complete as you progress
- `context/architecture.md` - Document design decisions
- `context/dependencies.md` - Track new dependencies
