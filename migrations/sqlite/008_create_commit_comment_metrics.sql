-- Create commit_metrics table
CREATE TABLE IF NOT EXISTS commit_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    team_id INTEGER NOT NULL,
    repository TEXT NOT NULL,
    commit_hash TEXT NOT NULL,
    author TEXT NOT NULL,
    message TEXT,
    created_at DATETIME NOT NULL,
    created_date DATE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    UNIQUE(team_id, repository, commit_hash)
);

-- Create indexes for commit_metrics
CREATE INDEX IF NOT EXISTS idx_commit_metrics_team_id ON commit_metrics(team_id);
CREATE INDEX IF NOT EXISTS idx_commit_metrics_created_date ON commit_metrics(created_date);
CREATE INDEX IF NOT EXISTS idx_commit_metrics_author ON commit_metrics(author);

-- Create comment_metrics table
CREATE TABLE IF NOT EXISTS comment_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    team_id INTEGER NOT NULL,
    repository TEXT NOT NULL,
    comment_id INTEGER NOT NULL,
    author TEXT NOT NULL,
    body TEXT,
    created_at DATETIME NOT NULL,
    created_date DATE,
    comment_type TEXT NOT NULL, -- 'commit', 'issue', 'pull_request_review', etc.
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    UNIQUE(team_id, repository, comment_id, comment_type)
);

-- Create indexes for comment_metrics
CREATE INDEX IF NOT EXISTS idx_comment_metrics_team_id ON comment_metrics(team_id);
CREATE INDEX IF NOT EXISTS idx_comment_metrics_created_date ON comment_metrics(created_date);
CREATE INDEX IF NOT EXISTS idx_comment_metrics_author ON comment_metrics(author);

-- Create views for aggregations

-- Commit Velocity View (Weekly)
CREATE VIEW IF NOT EXISTS view_team_commit_velocity AS
SELECT 
    team_id,
    strftime('%Y-%W', created_date) as week,
    COUNT(id) as commits_count
FROM commit_metrics
GROUP BY team_id, week;

-- Comment Activity View (Weekly)
CREATE VIEW IF NOT EXISTS view_team_comment_activity AS
SELECT 
    team_id,
    strftime('%Y-%W', created_date) as week,
    COUNT(id) as comments_count,
    comment_type
FROM comment_metrics
GROUP BY team_id, week, comment_type;

-- Member Commit Velocity View (Weekly)
CREATE VIEW IF NOT EXISTS view_member_commit_velocity AS
SELECT 
    team_id,
    author as github_username,
    strftime('%Y-%W', created_date) as week,
    COUNT(id) as commits_count
FROM commit_metrics
GROUP BY team_id, author, week;

-- Member Comment Activity View (Weekly)
CREATE VIEW IF NOT EXISTS view_member_comment_activity AS
SELECT 
    team_id,
    author as github_username,
    strftime('%Y-%W', created_date) as week,
    COUNT(id) as comments_count,
    comment_type
FROM comment_metrics
GROUP BY team_id, author, week, comment_type;
