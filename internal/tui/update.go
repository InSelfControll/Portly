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
		m.ruleSubMenuList.SetSize(msg.Width-4, msg.Height-8)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Back):
			switch m.screen {
			case ScreenAddNATRule, ScreenOpenPort, ScreenOpenIPPort, ScreenOpenIP:
				m.screen = ScreenAddRuleSelect // Go back to sub-menu
				return m, nil
			case ScreenAddRuleSelect, ScreenListRules, ScreenStatus, ScreenCheck, ScreenFirewall, ScreenSecurity:
				m.screen = ScreenMenu // Go back to main menu
				m.lastError = nil
				return m, nil
			case ScreenError, ScreenSuccess:
				// Return to menu on error/success screens
				m.screen = ScreenMenu
				m.lastError = nil
				m.successMsg = ""
				return m, nil
			}
		}
	}

	switch m.screen {
	case ScreenMenu:
		return m.updateMenu(msg)
	case ScreenAddRuleSelect:
		return m.updateAddRuleSelect(msg)
	case ScreenAddNATRule, ScreenOpenPort, ScreenOpenIPPort, ScreenOpenIP:
		return m.updateAddRule(msg)
	case ScreenListRules:
		return m.updateListRules(msg)
	case ScreenFirewall:
		return m.updateFirewall(msg)
	case ScreenSecurity:
		return m.updateSecurity(msg)
	case ScreenSecurityRules:
		return m.updateSecurityRules(msg)
	case ScreenStatus:
		return m.updateStatus(msg)
	case ScreenCheck:
		return m.updateCheck(msg)
	case ScreenLoading:
		return m.updateLoading(msg)
	case ScreenError, ScreenSuccess:
		// Handle any key to dismiss error/success screens
		if _, ok := msg.(tea.KeyMsg); ok {
			m.screen = ScreenMenu
			m.lastError = nil
			m.successMsg = ""
		}
		return m, nil
	}

	return m, nil
}
