package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

// View renders the UI
func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var content string
	switch m.screen {
	case ScreenMenu:
		content = m.viewMenu()
	case ScreenAddRuleSelect:
		content = m.viewAddRuleSelect()
	case ScreenAddNATRule, ScreenOpenPort, ScreenOpenIPPort, ScreenOpenIP:
		content = m.viewAddRule()
	case ScreenListRules:
		content = m.viewListRules()
	case ScreenFirewall:
		content = m.viewFirewall()
	case ScreenSecurity:
		content = m.viewSecurity()
	case ScreenSecurityRules:
		content = m.viewSecurityRules()
	case ScreenStatus:
		content = m.viewStatus()
	case ScreenCheck:
		content = m.viewCheck()
	case ScreenLoading:
		content = m.viewLoading()
	case ScreenError:
		content = m.viewError()
	case ScreenSuccess:
		content = m.viewSuccess()
	}

	statusBar := m.renderStatusBar()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		statusBar,
	)
}

func (m *Model) renderStatusBar() string {
	providerName := "none"
	if m.provider != nil {
		providerName = m.provider.Name()
	}

	status := fmt.Sprintf(" Portly | %s | %s | Provider: %s ",
		m.osInfo.Distribution,
		m.osInfo.Family,
		providerName,
	)

	return styles.StatusBar.Width(m.width).Render(status)
}

func (m *Model) viewAddRuleSelect() string {
	title := styles.Title.Render("Add New Rule")
	subtitle := styles.Subtitle.Render("Select the type of rule to create")

	menu := m.ruleSubMenuList.View()

	help := styles.Help.Render("↑/↓: navigate • enter: select • esc: back")

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
