-- Create team_memberships table
CREATE TABLE IF NOT EXISTS team_memberships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    team_id INTEGER NOT NULL,
    github_username TEXT NOT NULL,
    allocation_weight REAL DEFAULT 1.00 CHECK(allocation_weight >= 0.00 AND allocation_weight <= 1.00),
    joined_at DATETIME NOT NULL,
    left_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    UNIQUE(team_id, github_username, joined_at)
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_team_memberships_username ON team_memberships(github_username);
CREATE INDEX IF NOT EXISTS idx_team_memberships_team_id ON team_memberships(team_id);
