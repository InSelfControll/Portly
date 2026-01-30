package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

// ScreenSecurity is the security management screen

// updateSecurity handles security screen updates
func (m *Model) updateSecurity(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Back) {
			m.screen = ScreenMenu
			return m, nil
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("1"))) && m.osInfo.IsRHEL() {
			m.loadingMsg = "Setting SELinux to enforcing..."
			m.screen = ScreenLoading
			return m, m.toggleSELinux()
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("2"))) && m.osInfo.IsRHEL() {
			m.loadingMsg = "Setting SELinux to permissive..."
			m.screen = ScreenLoading
			return m, m.setSELinuxPermissive()
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("a"))) && m.osInfo.IsDebian() {
			m.loadingMsg = "Toggling AppArmor..."
			m.screen = ScreenLoading
			return m, m.toggleAppArmor()
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("v"))) && (m.osInfo.IsRHEL() || m.osInfo.IsDebian()) {
			m.screen = ScreenSecurityRules
			return m, m.loadSecurityRulesData()
		}
	}
	return m, nil
}

// viewSecurity renders the security management screen
func (m *Model) viewSecurity() string {
	title := styles.Title.Render("🔒 Security Management")

	var content string
	if m.osInfo.IsRHEL() {
		content = m.viewSELinux()
	} else if m.osInfo.IsDebian() {
		content = m.viewAppArmor()
	} else if m.osInfo.IsDarwin() {
		content = m.viewMacOSSecurity()
	} else {
		content = styles.Warning.Render("Security management not available on this platform")
	}

	var helpText string
	if m.osInfo.IsRHEL() || m.osInfo.IsDebian() {
		helpText = "esc: back • v: view rules"
	} else {
		helpText = "esc: back"
	}
	help := styles.Help.Render(helpText)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.Panel.Render(content),
		"",
		help,
	)
}

// viewSELinux shows SELinux management
func (m *Model) viewSELinux() string {
	status := m.getSELinuxStatus()

	modeStr := styles.Info.Render(status.mode)
	if status.mode == "enforcing" {
		modeStr = styles.Success.Render(status.mode)
	} else if status.mode == "permissive" {
		modeStr = styles.Warning.Render(status.mode)
	} else if status.mode == "disabled" {
		modeStr = styles.Error.Render(status.mode)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"SELinux Status:",
		fmt.Sprintf("  Mode: %s", modeStr),
		fmt.Sprintf("  Policy: %s", status.policy),
		"",
		styles.Button.Render("[1] Set Enforcing"),
		styles.Button.Render("[2] Set Permissive"),
	)
}

// viewAppArmor shows AppArmor management
func (m *Model) viewAppArmor() string {
	status := m.getAppArmorStatus()

	statusStr := styles.Error.Render("Not Loaded")
	if status.loaded {
		statusStr = styles.Success.Render("Loaded")
	}

	modeStr := styles.Error.Render(status.mode)
	if status.mode == "enforce" {
		modeStr = styles.Success.Render(status.mode)
	} else if status.mode == "complain" {
		modeStr = styles.Warning.Render(status.mode)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"AppArmor Status:",
		fmt.Sprintf("  Status: %s", statusStr),
		fmt.Sprintf("  Mode:   %s", modeStr),
		"",
		styles.Button.Render("[a] Toggle AppArmor"),
	)
}

// viewMacOSSecurity shows macOS SIP info
func (m *Model) viewMacOSSecurity() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"macOS Security:",
		"  System Integrity Protection (SIP)",
		"  SIP is managed by macOS and cannot be",
		"  disabled from within the running system.",
		"",
		styles.Info.Render("Boot into Recovery Mode to modify SIP"),
	)
}
