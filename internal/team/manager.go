package team

import (
	"fmt"
	"time"

	"github.com/dothanhlam/go-github-tracker/internal/config"
	"github.com/dothanhlam/go-github-tracker/internal/database"
)

// Manager handles team membership lookups
type Manager struct {
	db            *database.DB
	teams         map[int]*database.Team
	membershipMap map[string][]int // username -> team IDs
}

// NewManager creates a new team manager
func NewManager(db *database.DB, cfg *config.Config) (*Manager, error) {
	m := &Manager{
		db:            db,
		teams:         make(map[int]*database.Team),
		membershipMap: make(map[string][]int),
	}

	// Sync teams from config to database
	if err := m.syncTeams(cfg); err != nil {
		return nil, fmt.Errorf("failed to sync teams: %w", err)
	}

	// Load teams into memory
	if err := m.loadTeams(); err != nil {
		return nil, fmt.Errorf("failed to load teams: %w", err)
	}

	return m, nil
}

// syncTeams syncs teams from config to database
func (m *Manager) syncTeams(cfg *config.Config) error {
	for _, teamCfg := range cfg.Teams {
		// Upsert team
		team := &database.Team{
			ID:   teamCfg.TeamID,
			Name: teamCfg.Name,
		}

		query := `
			INSERT INTO teams (id, name, created_at, updated_at)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				name = excluded.name,
				updated_at = excluded.updated_at
		`
		now := time.Now()
		if _, err := m.db.Exec(query, team.ID, team.Name, now, now); err != nil {
			return fmt.Errorf("failed to upsert team %d: %w", team.ID, err)
		}

		// Upsert team memberships
		for _, member := range teamCfg.Members {
			query := `
				INSERT INTO team_memberships (team_id, github_username, allocation_weight, joined_at, created_at)
				VALUES (?, ?, ?, ?, ?)
				ON CONFLICT(team_id, github_username, joined_at) DO UPDATE SET
					allocation_weight = excluded.allocation_weight
			`
			if _, err := m.db.Exec(query, team.ID, member.Username, member.Allocation, now, now); err != nil {
				return fmt.Errorf("failed to upsert membership for %s: %w", member.Username, err)
			}
		}
	}

	return nil
}

// loadTeams loads teams and memberships into memory
func (m *Manager) loadTeams() error {
	// Load teams
	var teams []database.Team
	if err := m.db.Select(&teams, "SELECT * FROM teams"); err != nil {
		return fmt.Errorf("failed to load teams: %w", err)
	}

	for i := range teams {
		m.teams[teams[i].ID] = &teams[i]
	}

	// Load memberships
	var memberships []database.TeamMembership
	query := "SELECT * FROM team_memberships WHERE left_at IS NULL"
	if err := m.db.Select(&memberships, query); err != nil {
		return fmt.Errorf("failed to load memberships: %w", err)
	}

	for _, membership := range memberships {
		m.membershipMap[membership.GitHubUsername] = append(
			m.membershipMap[membership.GitHubUsername],
			membership.TeamID,
		)
	}

	return nil
}

// IsMember checks if a username is a member of any team
func (m *Manager) IsMember(username string) bool {
	_, exists := m.membershipMap[username]
	return exists
}

// GetTeamsForUser returns all team IDs for a username
func (m *Manager) GetTeamsForUser(username string) []int {
	return m.membershipMap[username]
}

// GetAllTeamIDs returns all team IDs
func (m *Manager) GetAllTeamIDs() []int {
	var ids []int
	for id := range m.teams {
		ids = append(ids, id)
	}
	return ids
}

// IsExternalReviewer checks if a reviewer is external to a team
func (m *Manager) IsExternalReviewer(reviewer string, teamID int) bool {
	teams := m.GetTeamsForUser(reviewer)
	for _, tid := range teams {
		if tid == teamID {
			return false
		}
	}
	return true
}
