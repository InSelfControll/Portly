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

// loadRules loads rules from state
func (m *Model) loadRules() tea.Cmd {
	return func() tea.Msg {
		if m.stateMgr == nil {
			return rulesMsg{[]models.AppliedRule{}}
		}
		return rulesMsg{m.stateMgr.ListActiveRules()}
	}
}

// rulesMsg carries rules data
type rulesMsg struct {
	rules []models.AppliedRule
}

// updateListRules handles list rules updates
func (m *Model) updateListRules(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case rulesMsg:
		m.rules = msg.rules
		return m, nil
		
	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("d", "delete"))) && len(m.rules) > 0 {
			// Delete selected rule
			if len(m.rules) > 0 {
				return m.deleteRule(0)
			}
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("r"))) {
			return m, m.loadRules()
		}
	}
	
	return m, nil
}

// deleteRule deletes a rule
func (m *Model) deleteRule(index int) (tea.Model, tea.Cmd) {
	if index >= len(m.rules) {
		return m, nil
	}
	
	rule := m.rules[index]
	m.loadingMsg = fmt.Sprintf("Removing rule %s...", rule.ID)
	m.screen = ScreenLoading
	
	return m, func() tea.Msg {
		if m.provider != nil {
			m.provider.RemoveNAT(m.ctx, rule.ID)
		}
		if m.stateMgr != nil {
			m.stateMgr.RemoveRule(rule.ID)
		}
		return successMsg{fmt.Sprintf("Rule %s removed", rule.ID)}
	}
}

// viewListRules renders the rules list
func (m *Model) viewListRules() string {
	title := styles.Title.Render("NAT Rules")
	subtitle := styles.Subtitle.Render(fmt.Sprintf("%d active rules", len(m.rules)))
	
	if len(m.rules) == 0 {
		content := styles.Info.Render("No active rules found. Press 'a' to add a rule.")
		help := styles.Help.Render("esc: back • a: add rule")
		
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
	
	// Render table
	var rows []string
	
	// Header
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.TableHeader.Width(10).Render("ID"),
		styles.TableHeader.Width(12).Render("Product"),
		styles.TableHeader.Width(10).Render("External"),
		styles.TableHeader.Width(18).Render("Internal"),
		styles.TableHeader.Width(8).Render("Proto"),
		styles.TableHeader.Width(10).Render("Status"),
	)
	rows = append(rows, header)
	rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.BorderColor)).Render(
		strings.Repeat("─", 70),
	))
	
	// Rows
	for _, rule := range m.rules {
		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.TableCell.Width(10).Render(rule.ID),
			styles.TableCell.Width(12).Render(rule.Product),
			styles.TableCell.Width(10).Render(fmt.Sprintf("%d", rule.ExternalPort)),
			styles.TableCell.Width(18).Render(fmt.Sprintf("%s:%d", rule.InternalIP, rule.InternalPort)),
			styles.TableCell.Width(8).Render(string(rule.Proto)),
			styles.TableCell.Width(10).Render(string(rule.Status)),
		)
		rows = append(rows, row)
	}
	
	table := lipgloss.JoinVertical(lipgloss.Left, rows...)
	help := styles.Help.Render("esc: back • d: delete first • r: refresh")
	
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
