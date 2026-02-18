package handlers

import (
	"net/http"

	"github.com/dothanhlam/go-github-tracker/internal/api/response"
	"github.com/dothanhlam/go-github-tracker/internal/database"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *database.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
}

// Handle processes health check requests
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	dbStatus := "connected"
	if err := h.db.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	resp := HealthResponse{
		Status:   "healthy",
		Database: dbStatus,
		Version:  "1.0.0",
	}

	response.JSON(w, http.StatusOK, resp)
}
