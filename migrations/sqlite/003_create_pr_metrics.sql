-- Create pr_metrics table
CREATE TABLE IF NOT EXISTS pr_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    team_id INTEGER NOT NULL,
    pr_number INTEGER NOT NULL,
    repository TEXT NOT NULL,
    author TEXT NOT NULL,
    title TEXT,
    created_at DATETIME NOT NULL,
    merged_at DATETIME,
    closed_at DATETIME,
    cycle_time_hours INTEGER,
    state TEXT CHECK(state IN ('open', 'closed', 'merged')),
    created_date DATE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    UNIQUE(team_id, repository, pr_number)
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_pr_metrics_team_id ON pr_metrics(team_id);
CREATE INDEX IF NOT EXISTS idx_pr_metrics_created_date ON pr_metrics(created_date);
CREATE INDEX IF NOT EXISTS idx_pr_metrics_merged_at ON pr_metrics(merged_at);
CREATE INDEX IF NOT EXISTS idx_pr_metrics_author ON pr_metrics(author);
