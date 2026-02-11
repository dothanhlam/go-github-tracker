-- Create pr_metrics table (PostgreSQL version)
CREATE TABLE IF NOT EXISTS pr_metrics (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    pr_number INTEGER NOT NULL,
    repository VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    title TEXT,
    created_at TIMESTAMP NOT NULL,
    merged_at TIMESTAMP,
    closed_at TIMESTAMP,
    cycle_time_hours INTEGER,
    state VARCHAR(50) CHECK(state IN ('open', 'closed', 'merged')),
    created_date DATE GENERATED ALWAYS AS (DATE(created_at)) STORED,
    UNIQUE(team_id, repository, pr_number)
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_pr_metrics_team_id ON pr_metrics(team_id);
CREATE INDEX IF NOT EXISTS idx_pr_metrics_created_date ON pr_metrics(created_date);
CREATE INDEX IF NOT EXISTS idx_pr_metrics_merged_at ON pr_metrics(merged_at);
CREATE INDEX IF NOT EXISTS idx_pr_metrics_author ON pr_metrics(author);
