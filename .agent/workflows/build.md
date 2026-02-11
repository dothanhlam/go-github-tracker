---
description: How to build the project
---

# Build Workflow

## Prerequisites

- Go 1.21+ installed
- Project dependencies downloaded (`go mod download`)

## Build Steps

1. **Clean previous builds**
   ```bash
   go clean
   ```

2. **Run tests**
   ```bash
   go test ./...
   ```

3. **Build the binary**
   ```bash
   go build -o bin/go-github-tracker ./cmd/...
   ```

4. **Verify the build**
   ```bash
   ./bin/go-github-tracker --version
   ```

## Build Options

### Development Build
```bash
go build -o bin/go-github-tracker ./cmd/...
```

### Production Build
```bash
go build -ldflags="-s -w" -o bin/go-github-tracker ./cmd/...
```

### Cross-Compilation
```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o bin/go-github-tracker-linux ./cmd/...

# For Windows
GOOS=windows GOARCH=amd64 go build -o bin/go-github-tracker.exe ./cmd/...
```
