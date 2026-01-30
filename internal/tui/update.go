package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menuList.SetSize(msg.Width-4, msg.Height-8)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Back) && m.screen != ScreenMenu:
			m.screen = ScreenMenu
			m.lastError = nil
			return m, nil
		}
	}

	switch m.screen {
	case ScreenMenu:
		return m.updateMenu(msg)
	case ScreenAddRule:
		return m.updateAddRule(msg)
	case ScreenListRules:
		return m.updateListRules(msg)
	case ScreenFirewall:
		return m.updateFirewall(msg)
	case ScreenSecurity:
		return m.updateSecurity(msg)
	case ScreenStatus:
		return m.updateStatus(msg)
	case ScreenCheck:
		return m.updateCheck(msg)
	case ScreenLoading:
		return m.updateLoading(msg)
	}

	return m, nil
}
