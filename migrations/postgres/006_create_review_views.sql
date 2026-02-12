-- Create view for review turnaround metrics (PostgreSQL version)
CREATE OR REPLACE VIEW view_review_turnaround AS
SELECT 
    team_id,
    DATE_TRUNC('month', created_at) as month,
    AVG(review_turnaround_hours) as avg_turnaround_hours,
    MIN(review_turnaround_hours) as min_turnaround_hours,
    MAX(review_turnaround_hours) as max_turnaround_hours,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY review_turnaround_hours) as median_turnaround_hours,
    PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY review_turnaround_hours) as p95_turnaround_hours,
    COUNT(*) as pr_count,
    COUNT(*) FILTER (WHERE review_turnaround_hours <= 24) as within_24h_count,
    COUNT(*) FILTER (WHERE review_turnaround_hours > 24) as over_24h_count
FROM pr_metrics
WHERE first_review_at IS NOT NULL
GROUP BY team_id, month
ORDER BY team_id, month DESC;

-- Create view for review engagement metrics (PostgreSQL version)
CREATE OR REPLACE VIEW view_review_engagement AS
SELECT 
    team_id,
    DATE_TRUNC('month', created_at) as month,
    AVG(review_comments_count) as avg_comments_per_pr,
    AVG(conversation_count) as avg_conversations_per_pr,
    AVG(reviewers_count) as avg_reviewers_per_pr,
    SUM(changes_requested_count)::DECIMAL * 100.0 / NULLIF(SUM(changes_requested_count + approved_count), 0) as changes_requested_rate,
    SUM(approved_count)::DECIMAL * 100.0 / NULLIF(SUM(changes_requested_count + approved_count), 0) as approval_rate,
    COUNT(*) as pr_count
FROM pr_metrics
GROUP BY team_id, month
ORDER BY team_id, month DESC;

-- Create view for knowledge sharing metrics (PostgreSQL version)
CREATE OR REPLACE VIEW view_knowledge_sharing AS
SELECT 
    team_id,
    DATE_TRUNC('month', created_at) as month,
    AVG(reviewers_count) as avg_reviewers,
    AVG(external_reviewers_count) as avg_external_reviewers,
    AVG(external_reviewers_count::DECIMAL * 100.0 / NULLIF(reviewers_count, 0)) as external_reviewer_rate,
    SUM(external_reviewers_count) as total_external_reviews,
    COUNT(*) FILTER (WHERE external_reviewers_count > 0) as prs_with_external_reviews,
    COUNT(*) as pr_count
FROM pr_metrics
WHERE reviewers_count > 0
GROUP BY team_id, month
ORDER BY team_id, month DESC;
