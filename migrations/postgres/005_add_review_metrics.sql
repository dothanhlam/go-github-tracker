-- Add review metrics columns to pr_metrics table (PostgreSQL version)
ALTER TABLE pr_metrics ADD COLUMN first_review_at TIMESTAMP;
ALTER TABLE pr_metrics ADD COLUMN review_turnaround_hours INTEGER;
ALTER TABLE pr_metrics ADD COLUMN review_comments_count INTEGER DEFAULT 0;
ALTER TABLE pr_metrics ADD COLUMN conversation_count INTEGER DEFAULT 0;
ALTER TABLE pr_metrics ADD COLUMN changes_requested_count INTEGER DEFAULT 0;
ALTER TABLE pr_metrics ADD COLUMN approved_count INTEGER DEFAULT 0;
ALTER TABLE pr_metrics ADD COLUMN reviewers_count INTEGER DEFAULT 0;
ALTER TABLE pr_metrics ADD COLUMN external_reviewers_count INTEGER DEFAULT 0;
ALTER TABLE pr_metrics ADD COLUMN reviewers_list JSONB; -- JSON array of reviewer usernames

-- Create indexes for new fields
CREATE INDEX IF NOT EXISTS idx_pr_metrics_first_review_at ON pr_metrics(first_review_at);
CREATE INDEX IF NOT EXISTS idx_pr_metrics_reviewers_count ON pr_metrics(reviewers_count);
CREATE INDEX IF NOT EXISTS idx_pr_metrics_review_turnaround ON pr_metrics(review_turnaround_hours);

-- Create GIN index for JSONB column (PostgreSQL-specific)
CREATE INDEX IF NOT EXISTS idx_pr_metrics_reviewers_list ON pr_metrics USING GIN(reviewers_list);
