package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dothanhlam/go-github-tracker/internal/config"
	"github.com/dothanhlam/go-github-tracker/internal/database"
	"github.com/dothanhlam/go-github-tracker/internal/mcp"
	mcpsdk "github.com/mark3labs/mcp-go/server"
)

func main() {
	// Support loading different .env files via ENV_FILE environment variable
	// This allows MCP config to specify which environment to use
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env.local" // Default to .env.local
	}
	
	// Set the environment file for config.Load() to use
	os.Setenv("ENV_FILE", envFile)
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.DBDriver, cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create MCP server
	server := mcpsdk.NewMCPServer(
		"dora-metrics-mcp",
		"1.0.0",
		mcpsdk.WithResourceCapabilities(true, false), // Read-only resources
		mcpsdk.WithToolCapabilities(true),            // Enable tools
	)

	// Initialize metrics handler
	handler := mcp.NewMetricsHandler(db)

	// Register resources
	handler.RegisterResources(server)

	// Register tools
	handler.RegisterTools(server)

	// Start server with stdio transport
	fmt.Fprintln(os.Stderr, "🚀 DORA Metrics MCP Server starting...")
	fmt.Fprintln(os.Stderr, "📊 Connected to database:", cfg.DBURL)
	fmt.Fprintln(os.Stderr, "✓ Resources and tools registered")

	if err := server.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
