package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dothanhlam/go-github-tracker/internal/api/response"
	"github.com/dothanhlam/go-github-tracker/internal/service"
	"github.com/go-chi/chi/v5"
)

// TeamsHandler handles team-related requests
type TeamsHandler struct {
	metricsService *service.MetricsService
}

// NewTeamsHandler creates a new teams handler
func NewTeamsHandler(metricsService *service.MetricsService) *TeamsHandler {
	return &TeamsHandler{
		metricsService: metricsService,
	}
}

// ListTeams handles GET /api/v1/teams
func (h *TeamsHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.metricsService.ListTeams()
	if err != nil {
		response.InternalError(w, "Failed to fetch teams")
		return
	}

	resp := map[string]interface{}{
		"teams": teams,
	}
	response.JSON(w, http.StatusOK, resp)
}

// GetVelocity handles GET /api/v1/teams/{id}/velocity
func (h *TeamsHandler) GetVelocity(w http.ResponseWriter, r *http.Request) {
	teamID, err := h.getTeamID(r)
	if err != nil {
		response.BadRequest(w, "Invalid team ID")
		return
	}

	startDate, endDate, granularity := h.parseDateParams(r)

	metrics, err := h.metricsService.GetTeamVelocity(teamID, startDate, endDate, granularity)
	if err != nil {
		if err.Error() == "team not found" {
			response.NotFound(w, "Team not found")
			return
		}
		response.InternalError(w, "Failed to fetch velocity metrics")
		return
	}

	response.JSON(w, http.StatusOK, metrics)
}

// GetLeadTime handles GET /api/v1/teams/{id}/lead-time
func (h *TeamsHandler) GetLeadTime(w http.ResponseWriter, r *http.Request) {
	teamID, err := h.getTeamID(r)
	if err != nil {
		response.BadRequest(w, "Invalid team ID")
		return
	}

	startDate, endDate, granularity := h.parseDateParams(r)

	metrics, err := h.metricsService.GetTeamLeadTime(teamID, startDate, endDate, granularity)
	if err != nil {
		if err.Error() == "team not found" {
			response.NotFound(w, "Team not found")
			return
		}
		response.InternalError(w, "Failed to fetch lead time metrics")
		return
	}

	response.JSON(w, http.StatusOK, metrics)
}

// GetReviewTurnaround handles GET /api/v1/teams/{id}/review-turnaround
func (h *TeamsHandler) GetReviewTurnaround(w http.ResponseWriter, r *http.Request) {
	teamID, err := h.getTeamID(r)
	if err != nil {
		response.BadRequest(w, "Invalid team ID")
		return
	}

	startDate, endDate, _ := h.parseDateParams(r)

	metrics, err := h.metricsService.GetReviewTurnaround(teamID, startDate, endDate)
	if err != nil {
		if err.Error() == "team not found" {
			response.NotFound(w, "Team not found")
			return
		}
		response.InternalError(w, "Failed to fetch review turnaround metrics")
		return
	}

	response.JSON(w, http.StatusOK, metrics)
}

// GetReviewEngagement handles GET /api/v1/teams/{id}/review-engagement
func (h *TeamsHandler) GetReviewEngagement(w http.ResponseWriter, r *http.Request) {
	teamID, err := h.getTeamID(r)
	if err != nil {
		response.BadRequest(w, "Invalid team ID")
		return
	}

	startDate, endDate, _ := h.parseDateParams(r)

	metrics, err := h.metricsService.GetReviewEngagement(teamID, startDate, endDate)
	if err != nil {
		if err.Error() == "team not found" {
			response.NotFound(w, "Team not found")
			return
		}
		response.InternalError(w, "Failed to fetch review engagement metrics")
		return
	}

	response.JSON(w, http.StatusOK, metrics)
}

// GetKnowledgeSharing handles GET /api/v1/teams/{id}/knowledge-sharing
func (h *TeamsHandler) GetKnowledgeSharing(w http.ResponseWriter, r *http.Request) {
	teamID, err := h.getTeamID(r)
	if err != nil {
		response.BadRequest(w, "Invalid team ID")
		return
	}

	startDate, endDate, _ := h.parseDateParams(r)

	metrics, err := h.metricsService.GetKnowledgeSharing(teamID, startDate, endDate)
	if err != nil {
		if err.Error() == "team not found" {
			response.NotFound(w, "Team not found")
			return
		}
		response.InternalError(w, "Failed to fetch knowledge sharing metrics")
		return
	}

	response.JSON(w, http.StatusOK, metrics)
}

// Helper functions

func (h *TeamsHandler) getTeamID(r *http.Request) (int, error) {
	teamIDStr := chi.URLParam(r, "id")
	return strconv.Atoi(teamIDStr)
}

func (h *TeamsHandler) parseDateParams(r *http.Request) (startDate, endDate time.Time, granularity string) {
	// Default to last 30 days
	endDate = time.Now()
	startDate = endDate.AddDate(0, 0, -30)
	granularity = "week"

	// Parse start_date
	if startStr := r.URL.Query().Get("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = t
		}
	}

	// Parse end_date
	if endStr := r.URL.Query().Get("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = t
		}
	}

	// Parse granularity
	if gran := r.URL.Query().Get("granularity"); gran != "" {
		if gran == "day" || gran == "week" || gran == "month" {
			granularity = gran
		}
	}

	return
}
