# PR Review Metrics - Quick Reference

## New Metrics Available

### 1. Review Turnaround Time
**Query the view:**
```sql
SELECT * FROM view_review_turnaround 
WHERE team_id = 1 
ORDER BY month DESC 
LIMIT 6;
```

**Metrics provided:**
- Average turnaround hours
- Min/max turnaround hours
- Count of PRs reviewed within 24 hours
- Count of PRs taking over 24 hours

---

### 2. Review Engagement
**Query the view:**
```sql
SELECT * FROM view_review_engagement 
WHERE team_id = 1 
ORDER BY month DESC 
LIMIT 6;
```

**Metrics provided:**
- Average comments per PR
- Average conversations per PR
- Average reviewers per PR
- Changes requested rate (%)
- Approval rate (%)

---

### 3. Knowledge Sharing
**Query the view:**
```sql
SELECT * FROM view_knowledge_sharing 
WHERE team_id = 1 
ORDER BY month DESC 
LIMIT 6;
```

**Metrics provided:**
- Average reviewers per PR
- Average external reviewers (outside team)
- External reviewer rate (%)
- Total external reviews
- PRs with external reviews

---

## Database Schema

### New Columns in `pr_metrics` Table

| Column | Type | Description |
|--------|------|-------------|
| `first_review_at` | DATETIME | Timestamp of first review |
| `review_turnaround_hours` | INTEGER | Hours from PR creation to first review |
| `review_comments_count` | INTEGER | Total review comments |
| `conversation_count` | INTEGER | Total conversation threads |
| `changes_requested_count` | INTEGER | Number of "changes requested" reviews |
| `approved_count` | INTEGER | Number of "approved" reviews |
| `reviewers_count` | INTEGER | Total unique reviewers |
| `external_reviewers_count` | INTEGER | Reviewers outside the team |
| `reviewers_list` | TEXT/JSONB | JSON array of reviewer usernames |

---

## GitHub API Data Required

To populate these metrics in Phase 2, fetch:

### PR Reviews
```
GET /repos/{owner}/{repo}/pulls/{pull_number}/reviews
```
Extract: `submitted_at`, `state`, `user.login`

### PR Comments
```
GET /repos/{owner}/{repo}/pulls/{pull_number}/comments
```
Extract: comment count, conversation threads

---

## Example Queries

### Find PRs with slow review turnaround
```sql
SELECT repository, pr_number, title, review_turnaround_hours
FROM pr_metrics
WHERE review_turnaround_hours > 48
ORDER BY review_turnaround_hours DESC;
```

### Find most engaged PRs (lots of comments)
```sql
SELECT repository, pr_number, title, review_comments_count, reviewers_count
FROM pr_metrics
WHERE review_comments_count > 10
ORDER BY review_comments_count DESC;
```

### Find PRs with cross-team collaboration
```sql
SELECT repository, pr_number, title, external_reviewers_count, reviewers_list
FROM pr_metrics
WHERE external_reviewers_count > 0
ORDER BY external_reviewers_count DESC;
```

### Monthly review quality metrics
```sql
SELECT 
    month,
    avg_turnaround_hours,
    avg_comments_per_pr,
    changes_requested_rate,
    avg_external_reviewers
FROM view_review_turnaround
JOIN view_review_engagement USING (team_id, month)
JOIN view_knowledge_sharing USING (team_id, month)
WHERE team_id = 1
ORDER BY month DESC;
```
