package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dothanhlam/go-github-tracker/internal/collector"
	"github.com/dothanhlam/go-github-tracker/internal/config"
	"github.com/dothanhlam/go-github-tracker/internal/database"
)

func main() {
	fmt.Println("🚀 DORA Metrics Collector - Phase 1 MVP")
	fmt.Println("========================================")

	// Load configuration
	fmt.Println("\n📋 Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	fmt.Printf("✓ Configuration loaded\n")
	fmt.Printf("  - Database Driver: %s\n", cfg.DBDriver)
	fmt.Printf("  - Database URL: %s\n", cfg.DBURL)
	fmt.Printf("  - Teams configured: %d\n", len(cfg.Teams))
	fmt.Printf("  - Repositories: %d\n", len(cfg.Repositories))

	// Ensure data directory exists for SQLite
	if cfg.DBDriver == "sqlite3" {
		if err := os.MkdirAll("./data", 0755); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}
	}

	// Connect to database
	fmt.Println("\n🔌 Connecting to database...")
	db, err := database.Connect(cfg.DBDriver, cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	fmt.Println("✓ Database connected")

	// Run migrations
	fmt.Println("\n📦 Running database migrations...")
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	fmt.Println("✓ All migrations applied")

	// Verify tables exist
	fmt.Println("\n🔍 Verifying database schema...")
	if err := verifySchema(db); err != nil {
		log.Fatalf("Schema verification failed: %v", err)
	}
	fmt.Println("✓ Schema verified")

	// Run Phase 2: Data Collection
	fmt.Println("\n🔄 Starting PR data collection...")
	if err := runCollector(cfg, db); err != nil {
		log.Fatalf("Collection failed: %v", err)
	}

	fmt.Println("\n✅ Collection complete!")
	fmt.Println("\nQuery your metrics:")
	fmt.Println("  sqlite3 ./data/dora_metrics.db \"SELECT * FROM view_team_velocity;\"")
	fmt.Println("  sqlite3 ./data/dora_metrics.db \"SELECT * FROM view_review_turnaround;\"")
}

// runCollector executes the PR collection process
func runCollector(cfg *config.Config, db *database.DB) error {
	// Import collector package
	collector, err := collector.New(cfg, db)
	if err != nil {
		return fmt.Errorf("failed to create collector: %w", err)
	}

	return collector.Run()
}

// verifySchema checks that all required tables exist
func verifySchema(db *database.DB) error {
	tables := []string{"teams", "team_memberships", "pr_metrics"}
	
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		if err := db.Get(&count, query); err != nil {
			return fmt.Errorf("table %s not found or inaccessible: %w", table, err)
		}
		fmt.Printf("  ✓ Table '%s' exists (rows: %d)\n", table, count)
	}

	return nil
}
