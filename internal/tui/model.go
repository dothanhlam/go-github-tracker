package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dothanhlam/go-github-tracker/internal/service"
)

var (
	docStyle         = lipgloss.NewStyle().Margin(1, 2)
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Padding(0, 1)
	bodyStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(1, 0)
	highlightStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	metricTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	metricValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Padding(0, 1)
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Padding(1, 0)
	paneStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
)

// TeamItem represents an item in the list
type TeamItem struct {
	team service.Team
}

func (i TeamItem) Title() string       { return i.team.Name }
func (i TeamItem) Description() string { return fmt.Sprintf("%d members", i.team.MemberCount) }
func (i TeamItem) FilterValue() string { return i.team.Name }

// Model is the main state of the TUI
type Model struct {
	metricsService *service.MetricsService
	teamList       list.Model
	selectedTeam   *service.Team
	metricsData    *MetricsData
	err            error
	loading        bool
	width          int
	height         int
}

type MetricsData struct {
	Velocity         *service.VelocityResponse
	LeadTime         *service.LeadTimeResponse
	ReviewTurnaround *service.ReviewTurnaroundResponse
	ReviewEngagement *service.ReviewEngagementResponse
	KnowledgeSharing *service.KnowledgeSharingResponse
}

// InitialModel initializes the state
func InitialModel(metricsService *service.MetricsService) *Model {
	// Initialize empty list
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Teams"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return &Model{
		metricsService: metricsService,
		teamList:       l,
		loading:        true,
	}
}

// Init acts a setup step
func (m *Model) Init() tea.Cmd {
	return m.fetchTeamsCmd()
}

// Messages
type teamsFetchedMsg struct {
	teams []service.Team
}
type metricsFetchedMsg struct {
	data *MetricsData
}
type errMsg struct{ err error }

// fetchTeamsCmd is a command to fetch all teams
func (m *Model) fetchTeamsCmd() tea.Cmd {
	return func() tea.Msg {
		teams, err := m.metricsService.ListTeams()
		if err != nil {
			return errMsg{err}
		}
		return teamsFetchedMsg{teams}
	}
}

// fetchMetricsCmd is a command to fetch metrics for a specific team
func (m *Model) fetchMetricsCmd(teamID int) tea.Cmd {
	return func() tea.Msg {
		// Use a 3-month window by default
		endDate := time.Now()
		startDate := endDate.AddDate(0, -3, 0)

		data := &MetricsData{}
		
		// Note: the MetricsService methods return (response, error)
		// We ignore individual errors to show what we can, but catch any major panic
		vel, _ := m.metricsService.GetTeamVelocity(teamID, startDate, endDate, "month")
		if vel != nil {
			data.Velocity = vel
		}

		lt, _ := m.metricsService.GetTeamLeadTime(teamID, startDate, endDate, "month")
		if lt != nil {
			data.LeadTime = lt
		}

		rt, _ := m.metricsService.GetReviewTurnaround(teamID, startDate, endDate)
		if rt != nil {
			data.ReviewTurnaround = rt
		}

		re, _ := m.metricsService.GetReviewEngagement(teamID, startDate, endDate)
		if re != nil {
			data.ReviewEngagement = re
		}

		ks, _ := m.metricsService.GetKnowledgeSharing(teamID, startDate, endDate)
		if ks != nil {
			data.KnowledgeSharing = ks
		}

		return metricsFetchedMsg{data}
	}
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// List gets roughly 1/3 of width
		h, v := docStyle.GetFrameSize()
		m.teamList.SetSize(m.width/3-h, m.height-v)

	case tea.KeyMsg:
		if m.teamList.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case teamsFetchedMsg:
		items := make([]list.Item, len(msg.teams))
		for i, t := range msg.teams {
			items[i] = TeamItem{team: t}
		}
		m.teamList.SetItems(items)
		m.loading = false
		
		// Fetch metrics for first team if available
		if len(msg.teams) > 0 {
			m.selectedTeam = &msg.teams[0]
			m.loading = true
			return m, m.fetchMetricsCmd(m.selectedTeam.ID)
		}

	case metricsFetchedMsg:
		m.metricsData = msg.data
		m.loading = false

	case errMsg:
		m.err = msg.err
		m.loading = false
	}

	var cmd tea.Cmd
	var listCmd tea.Cmd
	
	// Keep track of previously selected item index
	prevIndex := m.teamList.Index()
	m.teamList, listCmd = m.teamList.Update(msg)
	
	// If index changed, fetch metrics for new team
	if m.teamList.Index() != prevIndex {
		if i, ok := m.teamList.SelectedItem().(TeamItem); ok {
			m.selectedTeam = &i.team
			m.loading = true
			m.metricsData = nil
			cmd = m.fetchMetricsCmd(i.team.ID)
		}
	}

	return m, tea.Batch(cmd, listCmd)
}

