# Multi-Environment Configuration

## Overview

The DORA Metrics Tracker supports multiple environment configurations through environment-specific `.env` files.

## Environment Files

| File | Purpose | Default |
|------|---------|---------|
| `.env.local` | Local development | ✅ Yes |
| `.env.test` | Testing (in-memory DB) | |
| `.env.staging` | Staging environment | |
| `.env.production` | Production | |
| `.env` | Production fallback | |

## How It Works

The application uses the `APP_ENV` environment variable to determine which configuration file to load:

1. If `APP_ENV` is not set, defaults to `"local"`
2. Loads `.env.{APP_ENV}` (e.g., `.env.local`)
3. For production, falls back to `.env` if `.env.production` doesn't exist

## Usage

### Local Development (Default)
```bash
# Uses .env.local automatically
./bin/collector
```

### Test Environment
```bash
# Uses .env.test
APP_ENV=test ./bin/collector
```

### Production Environment
```bash
# Uses .env.production or .env
APP_ENV=production ./bin/collector
```

## Setup

### First Time Setup
```bash
# Copy example to local config
cp .env.example .env.local

# Edit with your settings
vim .env.local
```

### Creating New Environments
```bash
# Create staging config
cp .env.example .env.staging

# Edit staging-specific settings
vim .env.staging
```

## Configuration Examples

### .env.local (Local Development)
```bash
DB_DRIVER=sqlite3
DB_URL=./data/dora_metrics.db
GITHUB_PAT=
TEAM_CONFIG_JSON=[]
REPOSITORIES=
```

### .env.test (Testing)
```bash
DB_DRIVER=sqlite3
DB_URL=:memory:
GITHUB_PAT=test_token
TEAM_CONFIG_JSON=[{"team_id":1,"name":"TestTeam","members":[{"username":"testuser","allocation":1.0}]}]
REPOSITORIES=testowner/testrepo
```

### .env.production (Production)
```bash
DB_DRIVER=postgres
DB_URL=postgres://user:password@db.example.com:5432/dora_metrics?sslmode=require
GITHUB_PAT=${GITHUB_PAT_FROM_SECRETS_MANAGER}
TEAM_CONFIG_JSON=[...]
REPOSITORIES=company/repo1,company/repo2
```

## Git Ignore

All `.env.*` files are ignored by git except `.env.example`:

```gitignore
.env
.env.*
!.env.example
```

## Best Practices

1. **Never commit** `.env.*` files (except `.env.example`)
2. **Use `.env.local`** for local development
3. **Use `.env.test`** for automated tests
4. **Use system environment variables** in production (AWS Lambda, etc.)
5. **Keep `.env.example`** updated with all required variables

## Troubleshooting

### Warning: .env.local not found
```bash
# Create the file
cp .env.example .env.local
```

### Using system environment variables
If no `.env.*` file is found, the app will use system environment variables. This is useful for:
- CI/CD pipelines
- Docker containers
- AWS Lambda (with environment variables)

## Notes

- **gvm users**: The user manages Go versions with `gvm`. Commands may need `source ~/.bash_profile` to initialize gvm before running Go commands.
