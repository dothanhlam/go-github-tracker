package api

import (
	"net/http"
	"strings"

	"github.com/dothanhlam/go-github-tracker/internal/api/handlers"
	"github.com/dothanhlam/go-github-tracker/internal/api/middleware"
	"github.com/dothanhlam/go-github-tracker/internal/database"
	"github.com/dothanhlam/go-github-tracker/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server represents the API server
type Server struct {
	router         *chi.Mux
	db             *database.DB
	apiKeys        []string
	metricsService *service.MetricsService
}

// NewServer creates a new API server
func NewServer(db *database.DB, apiKeys string) *Server {
	s := &Server{
		router:         chi.NewRouter(),
		db:             db,
		apiKeys:        parseAPIKeys(apiKeys),
		metricsService: service.NewMetricsService(db),
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// Basic middleware
	s.router.Use(chimiddleware.RequestID)
	s.router.Use(chimiddleware.RealIP)
	s.router.Use(chimiddleware.Recoverer)
	s.router.Use(middleware.Logger)

	// CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Configure this for production
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API Key authentication
	s.router.Use(middleware.APIKeyAuth(s.apiKeys))
}

// setupRoutes configures routes
func (s *Server) setupRoutes() {
	// Health check (no auth required)
	healthHandler := handlers.NewHealthHandler(s.db)
	s.router.Get("/api/v1/health", healthHandler.Handle)

	// Teams endpoints
	teamsHandler := handlers.NewTeamsHandler(s.metricsService)
	s.router.Route("/api/v1/teams", func(r chi.Router) {
		r.Get("/", teamsHandler.ListTeams)
		r.Get("/{id}/velocity", teamsHandler.GetVelocity)
		r.Get("/{id}/lead-time", teamsHandler.GetLeadTime)
		r.Get("/{id}/review-turnaround", teamsHandler.GetReviewTurnaround)
		r.Get("/{id}/review-engagement", teamsHandler.GetReviewEngagement)
		r.Get("/{id}/knowledge-sharing", teamsHandler.GetKnowledgeSharing)
	})
}

// Handler returns the HTTP handler
func (s *Server) Handler() http.Handler {
	return s.router
}

// parseAPIKeys splits comma-separated API keys
func parseAPIKeys(keys string) []string {
	if keys == "" {
		return []string{}
	}

	parts := strings.Split(keys, ",")
	result := make([]string, 0, len(parts))
	for _, key := range parts {
		if trimmed := strings.TrimSpace(key); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
