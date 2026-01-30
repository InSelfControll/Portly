package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/platform"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

func (m *Model) updateCheck(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg.(tea.KeyMsg), keys.Back) {
			m.screen = ScreenMenu
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) viewCheck() string {
	title := styles.Title.Render("Configuration Check")

	var checks []string

	osCheck := fmt.Sprintf("[✓] OS: %s (%s)", m.osInfo.Distribution, m.osInfo.Family)
	checks = append(checks, styles.Success.Render(osCheck))

	if platform.IsRoot() {
		checks = append(checks, styles.Success.Render("[✓] Running as root"))
	} else {
		checks = append(checks, styles.Error.Render("[✗] Not running as root"))
	}

	if m.provider != nil {
		checks = append(checks, styles.Success.Render(fmt.Sprintf("[✓] Provider: %s", m.provider.Name())))
	} else {
		checks = append(checks, styles.Error.Render("[✗] No firewall provider"))
	}

	if m.stateMgr != nil {
		checks = append(checks, styles.Success.Render("[✓] State manager ready"))
	} else {
		checks = append(checks, styles.Warning.Render("[!] State manager unavailable"))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, checks...)
	help := styles.Help.Render("esc: back")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Panel.Render(content),
		"",
		help,
	)
}
