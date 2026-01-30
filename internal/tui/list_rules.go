package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// allRulesMsg carries all rules from the system
type allRulesMsg struct {
	natRules      []models.NATRule
	firewallRules []models.FirewallRule
}

// loadAllRules loads all rules from the firewall provider
func (m *Model) loadAllRules() tea.Cmd {
	return func() tea.Msg {
		var natRules []models.NATRule
		var fwRules []models.FirewallRule

		if m.provider != nil {
			// Load NAT rules from provider
			nat, err := m.provider.ListNATRules(m.ctx)
			if err == nil {
				natRules = nat
			}

			// Load firewall rules from provider
			fw, err := m.provider.ListFirewallRules(m.ctx)
			if err == nil {
				fwRules = fw
			}
		}

		return allRulesMsg{natRules, fwRules}
	}
}

// CombinedRule represents a unified view of any rule type
type CombinedRule struct {
	ID          string
	Type        string
	Product     string
	Details     string
	Description string
	Raw         interface{}
}

// getRulesVisibleHeight returns the number of rows that can fit on screen
func (m *Model) getRulesVisibleHeight() int {
	// Account for: title, subtitle, empty line, header, separator, empty line, help, status bar
	// Approximate: header (2) + margins (6) + help (1) + status (1) = ~10 lines
	availableHeight := m.height - 14
	if availableHeight < 5 {
		return 5 // Minimum visible rows
	}
	return availableHeight
}

// updateRulesScrollOffset ensures the view shows the relevant content
func (m *Model) updateRulesScrollOffset() {
	visibleHeight := m.getRulesVisibleHeight()
	var listLen int
	if m.ruleViewMode == "nat" {
		listLen = len(m.natRules)
	} else {
		listLen = len(m.firewallRules)
	}

	// Ensure scroll offset doesn't exceed list bounds
	if m.rulesScrollOffset > listLen-visibleHeight {
		m.rulesScrollOffset = listLen - visibleHeight
	}
	if m.rulesScrollOffset < 0 {
		m.rulesScrollOffset = 0
	}
}

// updateListRules handles list rules updates
func (m *Model) updateListRules(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case allRulesMsg:
		m.natRules = msg.natRules
		m.firewallRules = msg.firewallRules
		m.rulesScrollOffset = 0 // Reset scroll on reload
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("d", "delete"))) {
			return m.deleteFirstRule()
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("r"))) {
			return m, m.loadAllRules()
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("a"))) {
			m.screen = ScreenAddRuleSelect
			return m, nil
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("tab"))) {
			// Toggle between NAT and Firewall rules view
			if m.ruleViewMode == "nat" {
				m.ruleViewMode = "firewall"
			} else {
				m.ruleViewMode = "nat"
			}
			m.rulesScrollOffset = 0 // Reset scroll when switching views
			return m, nil
		}

		// Scrolling
		visibleHeight := m.getRulesVisibleHeight()
		var listLen int
		if m.ruleViewMode == "nat" {
			listLen = len(m.natRules)
		} else {
			listLen = len(m.firewallRules)
		}

		if key.Matches(msg, keys.Down) || key.Matches(msg, key.NewBinding(key.WithKeys("j"))) {
			m.rulesScrollOffset++
			if m.rulesScrollOffset > listLen-visibleHeight {
				m.rulesScrollOffset = listLen - visibleHeight
			}
			if m.rulesScrollOffset < 0 {
				m.rulesScrollOffset = 0
			}
			return m, nil
		}
		if key.Matches(msg, keys.Up) || key.Matches(msg, key.NewBinding(key.WithKeys("k"))) {
			m.rulesScrollOffset--
			if m.rulesScrollOffset < 0 {
				m.rulesScrollOffset = 0
			}
			return m, nil
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("pgdown"))) {
			m.rulesScrollOffset += visibleHeight
			if m.rulesScrollOffset > listLen-visibleHeight {
				m.rulesScrollOffset = listLen - visibleHeight
			}
			if m.rulesScrollOffset < 0 {
				m.rulesScrollOffset = 0
			}
			return m, nil
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("pgup"))) {
			m.rulesScrollOffset -= visibleHeight
			if m.rulesScrollOffset < 0 {
				m.rulesScrollOffset = 0
			}
			return m, nil
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("home"))) {
			m.rulesScrollOffset = 0
			return m, nil
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("end"))) {
			m.rulesScrollOffset = listLen - visibleHeight
			if m.rulesScrollOffset < 0 {
				m.rulesScrollOffset = 0
			}
			return m, nil
		}

		if key.Matches(msg, keys.Back) {
			m.screen = ScreenMenu
			m.rulesScrollOffset = 0
			return m, nil
		}
	}

	return m, nil
}

