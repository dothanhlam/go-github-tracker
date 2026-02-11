package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all migration files for the given driver
func (db *DB) RunMigrations() error {
	// Determine migration directory based on driver
	var migrationDir string
	switch db.driver {
	case "sqlite3":
		migrationDir = "migrations/sqlite"
	case "postgres":
		migrationDir = "migrations/postgres"
	default:
		return fmt.Errorf("unsupported driver: %s", db.driver)
	}

	// Read migration files from the filesystem (not embedded for now)
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	// Sort files to ensure they run in order
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Execute each migration file
	for _, filename := range sqlFiles {
		filePath := filepath.Join(migrationDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Execute the SQL
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		fmt.Printf("✓ Applied migration: %s\n", filename)
	}

	return nil
}
