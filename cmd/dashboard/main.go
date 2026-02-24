package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dothanhlam/go-github-tracker/internal/config"
	"github.com/dothanhlam/go-github-tracker/internal/database"
	"github.com/dothanhlam/go-github-tracker/internal/service"
	"github.com/dothanhlam/go-github-tracker/internal/tui"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Silence standard logger to not mess with TUI
	log.SetOutput(os.Stderr)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// For local TUI, ensure we are using sqlite3
	if cfg.DBDriver != "sqlite3" {
		fmt.Println("Warning: TUI is designed for local sqlite3 database, but config is using:", cfg.DBDriver)
	}

	// Connect to database
	db, err := database.Connect(cfg.DBDriver, cfg.DBURL)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	// Initialize metrics service
	metricsService := service.NewMetricsService(db)

	// Initialize TUI model
	m := tui.InitialModel(metricsService)

	// Start Bubbletea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI dashboard: %v\n", err)
		os.Exit(1)
	}
}
