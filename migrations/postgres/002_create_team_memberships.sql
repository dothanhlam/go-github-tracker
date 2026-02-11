-- Create team_memberships table (PostgreSQL version)
CREATE TABLE IF NOT EXISTS team_memberships (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    github_username VARCHAR(255) NOT NULL,
    allocation_weight DECIMAL(3,2) DEFAULT 1.00 CHECK(allocation_weight >= 0.00 AND allocation_weight <= 1.00),
    joined_at TIMESTAMP NOT NULL,
    left_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(team_id, github_username, joined_at)
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_team_memberships_username ON team_memberships(github_username);
CREATE INDEX IF NOT EXISTS idx_team_memberships_team_id ON team_memberships(team_id);
