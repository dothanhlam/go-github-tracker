package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// TeamMemberConfig represents a team member configuration
type TeamMemberConfig struct {
	Username   string  `json:"username"`
	Allocation float64 `json:"allocation"`
}

// TeamConfig represents a team configuration
type TeamConfig struct {
	TeamID  int                `json:"team_id"`
	Name    string             `json:"name"`
	Members []TeamMemberConfig `json:"members"`
}

// Config holds all application configuration
type Config struct {
	// Database configuration
	DBDriver string // "sqlite3" or "postgres"
	DBURL    string

	// GitHub configuration
	GitHubPAT string

	// Team configuration
	Teams []TeamConfig

	// Repositories to track
	Repositories []string
}

// Load loads configuration from environment variables
// It first attempts to load from a .env file if it exists
func Load() (*Config, error) {
	// Try to load .env file (ignore error if it doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		DBDriver:  getEnv("DB_DRIVER", "sqlite3"),
		DBURL:     getEnv("DB_URL", "./data/dora_metrics.db"),
		GitHubPAT: getEnv("GITHUB_PAT", ""),
	}

	// Parse team configuration
	teamConfigJSON := getEnv("TEAM_CONFIG_JSON", "[]")
	if err := json.Unmarshal([]byte(teamConfigJSON), &cfg.Teams); err != nil {
		return nil, fmt.Errorf("failed to parse TEAM_CONFIG_JSON: %w", err)
	}

	// Parse repositories
	reposStr := getEnv("REPOSITORIES", "")
	if reposStr != "" {
		cfg.Repositories = strings.Split(reposStr, ",")
		for i := range cfg.Repositories {
			cfg.Repositories[i] = strings.TrimSpace(cfg.Repositories[i])
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.DBDriver != "sqlite3" && c.DBDriver != "postgres" {
		return fmt.Errorf("DB_DRIVER must be 'sqlite3' or 'postgres', got: %s", c.DBDriver)
	}

	if c.DBURL == "" {
		return fmt.Errorf("DB_URL is required")
	}

	// GitHub PAT is optional for now (can be added later)
	// if c.GitHubPAT == "" {
	// 	return fmt.Errorf("GITHUB_PAT is required")
	// }

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