// deleteFirstRule deletes the first visible rule
func (m *Model) deleteFirstRule() (tea.Model, tea.Cmd) {
	if m.ruleViewMode == "nat" && len(m.natRules) > 0 {
		rule := m.natRules[0]
		m.loadingMsg = fmt.Sprintf("Removing NAT rule %s...", rule.ID)
		m.screen = ScreenLoading

		return m, func() tea.Msg {
			if m.provider != nil {
				m.provider.RemoveNAT(m.ctx, rule.ID)
			}
			return successMsg{fmt.Sprintf("NAT rule %s removed", rule.ID)}
		}
	}

	if m.ruleViewMode == "firewall" && len(m.firewallRules) > 0 {
		rule := m.firewallRules[0]
		m.loadingMsg = fmt.Sprintf("Removing firewall rule %s...", rule.ID)
		m.screen = ScreenLoading

		return m, func() tea.Msg {
			if m.provider != nil {
				m.provider.ClosePort(m.ctx, rule.ID)
			}
			return successMsg{fmt.Sprintf("Firewall rule %s removed", rule.ID)}
		}
	}

	return m, nil
}

// viewListRules renders the rules list
func (m *Model) viewListRules() string {
	modeLabel := "NAT Rules"
	if m.ruleViewMode == "firewall" {
		modeLabel = "Firewall Rules"
	}

	title := styles.Title.Render("List Rules")
	subtitle := styles.Subtitle.Render(fmt.Sprintf("%s (%d NAT, %d Firewall)", modeLabel, len(m.natRules), len(m.firewallRules)))

	if m.ruleViewMode == "nat" {
		return m.viewNATRules(title, subtitle)
	}
	return m.viewFirewallRules(title, subtitle)
}

// viewNATRules renders NAT rules table
func (m *Model) viewNATRules(title, subtitle string) string {
	if len(m.natRules) == 0 {
		var msg string
		if len(m.firewallRules) > 0 {
			msg = fmt.Sprintf("No NAT rules found, but %d firewall rule(s) exist.\nPress TAB to view firewall rules.", len(m.firewallRules))
		} else {
			msg = "No NAT rules found. Press 'a' to add a rule."
		}
		content := styles.Info.Render(msg)
		help := styles.Help.Render("esc: back • a: add • r: refresh • tab: switch view")

		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			subtitle,
			"",
			styles.Panel.Render(content),
			"",
			help,
		)
	}

	visibleHeight := m.getRulesVisibleHeight()
	m.updateRulesScrollOffset()
	startIdx := m.rulesScrollOffset
	endIdx := startIdx + visibleHeight
	if endIdx > len(m.natRules) {
		endIdx = len(m.natRules)
	}

	var rows []string

	// Header
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.TableHeader.Width(20).Render("ID"),
		styles.TableHeader.Width(12).Render("Product"),
		styles.TableHeader.Width(10).Render("External"),
		styles.TableHeader.Width(22).Render("Internal"),
		styles.TableHeader.Width(6).Render("Proto"),
	)
	rows = append(rows, header)
	rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.BorderColor)).Render(
		strings.Repeat("─", 72),
	))

	// Show scroll indicators if needed
	if startIdx > 0 {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.MutedColor)).Render("↑ more above"))
	}

	// Visible rows only
	for i := startIdx; i < endIdx; i++ {
		rule := m.natRules[i]
		// Truncate ID if too long for display
		displayID := rule.ID
		if len(displayID) > 20 {
			displayID = displayID[:17] + "..."
		}
		// Show product or "system" if empty
		product := rule.Product
		if product == "" {
			product = "system"
		}
		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.TableCell.Width(20).Render(displayID),
			styles.TableCell.Width(12).Render(product),
			styles.TableCell.Width(10).Render(fmt.Sprintf("%d", rule.ExternalPort)),
			styles.TableCell.Width(22).Render(fmt.Sprintf("%s:%d", rule.InternalIP, rule.InternalPort)),
			styles.TableCell.Width(6).Render(string(rule.Proto)),
		)
		rows = append(rows, row)
	}

	// Show scroll indicators if needed
	if endIdx < len(m.natRules) {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.MutedColor)).Render("↓ more below"))
	}

	table := lipgloss.JoinVertical(lipgloss.Left, rows...)
	help := styles.Help.Render("↑/↓/j/k: scroll • pgup/pgdn: page • home/end: jump • esc: back • d: delete first • r: refresh • tab: firewall view")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		styles.Panel.Render(table),
		"",
		help,
	)
}

