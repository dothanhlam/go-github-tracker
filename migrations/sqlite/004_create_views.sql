-- Create view for team velocity (PRs merged per week)
CREATE VIEW IF NOT EXISTS view_team_velocity AS
SELECT 
    team_id,
    DATE(merged_at, 'weekday 0', '-6 days') as week,
    COUNT(*) as prs_merged,
    AVG(cycle_time_hours) as avg_cycle_time_hours
FROM pr_metrics
WHERE merged_at IS NOT NULL
GROUP BY team_id, week
ORDER BY team_id, week DESC;

-- Create view for DORA lead time (median and P95)
-- Note: SQLite doesn't have PERCENTILE_CONT, so we use a simpler approach
CREATE VIEW IF NOT EXISTS view_dora_lead_time AS
SELECT 
    team_id,
    strftime('%Y-%m', merged_at) as month,
    AVG(cycle_time_hours) as avg_lead_time_hours,
    MIN(cycle_time_hours) as min_lead_time_hours,
    MAX(cycle_time_hours) as max_lead_time_hours,
    COUNT(*) as pr_count
FROM pr_metrics
WHERE merged_at IS NOT NULL
GROUP BY team_id, month
ORDER BY team_id, month DESC;
