package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all migration files for the given driver.
// This function implements a migration tracking system to ensure each migration
// runs exactly once, preventing errors like "duplicate column name" when the
// application restarts.
//
// Migration Tracking Flow:
// 1. Create schema_migrations table (if it doesn't exist)
// 2. For each migration file:
//    - Check if it's already been applied
//    - Skip if already applied
//    - Execute and record if not applied
func (db *DB) RunMigrations() error {
	// STEP 1: Create the schema_migrations table to track which migrations have been applied
	// This table stores the version (filename) and timestamp of each applied migration
	if err := db.createSchemaMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// STEP 2: Determine migration directory based on database driver
	// Different databases (SQLite vs PostgreSQL) may have different SQL syntax
	var migrationDir string
	switch db.driver {
	case "sqlite3":
		migrationDir = "migrations/sqlite"
	case "postgres":
		migrationDir = "migrations/postgres"
	default:
		return fmt.Errorf("unsupported driver: %s", db.driver)
	}

	// STEP 3: Read all migration files from the filesystem
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	// STEP 4: Filter and sort SQL files to ensure they run in the correct order
	// Migration files are named with numeric prefixes (001_, 002_, etc.) to control execution order
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// STEP 5: Execute each migration file (but only if it hasn't been applied yet)
	for _, filename := range sqlFiles {
		// Check if this migration has already been applied
		// This prevents re-running migrations that use ALTER TABLE or other non-idempotent operations
		applied, err := db.isMigrationApplied(filename)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", filename, err)
		}

		// Skip migrations that have already been applied
		if applied {
			fmt.Printf("⊘ Skipped (already applied): %s\n", filename)
			continue
		}

		// Read the migration file content
		filePath := filepath.Join(migrationDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Execute the SQL migration
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record that this migration has been successfully applied
		// This ensures it won't run again on the next application start
		if err := db.recordMigration(filename); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		fmt.Printf("✓ Applied migration: %s\n", filename)
	}

	return nil
}

// createSchemaMigrationsTable creates the table used to track which migrations have been applied.
// This table has two columns:
// - version: The migration filename (e.g., "005_add_review_metrics.sql")
// - applied_at: Timestamp when the migration was executed
//
// The table is created with "IF NOT EXISTS" so it's safe to call on every startup.
func (db *DB) createSchemaMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

// isMigrationApplied checks if a specific migration has already been applied.
// It queries the schema_migrations table to see if the migration version exists.
//
// Parameters:
//   - version: The migration filename (e.g., "005_add_review_metrics.sql")
//
// Returns:
//   - true if the migration has been applied
//   - false if the migration has not been applied
//   - error if the database query fails
func (db *DB) isMigrationApplied(version string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM schema_migrations WHERE version = ?"
	if err := db.Get(&count, query, version); err != nil {
		return false, err
	}
	return count > 0, nil
}

// recordMigration records that a migration has been successfully applied.
// This inserts a new row into the schema_migrations table with the migration version.
//
// Parameters:
//   - version: The migration filename (e.g., "005_add_review_metrics.sql")
//
// Returns:
//   - error if the insert fails (e.g., if the migration was already recorded)
func (db *DB) recordMigration(version string) error {
	query := "INSERT INTO schema_migrations (version) VALUES (?)"
	_, err := db.Exec(query, version)
	return err
}
