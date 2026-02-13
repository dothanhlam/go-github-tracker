package config

import (
	"encoding/json"
	"testing"
)

// TestTeamConfigJSONParsing tests parsing team configuration from JSON
func TestTeamConfigJSONParsing(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		wantLen int
	}{
		{
			name:    "empty array",
			json:    "[]",
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "single team with members",
			json: `[{
				"name": "Test Team",
				"members": [
					{"username": "user1", "allocation": 1.0},
					{"username": "user2", "allocation": 0.5}
				]
			}]`,
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "multiple teams",
			json: `[
				{
					"name": "Team A",
					"members": [{"username": "user1", "allocation": 1.0}]
				},
				{
					"name": "Team B",
					"members": [{"username": "user2", "allocation": 1.0}]
				}
			]`,
			wantErr: false,
			wantLen: 2,
		},
		{
			name:    "invalid JSON",
			json:    "{invalid json}",
			wantErr: true,
			wantLen: 0,
		},
		{
			name: "team without team_id (auto-generated)",
			json: `[{
				"name": "Auto ID Team",
				"members": [{"username": "user1", "allocation": 1.0}]
			}]`,
			wantErr: false,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var teams []TeamConfig
			err := json.Unmarshal([]byte(tt.json), &teams)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(teams) != tt.wantLen {
				t.Errorf("parsed %d teams, want %d", len(teams), tt.wantLen)
			}
		})
	}
}

// TestTeamConfigStructure tests the team configuration structure
func TestTeamConfigStructure(t *testing.T) {
	jsonStr := `{
		"name": "Data Platform Team",
		"members": [
			{"username": "alice", "allocation": 1.0},
			{"username": "bob", "allocation": 0.5}
		]
	}`

	var team TeamConfig
	err := json.Unmarshal([]byte(jsonStr), &team)
	if err != nil {
		t.Fatalf("failed to parse team config: %v", err)
	}

	// Verify team name
	if team.Name != "Data Platform Team" {
		t.Errorf("team.Name = %v, want 'Data Platform Team'", team.Name)
	}

	// Verify members count
	if len(team.Members) != 2 {
		t.Errorf("len(team.Members) = %v, want 2", len(team.Members))
	}

	// Verify first member
	if team.Members[0].Username != "alice" {
		t.Errorf("team.Members[0].Username = %v, want 'alice'", team.Members[0].Username)
	}
	if team.Members[0].Allocation != 1.0 {
		t.Errorf("team.Members[0].Allocation = %v, want 1.0", team.Members[0].Allocation)
	}

	// Verify second member
	if team.Members[1].Username != "bob" {
		t.Errorf("team.Members[1].Username = %v, want 'bob'", team.Members[1].Username)
	}
	if team.Members[1].Allocation != 0.5 {
		t.Errorf("team.Members[1].Allocation = %v, want 0.5", team.Members[1].Allocation)
	}
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid sqlite config",
			config: &Config{
				DBDriver: "sqlite3",
				DBURL:    "./data/test.db",
			},
			wantErr: false,
		},
		{
			name: "valid postgres config",
			config: &Config{
				DBDriver: "postgres",
				DBURL:    "postgresql://localhost/test",
			},
			wantErr: false,
		},
		{
			name: "invalid driver",
			config: &Config{
				DBDriver: "mysql",
				DBURL:    "localhost:3306",
			},
			wantErr: true,
		},
		{
			name: "missing DB URL",
			config: &Config{
				DBDriver: "sqlite3",
				DBURL:    "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTeamMemberAllocation tests various allocation values
func TestTeamMemberAllocation(t *testing.T) {
	tests := []struct {
		name       string
		allocation float64
	}{
		{"full time", 1.0},
		{"half time", 0.5},
		{"quarter time", 0.25},
		{"80 percent", 0.8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			member := TeamMemberConfig{
				Username:   "testuser",
				Allocation: tt.allocation,
			}

			if member.Allocation != tt.allocation {
				t.Errorf("member.Allocation = %v, want %v", member.Allocation, tt.allocation)
			}
		})
	}
}
