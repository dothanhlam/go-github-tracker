// ===== DORA Metrics Dashboard App =====

(function () {
  'use strict';

  // --- State ---
  let apiKey = localStorage.getItem('dora_api_key') || '';
  let currentTeamId = null;
  let charts = {};

  // --- DOM refs ---
  const teamListEl = document.getElementById('teamList');
  const emptyStateEl = document.getElementById('emptyState');
  const dashboardEl = document.getElementById('dashboardContent');
  const teamNameEl = document.getElementById('teamName');
  const teamMetaEl = document.getElementById('teamMeta');
  const startDateEl = document.getElementById('startDate');
  const endDateEl = document.getElementById('endDate');
  const refreshBtn = document.getElementById('refreshBtn');
  const statusBadge = document.getElementById('connectionStatus');

  // --- Chart.js defaults ---
  Chart.defaults.font.family = "-apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif";
  Chart.defaults.font.size = 12;
  Chart.defaults.plugins.legend.labels.usePointStyle = true;
  Chart.defaults.plugins.legend.labels.pointStyleWidth = 10;
  Chart.defaults.animation.duration = 400;

  const COLORS = {
    green: '#1a7f37',
    greenBg: 'rgba(26,127,55,0.1)',
    blue: '#0969da',
    blueBg: 'rgba(9,105,218,0.1)',
    purple: '#8250df',
    purpleBg: 'rgba(130,80,223,0.1)',
    orange: '#bf8700',
    orangeBg: 'rgba(191,135,0,0.1)',
    red: '#cf222e',
    redBg: 'rgba(207,34,46,0.1)',
    gray: '#656d76',
    grayBg: 'rgba(101,109,118,0.1)',
  };

  // --- Init ---
  function init() {
    // Set default dates (last 30 days)
    const now = new Date();
    const thirtyDaysAgo = new Date(now);
    thirtyDaysAgo.setDate(now.getDate() - 30);
    startDateEl.value = formatDate(thirtyDaysAgo);
    endDateEl.value = formatDate(now);

    // Wire up events
    refreshBtn.addEventListener('click', () => loadTeamData(currentTeamId));
    document.getElementById('saveApiKeyBtn').addEventListener('click', saveApiKey);
    document.getElementById('apiKeyInput').addEventListener('keydown', (e) => {
      if (e.key === 'Enter') saveApiKey();
    });

    // Load
    if (apiKey) {
      document.getElementById('apiKeyInput').value = apiKey;
      loadTeams();
    } else {
      showStatus('No API key', 'secondary');
      teamListEl.innerHTML = '<div class="text-muted small p-3">Set your API key to get started.</div>';
      // Show modal
      const modal = new bootstrap.Modal(document.getElementById('apiKeyModal'));
      modal.show();
    }
  }

  // --- API Key ---
  function saveApiKey() {
    const input = document.getElementById('apiKeyInput');
    apiKey = input.value.trim();
    if (!apiKey) return;
    localStorage.setItem('dora_api_key', apiKey);
    bootstrap.Modal.getInstance(document.getElementById('apiKeyModal')).hide();
    loadTeams();
  }

  // --- API Helpers ---
  async function apiFetch(path) {
    const res = await fetch(path, {
      headers: { 'X-API-Key': apiKey }
    });
    if (res.status === 401) {
      showStatus('Unauthorized', 'danger');
      throw new Error('Unauthorized');
    }
    if (!res.ok) throw new Error(`API error: ${res.status}`);
    return res.json();
  }

  function showStatus(text, variant) {
    statusBadge.textContent = text;
    statusBadge.className = `badge bg-${variant}`;
  }

  function formatDate(d) {
    return d.toISOString().split('T')[0];
  }

  // --- Load Teams ---
  async function loadTeams() {
    try {
      showStatus('Connecting…', 'warning');
      const data = await apiFetch('/api/v1/teams');
      showStatus('Connected', 'success');

      const teams = data.teams || [];
      if (teams.length === 0) {
        teamListEl.innerHTML = '<div class="text-muted small p-3">No teams found.</div>';
        return;
      }

      teamListEl.innerHTML = '';
      teams.forEach(team => {
        const item = document.createElement('a');
        item.href = '#';
        item.className = 'list-group-item list-group-item-action';
        item.textContent = team.name || `Team ${team.id}`;
        item.dataset.teamId = team.id;
        item.dataset.teamName = team.name || `Team ${team.id}`;
        item.dataset.memberCount = team.member_count || 0;
        item.addEventListener('click', (e) => {
          e.preventDefault();
          selectTeam(team.id, team.name || `Team ${team.id}`, team.member_count || 0);
        });
        teamListEl.appendChild(item);
      });

      // Auto-select first team
      if (teams.length > 0) {
        selectTeam(teams[0].id, teams[0].name || `Team ${teams[0].id}`, teams[0].member_count || 0);
      }

    } catch (err) {
      if (err.message !== 'Unauthorized') {
        showStatus('Error', 'danger');
      }
      teamListEl.innerHTML = '<div class="text-danger small p-3">Failed to load teams. Check API key and server.</div>';
    }
  }

  // --- Select Team ---
  function selectTeam(teamId, teamName, memberCount) {
    currentTeamId = teamId;
    teamNameEl.textContent = teamName;
    teamMetaEl.textContent = `${memberCount} member${memberCount !== 1 ? 's' : ''} · Team ID: ${teamId}`;

    // Update sidebar active state
    document.querySelectorAll('#teamList .list-group-item').forEach(el => {
      el.classList.toggle('active', parseInt(el.dataset.teamId) === teamId);
    });

    emptyStateEl.style.display = 'none';
    dashboardEl.style.display = 'block';

    loadTeamData(teamId);
  }

  // --- Load Team Data ---
  async function loadTeamData(teamId) {
    if (!teamId) return;

    const start = startDateEl.value;
    const end = endDateEl.value;
    const qs = `start_date=${start}&end_date=${end}`;

    // Show loading on charts
    Object.values(charts).forEach(c => c.destroy());
    charts = {};

    try {
      const [velocity, leadTime, turnaround, engagement, knowledgeSharing, commits] = await Promise.all([
        apiFetch(`/api/v1/teams/${teamId}/velocity?${qs}`),
        apiFetch(`/api/v1/teams/${teamId}/lead-time?${qs}`),
        apiFetch(`/api/v1/teams/${teamId}/review-turnaround?${qs}`),
        apiFetch(`/api/v1/teams/${teamId}/review-engagement?${qs}`),
        apiFetch(`/api/v1/teams/${teamId}/knowledge-sharing?${qs}`),
        apiFetch(`/api/v1/teams/${teamId}/commits?${qs}`),
      ]);

      renderSummaryStats(velocity, turnaround, engagement);
      renderVelocityChart(velocity);
      renderLeadTimeChart(leadTime);
      renderTurnaroundChart(turnaround);
      renderEngagementChart(engagement);
      renderKnowledgeSharingChart(knowledgeSharing);
      renderCommitChart(commits);

    } catch (err) {
      console.error('Failed to load team data:', err);
    }
  }

  // --- Summary Stats ---
  function renderSummaryStats(velocity, turnaround, engagement) {
    const vMetrics = velocity.metrics || [];
    const tMetrics = turnaround.metrics || [];
    const eMetrics = engagement.metrics || [];

    // Total PRs merged
    const totalPRs = vMetrics.reduce((sum, m) => sum + (m.prs_merged || 0), 0);
    document.getElementById('statPRsMerged').textContent = totalPRs;

    // Average cycle time
    const cycleTimes = vMetrics.filter(m => m.avg_cycle_time_hours > 0).map(m => m.avg_cycle_time_hours);
    const avgCycle = cycleTimes.length > 0 ? (cycleTimes.reduce((a, b) => a + b, 0) / cycleTimes.length).toFixed(1) : '—';
    document.getElementById('statAvgCycleTime').textContent = avgCycle;

    // Average review turnaround
    const turnarounds = tMetrics.filter(m => m.avg_turnaround_hours > 0).map(m => m.avg_turnaround_hours);
    const avgTurn = turnarounds.length > 0 ? (turnarounds.reduce((a, b) => a + b, 0) / turnarounds.length).toFixed(1) : '—';
    document.getElementById('statAvgReviewTime').textContent = avgTurn;

    // Total unique reviewers (max across periods)
    const maxReviewers = eMetrics.reduce((max, m) => Math.max(max, m.unique_reviewers || 0), 0);
    document.getElementById('statUniqueReviewers').textContent = maxReviewers || '—';
  }

  // --- Charts ---

  function renderVelocityChart(data) {
    const metrics = data.metrics || [];
    const ctx = document.getElementById('velocityChart').getContext('2d');
    charts.velocity = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: metrics.map(m => m.period),
        datasets: [{
          label: 'PRs Merged',
          data: metrics.map(m => m.prs_merged),
          backgroundColor: COLORS.greenBg,
          borderColor: COLORS.green,
          borderWidth: 2,
          borderRadius: 4,
        }, {
          label: 'Avg Cycle Time (hrs)',
          data: metrics.map(m => m.avg_cycle_time_hours),
          type: 'line',
          borderColor: COLORS.blue,
          backgroundColor: COLORS.blueBg,
          tension: 0.3,
          yAxisID: 'y1',
          pointRadius: 3,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: { beginAtZero: true, title: { display: true, text: 'PRs Merged' } },
          y1: { beginAtZero: true, position: 'right', grid: { drawOnChartArea: false }, title: { display: true, text: 'Cycle Time (hrs)' } },
        }
      }
    });
  }

  function renderLeadTimeChart(data) {
    const metrics = data.metrics || [];
    const ctx = document.getElementById('leadTimeChart').getContext('2d');
    charts.leadTime = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: metrics.map(m => m.period),
        datasets: [{
          label: 'Median Lead Time (hrs)',
          data: metrics.map(m => m.median_lead_time_hours),
          backgroundColor: COLORS.purpleBg,
          borderColor: COLORS.purple,
          borderWidth: 2,
          borderRadius: 4,
        }, {
          label: 'P95 Lead Time (hrs)',
          data: metrics.map(m => m.p95_lead_time_hours),
          backgroundColor: COLORS.redBg,
          borderColor: COLORS.red,
          borderWidth: 2,
          borderRadius: 4,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: { y: { beginAtZero: true, title: { display: true, text: 'Hours' } } }
      }
    });
  }

  function renderTurnaroundChart(data) {
    const metrics = data.metrics || [];
    const ctx = document.getElementById('reviewTurnaroundChart').getContext('2d');
    charts.turnaround = new Chart(ctx, {
      type: 'line',
      data: {
        labels: metrics.map(m => m.period),
        datasets: [{
          label: 'Avg Turnaround (hrs)',
          data: metrics.map(m => m.avg_turnaround_hours),
          borderColor: COLORS.purple,
          backgroundColor: COLORS.purpleBg,
          fill: true,
          tension: 0.3,
          pointRadius: 3,
        }, {
          label: 'Median Turnaround (hrs)',
          data: metrics.map(m => m.median_turnaround_hours),
          borderColor: COLORS.blue,
          backgroundColor: COLORS.blueBg,
          fill: false,
          tension: 0.3,
          pointRadius: 3,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: { y: { beginAtZero: true, title: { display: true, text: 'Hours' } } }
      }
    });
  }

  function renderEngagementChart(data) {
    const metrics = data.metrics || [];
    const ctx = document.getElementById('engagementChart').getContext('2d');
    charts.engagement = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: metrics.map(m => m.period),
        datasets: [{
          label: 'Total Reviews',
          data: metrics.map(m => m.total_reviews),
          backgroundColor: COLORS.orangeBg,
          borderColor: COLORS.orange,
          borderWidth: 2,
          borderRadius: 4,
        }, {
          label: 'Unique Reviewers',
          data: metrics.map(m => m.unique_reviewers),
          backgroundColor: COLORS.blueBg,
          borderColor: COLORS.blue,
          borderWidth: 2,
          borderRadius: 4,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: { y: { beginAtZero: true } }
      }
    });
  }

  function renderKnowledgeSharingChart(data) {
    const metrics = data.metrics || [];
    const ctx = document.getElementById('knowledgeSharingChart').getContext('2d');
    charts.knowledgeSharing = new Chart(ctx, {
      type: 'line',
      data: {
        labels: metrics.map(m => m.period),
        datasets: [{
          label: 'Cross-Team Reviews',
          data: metrics.map(m => m.cross_team_reviews),
          borderColor: COLORS.green,
          backgroundColor: COLORS.greenBg,
          fill: true,
          tension: 0.3,
          pointRadius: 3,
        }, {
          label: 'Knowledge Sharing Score',
          data: metrics.map(m => m.knowledge_sharing_score),
          borderColor: COLORS.orange,
          backgroundColor: COLORS.orangeBg,
          fill: false,
          tension: 0.3,
          yAxisID: 'y1',
          pointRadius: 3,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: { beginAtZero: true, title: { display: true, text: 'Reviews' } },
          y1: { beginAtZero: true, position: 'right', grid: { drawOnChartArea: false }, title: { display: true, text: 'Score' } },
        }
      }
    });
  }

  function renderCommitChart(data) {
    const metrics = data.metrics || [];
    const ctx = document.getElementById('commitChart').getContext('2d');
    charts.commits = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: metrics.map(m => m.period),
        datasets: [{
          label: 'Commits',
          data: metrics.map(m => m.commits_count),
          backgroundColor: COLORS.grayBg,
          borderColor: COLORS.gray,
          borderWidth: 2,
          borderRadius: 4,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: { y: { beginAtZero: true, title: { display: true, text: 'Commits' } } }
      }
    });
  }

  // --- Boot ---
  document.addEventListener('DOMContentLoaded', init);

})();
