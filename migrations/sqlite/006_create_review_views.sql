-- Create view for review turnaround metrics
CREATE VIEW IF NOT EXISTS view_review_turnaround AS
SELECT 
    team_id,
    strftime('%Y-%m', created_at) as month,
    AVG(review_turnaround_hours) as avg_turnaround_hours,
    MIN(review_turnaround_hours) as min_turnaround_hours,
    MAX(review_turnaround_hours) as max_turnaround_hours,
    COUNT(*) as pr_count,
    COUNT(CASE WHEN review_turnaround_hours <= 24 THEN 1 END) as within_24h_count,
    COUNT(CASE WHEN review_turnaround_hours > 24 THEN 1 END) as over_24h_count
FROM pr_metrics
WHERE first_review_at IS NOT NULL
GROUP BY team_id, month
ORDER BY team_id, month DESC;

-- Create view for review engagement metrics
CREATE VIEW IF NOT EXISTS view_review_engagement AS
SELECT 
    team_id,
    strftime('%Y-%m', created_at) as month,
    AVG(review_comments_count) as avg_comments_per_pr,
    AVG(conversation_count) as avg_conversations_per_pr,
    AVG(reviewers_count) as avg_reviewers_per_pr,
    SUM(changes_requested_count) * 100.0 / NULLIF(SUM(changes_requested_count + approved_count), 0) as changes_requested_rate,
    SUM(approved_count) * 100.0 / NULLIF(SUM(changes_requested_count + approved_count), 0) as approval_rate,
    COUNT(*) as pr_count
FROM pr_metrics
GROUP BY team_id, month
ORDER BY team_id, month DESC;

-- Create view for knowledge sharing metrics
CREATE VIEW IF NOT EXISTS view_knowledge_sharing AS
SELECT 
    team_id,
    strftime('%Y-%m', created_at) as month,
    AVG(reviewers_count) as avg_reviewers,
    AVG(external_reviewers_count) as avg_external_reviewers,
    AVG(external_reviewers_count * 100.0 / NULLIF(reviewers_count, 0)) as external_reviewer_rate,
    SUM(external_reviewers_count) as total_external_reviews,
    COUNT(CASE WHEN external_reviewers_count > 0 THEN 1 END) as prs_with_external_reviews,
    COUNT(*) as pr_count
FROM pr_metrics
WHERE reviewers_count > 0
GROUP BY team_id, month
ORDER BY team_id, month DESC;
