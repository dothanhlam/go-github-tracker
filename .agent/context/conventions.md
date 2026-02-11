# Coding Conventions

## Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

## Project Conventions

### File Organization

- Keep related code together in packages
- Use clear, descriptive package names
- Avoid circular dependencies

### Naming Conventions

- Use `camelCase` for unexported identifiers
- Use `PascalCase` for exported identifiers
- Use descriptive names that convey purpose
- Avoid abbreviations unless widely understood

### Error Handling

- Always check and handle errors
- Provide context when wrapping errors
- Use custom error types for domain-specific errors

### Testing

- Write tests for all public APIs
- Use table-driven tests where appropriate
- Keep tests close to the code they test (`*_test.go`)

### Comments

- Document all exported functions, types, and constants
- Use complete sentences in comments
- Explain "why" not "what" in implementation comments

## Code Review Checklist

- [ ] Code follows Go conventions
- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Error handling is appropriate
- [ ] No unnecessary complexity
