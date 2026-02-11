# Roadmap

## Project Timeline

**Project Start**: February 2026  
**Target MVP**: March 2026 (4 weeks)  
**Target Production**: April 2026 (8 weeks)

---

## Milestones

### Milestone 1: Foundation & Setup
**Target Date**: Week 1 (Feb 11-17, 2026)  
**Status**: ✅ Complete

#### Deliverables
- [x] Project structure established
- [x] Go module initialized
- [x] Documentation framework created
- [x] Core dependencies installed and configured
- [x] PostgreSQL database schema created
- [x] SQLite database schema created (for local dev)
- [x] Database migration scripts implemented
- [x] Configuration management system built
- [x] Local development environment documented

#### Success Criteria
- ✅ All dependencies compile without errors
- ✅ Database schema can be provisioned from scratch
- ✅ Configuration can be loaded from environment variables
- ✅ Developer can run the project locally

---

### Milestone 2: Data Collection (Collector MVP)
**Target Date**: Week 2-3 (Feb 18 - Mar 3, 2026)  
**Status**: ⚪ Not Started

#### Deliverables
- [ ] GitHub OAuth2 authentication implemented
- [ ] Team membership sync module completed
- [ ] GitHub PR polling logic implemented
- [ ] PR filtering and attribution engine built
- [ ] Database persistence layer completed
- [ ] Idempotent upsert operations working
- [ ] Basic error handling and logging
- [ ] Unit tests for core logic (90%+ coverage)

#### Success Criteria
- ✅ Can authenticate with GitHub API
- ✅ Can fetch PRs from configured repositories
- ✅ Can correctly filter PRs by team membership
- ✅ Can store PR metrics in database without duplicates
- ✅ Can handle GitHub rate limits gracefully
- ✅ Re-running the collector doesn't create duplicate data

---

### Milestone 3: Metrics & API
**Target Date**: Week 3-4 (Mar 4-10, 2026)  
**Status**: ⚪ Not Started

#### Deliverables
- [ ] SQL views for team velocity created and tested
- [ ] SQL views for DORA lead time created and tested
- [ ] Lambda handler for API Gateway implemented
- [ ] API endpoint: `GET /metrics/team/{team_id}/velocity`
- [ ] API endpoint: `GET /metrics/team/{team_id}/lead-time`
- [ ] API authentication (API key) implemented
- [ ] API integration tests completed
- [ ] API documentation created

#### Success Criteria
- ✅ SQL views return accurate metrics
- ✅ API endpoints return correct JSON responses
- ✅ API response time < 500ms for typical queries
- ✅ Unauthorized requests are rejected
- ✅ API handles edge cases (no data, invalid team_id)

---

### Milestone 4: Deployment & Automation
**Target Date**: Week 4-5 (Mar 11-17, 2026)  
**Status**: ⚪ Not Started

#### Deliverables
- [ ] AWS Lambda function created (ARM64)
- [ ] EventBridge scheduled rule configured (every 4 hours)
- [ ] GitHub Actions CI/CD workflow implemented
- [ ] Automated tests run on every commit
- [ ] Automated deployment to Lambda on merge to main
- [ ] CloudWatch logging configured
- [ ] CloudWatch metrics and alarms set up
- [ ] Infrastructure documented

#### Success Criteria
- ✅ Lambda function runs successfully on schedule
- ✅ CI/CD pipeline deploys without manual intervention
- ✅ CloudWatch logs show successful execution
- ✅ Alarms trigger on failures
- ✅ System runs for 1 week without issues

---

### Milestone 5: Production Hardening
**Target Date**: Week 6-8 (Mar 18 - Apr 7, 2026)  
**Status**: ⚪ Not Started

#### Deliverables
- [ ] GitHub PAT moved to AWS Secrets Manager
- [ ] Database connection pooling optimized
- [ ] Retry logic and exponential backoff refined
- [ ] Dead letter queue for failed processing
- [ ] Performance testing completed
- [ ] Security audit completed
- [ ] Load testing completed
- [ ] Production runbook created
- [ ] Monitoring dashboard created

#### Success Criteria
- ✅ Zero data loss over 1 week of operation
- ✅ Zero duplicate data over 1 week of operation
- ✅ 99%+ uptime for scheduled runs
- ✅ All security best practices implemented
- ✅ Performance meets NFR requirements
- ✅ Team can troubleshoot issues using runbook

---

## Future Enhancements

### Phase 6: Advanced Metrics (Q2 2026)
- [ ] Deployment frequency tracking
- [ ] Change failure rate calculation
- [ ] Mean time to recovery (MTTR) tracking
- [ ] Custom metric definitions
- [ ] Historical trend analysis
- [ ] Anomaly detection

### Phase 7: Frontend Dashboard (Q3 2026)
- [ ] Vue.js dashboard application
- [ ] Team selection and filtering UI
- [ ] Date range selection
- [ ] Interactive charts (velocity, lead time)
- [ ] Drill-down to individual PRs
- [ ] Export metrics to CSV/PDF
- [ ] User authentication and authorization

### Phase 8: Advanced Features (Q4 2026)
- [ ] Slack/Teams notifications for metric thresholds
- [ ] Automated team membership sync from HR systems
- [ ] Multi-repository aggregation views
- [ ] Comparison across teams
- [ ] Goal setting and tracking
- [ ] Integration with Jira for issue tracking
- [ ] Support for GitLab and Bitbucket

### Phase 9: Enterprise Features (2027)
- [ ] Multi-tenant support
- [ ] Role-based access control (RBAC)
- [ ] Audit logging
- [ ] Custom branding
- [ ] SLA monitoring
- [ ] Advanced analytics and ML predictions

---

## Dependencies & Risks

### Critical Dependencies
- GitHub API availability and stability
- AWS Lambda and EventBridge reliability
- PostgreSQL database availability

### Known Risks
- **GitHub Rate Limits**: Mitigation - implement aggressive caching and respect limits
- **Lambda Cold Starts**: Mitigation - acceptable for scheduled jobs, consider provisioned concurrency if needed
- **Data Quality**: Mitigation - implement validation and data quality checks
- **Team Adoption**: Mitigation - ensure metrics are accurate and valuable to engineering managers