// viewFirewallRules renders firewall rules table
func (m *Model) viewFirewallRules(title, subtitle string) string {
	if len(m.firewallRules) == 0 {
		var msg string
		if len(m.natRules) > 0 {
			msg = fmt.Sprintf("No firewall rules found, but %d NAT rule(s) exist.\nPress TAB to view NAT rules.", len(m.natRules))
		} else {
			msg = "No firewall rules found. Press 'a' to add a rule."
		}
		content := styles.Info.Render(msg)
		help := styles.Help.Render("esc: back • a: add • r: refresh • tab: switch view")

		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			subtitle,
			"",
			styles.Panel.Render(content),
			"",
			help,
		)
	}

	visibleHeight := m.getRulesVisibleHeight()
	m.updateRulesScrollOffset()
	startIdx := m.rulesScrollOffset
	endIdx := startIdx + visibleHeight
	if endIdx > len(m.firewallRules) {
		endIdx = len(m.firewallRules)
	}

	var rows []string

	// Header
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.TableHeader.Width(20).Render("ID"),
		styles.TableHeader.Width(12).Render("Type"),
		styles.TableHeader.Width(6).Render("Port"),
		styles.TableHeader.Width(6).Render("Proto"),
		styles.TableHeader.Width(16).Render("Source"),
		styles.TableHeader.Width(12).Render("Product"),
	)
	rows = append(rows, header)
	rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.BorderColor)).Render(
		strings.Repeat("─", 74),
	))

	// Show scroll indicators if needed
	if startIdx > 0 {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.MutedColor)).Render("↑ more above"))
	}

	// Visible rows only
	for i := startIdx; i < endIdx; i++ {
		rule := m.firewallRules[i]
		source := "any"
		if rule.SourceIP != "" {
			source = rule.SourceIP
		}

		typeLabel := string(rule.Type)
		if typeLabel == "" {
			typeLabel = "port"
		}

		portStr := fmt.Sprintf("%d", rule.Port)
		if rule.Port == 0 {
			if rule.Type == models.RuleTypeTrustIP {
				portStr = "all"
			} else {
				portStr = "-"
			}
		}

		// Truncate ID if too long for display
		displayID := rule.ID
		if len(displayID) > 20 {
			displayID = displayID[:17] + "..."
		}

		// Show product or "system" if empty
		product := rule.Product
		if product == "" {
			product = "system"
		}

		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.TableCell.Width(20).Render(displayID),
			styles.TableCell.Width(12).Render(typeLabel),
			styles.TableCell.Width(6).Render(portStr),
			styles.TableCell.Width(6).Render(string(rule.Protocol)),
			styles.TableCell.Width(16).Render(source),
			styles.TableCell.Width(12).Render(product),
		)
		rows = append(rows, row)
	}

	// Show scroll indicators if needed
	if endIdx < len(m.firewallRules) {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.MutedColor)).Render("↓ more below"))
	}

	table := lipgloss.JoinVertical(lipgloss.Left, rows...)
	help := styles.Help.Render("↑/↓/j/k: scroll • pgup/pgdn: page • home/end: jump • esc: back • d: delete first • r: refresh • tab: NAT view")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		styles.Panel.Render(table),
		"",
		help,
	)
}
