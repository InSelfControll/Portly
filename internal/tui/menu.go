package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

// updateMenu handles menu updates
func (m *Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Enter) {
			item, ok := m.menuList.SelectedItem().(menuItem)
			if ok {
				if item.screen == -1 {
					return m, tea.Quit
				}
				m.screen = item.screen
				
				// Initialize data for the selected screen
				switch item.screen {
				case ScreenListRules:
					return m, m.loadRules()
				case ScreenStatus:
					return m, nil
				case ScreenCheck:
					return m, nil
				}
			}
		}
	}

	var cmd tea.Cmd
	m.menuList, cmd = m.menuList.Update(msg)
	return m, cmd
}

// viewMenu renders the menu
func (m *Model) viewMenu() string {
	title := styles.Title.Render("🔥 Unified Firewall Orchestrator")
	subtitle := styles.Subtitle.Render("Manage NAT rules and security policies")
	
	menu := m.menuList.View()
	
	help := styles.Help.Render("↑/↓: navigate • enter: select • q: quit")
	
	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		"",
		styles.Panel.Render(menu),
		"",
		help,
	)
}
