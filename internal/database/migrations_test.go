package database

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestSchemaMigrationsTable tests that the schema_migrations table is created
func TestSchemaMigrationsTable(t *testing.T) {
	// Create a temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Connect("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Create schema_migrations table
	err = db.createSchemaMigrationsTable()
	if err != nil {
		t.Fatalf("failed to create schema_migrations table: %v", err)
	}

	// Verify table exists
	var tableName string
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name='schema_migrations'"
	err = db.Get(&tableName, query)
	if err != nil {
		t.Fatalf("schema_migrations table not found: %v", err)
	}

	if tableName != "schema_migrations" {
		t.Errorf("table name = %v, want 'schema_migrations'", tableName)
	}
}

// TestMigrationTracking tests the migration tracking functionality
func TestMigrationTracking(t *testing.T) {
	// Create a temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Connect("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Create schema_migrations table
	err = db.createSchemaMigrationsTable()
	if err != nil {
		t.Fatalf("failed to create schema_migrations table: %v", err)
	}

	// Test 1: Migration should not be applied initially
	applied, err := db.isMigrationApplied("001_test_migration.sql")
	if err != nil {
		t.Fatalf("failed to check migration status: %v", err)
	}
	if applied {
		t.Error("migration should not be applied initially")
	}

	// Test 2: Record migration
	err = db.recordMigration("001_test_migration.sql")
	if err != nil {
		t.Fatalf("failed to record migration: %v", err)
	}

	// Test 3: Migration should now be applied
	applied, err = db.isMigrationApplied("001_test_migration.sql")
	if err != nil {
		t.Fatalf("failed to check migration status: %v", err)
	}
	if !applied {
		t.Error("migration should be applied after recording")
	}

	// Test 4: Verify migration is in database
	var count int
	query := "SELECT COUNT(*) FROM schema_migrations WHERE version = ?"
	err = db.Get(&count, query, "001_test_migration.sql")
	if err != nil {
		t.Fatalf("failed to query schema_migrations: %v", err)
	}
	if count != 1 {
		t.Errorf("migration count = %v, want 1", count)
	}
}

// TestMultipleMigrations tests tracking multiple migrations
func TestMultipleMigrations(t *testing.T) {
	// Create a temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Connect("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Create schema_migrations table
	err = db.createSchemaMigrationsTable()
	if err != nil {
		t.Fatalf("failed to create schema_migrations table: %v", err)
	}

	// Record multiple migrations
	migrations := []string{
		"001_create_teams.sql",
		"002_create_team_memberships.sql",
		"003_create_pr_metrics.sql",
	}

	for _, migration := range migrations {
		err = db.recordMigration(migration)
		if err != nil {
			t.Fatalf("failed to record migration %s: %v", migration, err)
		}
	}

	// Verify all migrations are recorded
	for _, migration := range migrations {
		applied, err := db.isMigrationApplied(migration)
		if err != nil {
			t.Fatalf("failed to check migration %s: %v", migration, err)
		}
		if !applied {
			t.Errorf("migration %s should be applied", migration)
		}
	}

	// Verify total count
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM schema_migrations")
	if err != nil {
		t.Fatalf("failed to count migrations: %v", err)
	}
	if count != len(migrations) {
		t.Errorf("migration count = %v, want %v", count, len(migrations))
	}
}

// TestDatabaseConnection tests database connection with different drivers
func TestDatabaseConnection(t *testing.T) {
	tests := []struct {
		name    string
		driver  string
		dsn     string
		wantErr bool
	}{
		{
			name:    "valid sqlite3 in-memory",
			driver:  "sqlite3",
			dsn:     ":memory:",
			wantErr: false,
		},
		{
			name:    "valid sqlite3 file",
			driver:  "sqlite3",
			dsn:     filepath.Join(t.TempDir(), "test.db"),
			wantErr: false,
		},
		{
			name:    "invalid driver",
			driver:  "invalid",
			dsn:     ":memory:",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.driver, tt.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if db != nil {
				db.Close()
			}
		})
	}
}

// TestMigrationIdempotency tests that creating schema_migrations table is idempotent
func TestMigrationIdempotency(t *testing.T) {
	// Create a temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Connect("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Create table multiple times - should not error
	for i := 0; i < 3; i++ {
		err = db.createSchemaMigrationsTable()
		if err != nil {
			t.Fatalf("iteration %d: failed to create schema_migrations table: %v", i, err)
		}
	}

	// Verify table still exists and is correct
	var tableName string
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name='schema_migrations'"
	err = db.Get(&tableName, query)
	if err != nil {
		t.Fatalf("schema_migrations table not found: %v", err)
	}
}

// TestRunMigrationsWithNoFiles tests running migrations when no migration files exist
func TestRunMigrationsWithNoFiles(t *testing.T) {
	// Create a temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Connect("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Create empty migration directory
	migrationDir := filepath.Join(tmpDir, "migrations", "sqlite")
	err = os.MkdirAll(migrationDir, 0755)
	if err != nil {
		t.Fatalf("failed to create migration directory: %v", err)
	}

	// Change to temp directory so migrations can be found
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Run migrations - should succeed with no files
	err = db.RunMigrations()
	if err != nil {
		t.Errorf("RunMigrations() with no files should not error: %v", err)
	}

	// Verify schema_migrations table was created
	var tableName string
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name='schema_migrations'"
	err = db.Get(&tableName, query)
	if err != nil {
		t.Fatalf("schema_migrations table should be created: %v", err)
	}
}
