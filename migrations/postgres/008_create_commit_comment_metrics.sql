-- Create commit_metrics table
CREATE TABLE IF NOT EXISTS commit_metrics (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    repository VARCHAR(255) NOT NULL,
    commit_hash VARCHAR(64) NOT NULL,
    author VARCHAR(255) NOT NULL,
    message TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_date DATE NOT NULL,
    UNIQUE(team_id, repository, commit_hash)
);

-- Create indexes for commit_metrics
CREATE INDEX IF NOT EXISTS idx_commit_metrics_team_id ON commit_metrics(team_id);
CREATE INDEX IF NOT EXISTS idx_commit_metrics_created_date ON commit_metrics(created_date);
CREATE INDEX IF NOT EXISTS idx_commit_metrics_author ON commit_metrics(author);

-- Create comment_metrics table
CREATE TABLE IF NOT EXISTS comment_metrics (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    repository VARCHAR(255) NOT NULL,
    comment_id BIGINT NOT NULL,
    author VARCHAR(255) NOT NULL,
    body TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_date DATE NOT NULL,
    comment_type VARCHAR(50) NOT NULL, -- 'commit', 'issue', 'pull_request_review', etc.
    UNIQUE(team_id, repository, comment_id, comment_type)
);

-- Create indexes for comment_metrics
CREATE INDEX IF NOT EXISTS idx_comment_metrics_team_id ON comment_metrics(team_id);
CREATE INDEX IF NOT EXISTS idx_comment_metrics_created_date ON comment_metrics(created_date);
CREATE INDEX IF NOT EXISTS idx_comment_metrics_author ON comment_metrics(author);

-- Create views for aggregations

-- Commit Velocity View (Weekly)
CREATE OR REPLACE VIEW view_team_commit_velocity AS
SELECT 
    team_id,
    TO_CHAR(created_date, 'IYYY-IW') as week,
    COUNT(id) as commits_count
FROM commit_metrics
GROUP BY team_id, TO_CHAR(created_date, 'IYYY-IW');

-- Comment Activity View (Weekly)
CREATE OR REPLACE VIEW view_team_comment_activity AS
SELECT 
    team_id,
    TO_CHAR(created_date, 'IYYY-IW') as week,
    COUNT(id) as comments_count,
    comment_type
FROM comment_metrics
GROUP BY team_id, TO_CHAR(created_date, 'IYYY-IW'), comment_type;

-- Member Commit Velocity View (Weekly)
CREATE OR REPLACE VIEW view_member_commit_velocity AS
SELECT 
    team_id,
    author as github_username,
    TO_CHAR(created_date, 'IYYY-IW') as week,
    COUNT(id) as commits_count
FROM commit_metrics
GROUP BY team_id, author, TO_CHAR(created_date, 'IYYY-IW');

-- Member Comment Activity View (Weekly)
CREATE OR REPLACE VIEW view_member_comment_activity AS
SELECT 
    team_id,
    author as github_username,
    TO_CHAR(created_date, 'IYYY-IW') as week,
    COUNT(id) as comments_count,
    comment_type
FROM comment_metrics
GROUP BY team_id, author, TO_CHAR(created_date, 'IYYY-IW'), comment_type;
