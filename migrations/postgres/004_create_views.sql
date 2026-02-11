-- Create view for team velocity (PRs merged per week)
CREATE OR REPLACE VIEW view_team_velocity AS
SELECT 
    team_id,
    DATE_TRUNC('week', merged_at) as week,
    COUNT(*) as prs_merged,
    AVG(cycle_time_hours) as avg_cycle_time_hours
FROM pr_metrics
WHERE merged_at IS NOT NULL
GROUP BY team_id, week
ORDER BY team_id, week DESC;

-- Create view for DORA lead time (median and P95)
CREATE OR REPLACE VIEW view_dora_lead_time AS
SELECT 
    team_id,
    DATE_TRUNC('month', merged_at) as month,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY cycle_time_hours) as median_lead_time_hours,
    PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY cycle_time_hours) as p95_lead_time_hours,
    COUNT(*) as pr_count
FROM pr_metrics
WHERE merged_at IS NOT NULL
GROUP BY team_id, month
ORDER BY team_id, month DESC;
