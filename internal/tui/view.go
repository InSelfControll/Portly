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
	case ScreenAddRule:
		content = m.viewAddRule()
	case ScreenListRules:
		content = m.viewListRules()
	case ScreenFirewall:
		content = m.viewFirewall()
	case ScreenSecurity:
		content = m.viewSecurity()
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
