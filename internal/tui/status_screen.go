package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/installer"
	"github.com/orchestrator/unified-firewall/internal/platform"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

func (m *Model) updateStatus(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg.(tea.KeyMsg), keys.Back) {
			m.screen = ScreenMenu
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) viewStatus() string {
	title := styles.Title.Render("System Status")

	osContent := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("OS:      %s", m.osInfo.Distribution),
		fmt.Sprintf("Family:  %s", m.osInfo.Family),
		fmt.Sprintf("Version: %s", m.osInfo.Version),
	)
	osPanel := styles.Panel.Width(35).Render(osContent)

	providerName := "Not available"
	if m.provider != nil {
		providerName = m.provider.Name()
	}

	rootStatus := styles.Error.Render("No")
	if platform.IsRoot() {
		rootStatus = styles.Success.Render("Yes")
	}

	providerContent := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("Provider: %s", providerName),
		fmt.Sprintf("Root:     %s", rootStatus),
	)
	providerPanel := styles.Panel.Width(35).Render(providerContent)

	var productLines []string
	products := installer.GetSupportedProducts()
	for _, name := range products {
		_, installed := platform.IsProductInstalled(name)
		if installed {
			productLines = append(productLines, styles.Success.Render("✓ "+name))
		} else {
			productLines = append(productLines, styles.MutedColor+"✗ "+name)
		}
	}
	productsContent := lipgloss.JoinVertical(lipgloss.Left, productLines...)
	productsPanel := styles.Panel.Width(35).Render(productsContent)

	stateContent := "State not available"
	if m.stateMgr != nil {
		allRules := m.stateMgr.ListRules()
		activeRules := m.stateMgr.ListActiveRules()
		stateContent = fmt.Sprintf("Total: %d | Active: %d", len(allRules), len(activeRules))
	}
	statePanel := styles.Panel.Width(35).Render(stateContent)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, osPanel, "  ", providerPanel)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, productsPanel, "  ", statePanel)

	help := styles.Help.Render("esc: back")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		topRow,
		"",
		bottomRow,
		"",
		help,
	)
}