// View renders the TUI
func (m *Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v\nPress q to quit", m.err))
	}

	// Left pane: Team list
	listPane := paneStyle.
		Width(m.width/3 - 4).
		Height(m.height - 4).
		Render(m.teamList.View())

	// Right pane: Metrics view
	var details string
	if m.loading && m.selectedTeam == nil {
		details = "Loading teams..."
	} else if m.loading {
		details = fmt.Sprintf("Loading metrics for %s...", m.selectedTeam.Name)
	} else if m.selectedTeam == nil {
		details = "No teams found."
	} else if m.metricsData != nil {
		details = m.renderMetricsView()
	}

	rightPaneWidth := (m.width * 2 / 3) - 4
	detailsPane := paneStyle.
		Width(rightPaneWidth).
		Height(m.height - 4).
		Render(details)

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, listPane, detailsPane))
}

func (m *Model) renderMetricsView() string {
	s := ""
	s += titleStyle.Render(fmt.Sprintf("📊 DORA Metrics Dashboard :: %s", m.selectedTeam.Name)) + "\n\n"

	// Velocity
	s += metricTitleStyle.Render("🚀 Team Velocity") + "\n"
	if m.metricsData.Velocity != nil && len(m.metricsData.Velocity.Metrics) > 0 {
		for _, v := range m.metricsData.Velocity.Metrics {
			s += fmt.Sprintf("  %s: %s PRs merged, %s cycle time\n", 
				bodyStyle.Render(v.Period), 
				highlightStyle.Render(fmt.Sprintf("%d", v.PRsMerged)),
				highlightStyle.Render(fmt.Sprintf("%.1fh", v.AvgCycleTimeHrs)))
		}
	} else {
		s += bodyStyle.Render("  No data available") + "\n"
	}
	s += "\n"

	// Lead Time
	s += metricTitleStyle.Render("⏱️  Lead Time for Changes") + "\n"
	if m.metricsData.LeadTime != nil && len(m.metricsData.LeadTime.Metrics) > 0 {
		for _, lt := range m.metricsData.LeadTime.Metrics {
			s += fmt.Sprintf("  %s: %s median, %s p95\n", 
				bodyStyle.Render(lt.Period),
				highlightStyle.Render(fmt.Sprintf("%.1fh", lt.MedianLeadTimeHrs)),
				highlightStyle.Render(fmt.Sprintf("%.1fh", lt.P95LeadTimeHrs)))
		}
	} else {
		s += bodyStyle.Render("  No data available") + "\n"
	}
	s += "\n"

	// Technical Quality / Review Turnaround
	s += metricTitleStyle.Render("🔍 Review Turnaround") + "\n"
	if m.metricsData.ReviewTurnaround != nil && len(m.metricsData.ReviewTurnaround.Metrics) > 0 {
		for _, rt := range m.metricsData.ReviewTurnaround.Metrics {
			s += fmt.Sprintf("  %s: %s avg turnaround\n", 
				bodyStyle.Render(rt.Period),
				highlightStyle.Render(fmt.Sprintf("%.1fh", rt.AvgTurnaroundHrs)))
		}
	} else {
		s += bodyStyle.Render("  No data available") + "\n"
	}
	s += "\n"

	// Review Engagement
	s += metricTitleStyle.Render("💬 Review Engagement") + "\n"
	if m.metricsData.ReviewEngagement != nil && len(m.metricsData.ReviewEngagement.Metrics) > 0 {
		for _, re := range m.metricsData.ReviewEngagement.Metrics {
			s += fmt.Sprintf("  %s: %s avg reviews/PR\n", 
				bodyStyle.Render(re.Period),
				highlightStyle.Render(fmt.Sprintf("%.1f", re.AvgReviewsPerPR)))
		}
	} else {
		s += bodyStyle.Render("  No data available") + "\n"
	}

	return s
}
