package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

// ScreenFirewall is the firewall management screen

// updateFirewall handles firewall screen updates
func (m *Model) updateFirewall(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Back) {
			m.screen = ScreenMenu
			return m, nil
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("1"))) {
			m.loadingMsg = "Starting firewall..."
			m.screen = ScreenLoading
			return m, m.toggleFirewall(true)
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("2"))) {
			m.loadingMsg = "Stopping firewall..."
			m.screen = ScreenLoading
			return m, m.toggleFirewall(false)
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("i"))) {
			m.loadingMsg = "Installing firewall..."
			m.screen = ScreenLoading
			return m, m.installFirewall()
		}
	}
	return m, nil
}

// viewFirewall renders the firewall management screen
func (m *Model) viewFirewall() string {
	title := styles.Title.Render("🔥 Firewall Management")
	status := m.getFirewallStatus()

	var content string
	if m.provider == nil {
		content = m.viewNoFirewall()
	} else {
		content = m.viewFirewallStatus(status)
	}

	help := styles.Help.Render("1: start • 2: stop • i: install • esc: back")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Panel.Render(content),
		"",
		help,
	)
}

// viewNoFirewall shows when no firewall is detected
func (m *Model) viewNoFirewall() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		styles.Error.Render("⚠️  No firewall provider detected"),
		"",
		"Your system may not have a supported firewall installed.",
		"Supported: firewalld (RHEL), nftables (Ubuntu), pf (macOS)",
		"",
		styles.Button.Render("[i] Install Firewall"),
	)
}

// viewFirewallStatus shows firewall status
func (m *Model) viewFirewallStatus(status firewallStatus) string {
	providerName := m.provider.Name()
	statusStr := styles.Success.Render("✓ Running")
	if !status.running {
		statusStr = styles.Error.Render("✗ Stopped")
	}

	var toggleBtn string
	if status.running {
		toggleBtn = styles.Button.Render("[2] Stop Firewall")
	} else {
		toggleBtn = styles.Button.Render("[1] Start Firewall")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("Provider: %s", providerName),
		fmt.Sprintf("Status:   %s", statusStr),
		fmt.Sprintf("Enabled:  %v", status.enabled),
		"",
		toggleBtn,
	)
}
