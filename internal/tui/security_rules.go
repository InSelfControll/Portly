package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

// securityRulesMsg carries security configuration data
type securityRulesMsg struct {
	seLinuxBooleans  []SELinuxBoolean
	appArmorProfiles []AppArmorProfile
}

// SELinuxBoolean represents an SELinux boolean setting
type SELinuxBoolean struct {
	Name        string
	State       string // on, off
	Default     string
	Description string
}

// AppArmorProfile represents an AppArmor profile status
type AppArmorProfile struct {
	Name string
	Mode string // enforce, complain, disabled
}

// loadSecurityRules loads security configuration
type loadSecurityRules struct{}

func (m *Model) loadSecurityRulesData() tea.Cmd {
	return func() tea.Msg {
		msg := securityRulesMsg{}

		if m.osInfo.IsRHEL() {
			msg.seLinuxBooleans = m.getSELinuxBooleans()
		} else if m.osInfo.IsDebian() {
			msg.appArmorProfiles = m.getAppArmorProfiles()
		}

		return msg
	}
}

// getSELinuxBooleans gets relevant SELinux booleans
func (m *Model) getSELinuxBooleans() []SELinuxBoolean {
	var booleans []SELinuxBoolean

	// List of container/network related booleans we care about
	relevantBooleans := []string{
		"container_manage_cgroup",
		"container_use_cephfs",
		"domain_can_mmap_files",
		"httpd_can_network_connect",
		"httpd_can_network_relay",
		"nis_enabled",
		"daemons_enable_cluster_mode",
	}

	for _, name := range relevantBooleans {
		cmd := exec.Command("getsebool", name)
		out, err := cmd.Output()
		if err != nil {
			continue
		}

		// Parse output like "container_manage_cgroup --> on"
		parts := strings.Split(string(out), "-->")
		if len(parts) == 2 {
			state := strings.TrimSpace(parts[1])
			booleans = append(booleans, SELinuxBoolean{
				Name:    name,
				State:   state,
				Default: state,
			})
		}
	}

	return booleans
}

// getAppArmorProfiles gets AppArmor profile status
func (m *Model) getAppArmorProfiles() []AppArmorProfile {
	var profiles []AppArmorProfile

	// Try aa-status first
	cmd := exec.Command("aa-status")
	out, err := cmd.Output()
	if err != nil {
		return profiles
	}

	lines := strings.Split(string(out), "\n")
	inComplain := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "profiles are in enforce mode") {
			inComplain = false
			continue
		}
		if strings.HasPrefix(line, "profiles are in complain mode") {
			inComplain = true
			continue
		}
		if strings.HasPrefix(line, "processes") || strings.HasPrefix(line, "0 processes") {
			inComplain = false
			continue
		}

		// Skip empty lines and headers
		if line == "" || strings.HasSuffix(line, ":") {
			continue
		}

		// Parse profile name (remove trailing space and (enforce) or (complain))
		profileName := strings.TrimSpace(line)
		profileName = strings.TrimSuffix(profileName, " (enforce)")
		profileName = strings.TrimSuffix(profileName, " (complain)")

		if profileName != "" && !strings.Contains(profileName, " ") {
			mode := "enforce"
			if inComplain {
				mode = "complain"
			}
			profiles = append(profiles, AppArmorProfile{
				Name: profileName,
				Mode: mode,
			})
		}
	}

	return profiles
}

// toggleSELinuxBoolean toggles an SELinux boolean on/off
func (m *Model) toggleSELinuxBoolean(index int) tea.Cmd {
	return func() tea.Msg {
		if index >= len(m.seLinuxBooleans) {
			return errMsg{fmt.Errorf("invalid boolean index")}
		}

		boolean := m.seLinuxBooleans[index]
		newState := "on"
		if boolean.State == "on" {
			newState = "off"
		}

		cmd := exec.Command("setsebool", "-P", boolean.Name, newState)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return errMsg{fmt.Errorf("failed to toggle %s: %v (output: %s)", boolean.Name, err, string(output))}
		}

		return successMsg{fmt.Sprintf("%s set to %s", boolean.Name, newState)}
	}
}

// toggleAppArmorProfile toggles an AppArmor profile between enforce and complain
func (m *Model) toggleAppArmorProfile(index int) tea.Cmd {
	return func() tea.Msg {
		if index >= len(m.appArmorProfiles) {
			return errMsg{fmt.Errorf("invalid profile index")}
		}

		profile := m.appArmorProfiles[index]
		newMode := "enforce"
		if profile.Mode == "enforce" {
			newMode = "complain"
		}

		// Use aa-complain or aa-enforce
		var cmd *exec.Cmd
		if newMode == "enforce" {
			cmd = exec.Command("aa-enforce", profile.Name)
		} else {
			cmd = exec.Command("aa-complain", profile.Name)
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			return errMsg{fmt.Errorf("failed to set %s to %s: %v (output: %s)", profile.Name, newMode, err, string(output))}
		}

		return successMsg{fmt.Sprintf("%s set to %s mode", profile.Name, newMode)}
	}
}

// disableAppArmorProfile disables an AppArmor profile
func (m *Model) disableAppArmorProfile(index int) tea.Cmd {
	return func() tea.Msg {
		if index >= len(m.appArmorProfiles) {
			return errMsg{fmt.Errorf("invalid profile index")}
		}

		profile := m.appArmorProfiles[index]

		cmd := exec.Command("aa-disable", profile.Name)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return errMsg{fmt.Errorf("failed to disable %s: %v (output: %s)", profile.Name, err, string(output))}
		}

		return successMsg{fmt.Sprintf("%s disabled", profile.Name)}
	}
}

// updateScrollOffset ensures the selection is visible within the scroll window
func (m *Model) updateSecurityScrollOffset(visibleHeight int) {
	if m.securitySelectionIdx < m.securityScrollOffset {
		// Selection is above the visible area
		m.securityScrollOffset = m.securitySelectionIdx
	} else if m.securitySelectionIdx >= m.securityScrollOffset+visibleHeight {
		// Selection is below the visible area
		m.securityScrollOffset = m.securitySelectionIdx - visibleHeight + 1
	}
}

// getVisibleHeight returns the number of rows that can fit on screen
func (m *Model) getSecurityVisibleHeight() int {
	// Account for: title, subtitle, empty line, header, separator, empty line, help, status bar
	// Approximate: header (2) + margins (4) + help (1) + status (1) = ~8 lines
	availableHeight := m.height - 12
	if availableHeight < 5 {
		return 5 // Minimum visible rows
	}
	return availableHeight
}

// updateSecurityRules handles security rules screen updates
func (m *Model) updateSecurityRules(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case securityRulesMsg:
		m.seLinuxBooleans = msg.seLinuxBooleans
		m.appArmorProfiles = msg.appArmorProfiles
		// Reset selection when loading new data
		m.securitySelectionIdx = 0
		m.securityScrollOffset = 0
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, keys.Back) {
			m.screen = ScreenSecurity
			m.securitySelectionIdx = 0
			m.securityScrollOffset = 0
			return m, nil
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("r"))) {
			return m, m.loadSecurityRulesData()
		}

		// Get visible height for scrolling
		visibleHeight := m.getSecurityVisibleHeight()
		var listLen int
		if m.osInfo.IsRHEL() {
			listLen = len(m.seLinuxBooleans)
		} else {
			listLen = len(m.appArmorProfiles)
		}

		// Navigation
		if key.Matches(msg, keys.Down) {
			if listLen > 0 {
				m.securitySelectionIdx = (m.securitySelectionIdx + 1) % listLen
				m.updateSecurityScrollOffset(visibleHeight)
			}
			return m, nil
		}
		if key.Matches(msg, keys.Up) {
			if listLen > 0 {
				m.securitySelectionIdx--
				if m.securitySelectionIdx < 0 {
					m.securitySelectionIdx = listLen - 1
				}
				m.updateSecurityScrollOffset(visibleHeight)
			}
			return m, nil
		}

		// Page navigation
		if key.Matches(msg, key.NewBinding(key.WithKeys("pgdown"))) {
			if listLen > 0 {
				m.securitySelectionIdx += visibleHeight
				if m.securitySelectionIdx >= listLen {
					m.securitySelectionIdx = listLen - 1
				}
				m.updateSecurityScrollOffset(visibleHeight)
			}
			return m, nil
		}
		if key.Matches(msg, key.NewBinding(key.WithKeys("pgup"))) {
			if listLen > 0 {
				m.securitySelectionIdx -= visibleHeight
				if m.securitySelectionIdx < 0 {
					m.securitySelectionIdx = 0
				}
				m.updateSecurityScrollOffset(visibleHeight)
			}
			return m, nil
		}

		// Toggle action (Enter or Space)
		if key.Matches(msg, keys.Enter) || msg.String() == " " {
			if m.osInfo.IsRHEL() && len(m.seLinuxBooleans) > 0 {
				m.loadingMsg = fmt.Sprintf("Toggling %s...", m.seLinuxBooleans[m.securitySelectionIdx].Name)
				m.screen = ScreenLoading
				return m, m.toggleSELinuxBoolean(m.securitySelectionIdx)
			} else if m.osInfo.IsDebian() && len(m.appArmorProfiles) > 0 {
				m.loadingMsg = fmt.Sprintf("Toggling %s...", m.appArmorProfiles[m.securitySelectionIdx].Name)
				m.screen = ScreenLoading
				return m, m.toggleAppArmorProfile(m.securitySelectionIdx)
			}
		}

		// Disable action (only for AppArmor - 'd' key)
		if msg.String() == "d" && m.osInfo.IsDebian() && len(m.appArmorProfiles) > 0 {
			m.loadingMsg = fmt.Sprintf("Disabling %s...", m.appArmorProfiles[m.securitySelectionIdx].Name)
			m.screen = ScreenLoading
			return m, m.disableAppArmorProfile(m.securitySelectionIdx)
		}
	}

	return m, nil
}

// viewSecurityRules renders the security rules screen
func (m *Model) viewSecurityRules() string {
	title := styles.Title.Render("🔒 Security Rules")

	if m.osInfo.IsRHEL() {
		return m.viewSELinuxRules(title)
	} else if m.osInfo.IsDebian() {
		return m.viewAppArmorRules(title)
	}

	// Not supported
	subtitle := styles.Subtitle.Render("Security rules view not available on this platform")
	content := styles.Warning.Render("SELinux/AppArmor rules management is only available on RHEL/Debian based systems")
	help := styles.Help.Render("esc: back")

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

// viewSELinuxRules renders SELinux booleans
func (m *Model) viewSELinuxRules(title string) string {
	subtitle := styles.Subtitle.Render(fmt.Sprintf("SELinux Booleans (%d configured)", len(m.seLinuxBooleans)))

	if len(m.seLinuxBooleans) == 0 {
		content := styles.Info.Render("No SELinux booleans found or SELinux is disabled.")
		help := styles.Help.Render("esc: back • r: refresh")

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

	visibleHeight := m.getSecurityVisibleHeight()
	startIdx := m.securityScrollOffset
	endIdx := startIdx + visibleHeight
	if endIdx > len(m.seLinuxBooleans) {
		endIdx = len(m.seLinuxBooleans)
	}

	var rows []string

	// Header
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.TableHeader.Width(5).Render(""),
		styles.TableHeader.Width(35).Render("Boolean"),
		styles.TableHeader.Width(10).Render("State"),
	)
	rows = append(rows, header)
	rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.BorderColor)).Render(
		strings.Repeat("─", 55),
	))

	// Show scroll indicators if needed
	if startIdx > 0 {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.MutedColor)).Render("↑ more above"))
	}

	// Visible rows only
	for i := startIdx; i < endIdx; i++ {
		b := m.seLinuxBooleans[i]
		stateStyle := styles.Error.Render
		if b.State == "on" {
			stateStyle = styles.Success.Render
		}

		// Selection indicator
		selector := "  "
		if i == m.securitySelectionIdx {
			selector = styles.Success.Render("> ")
		}

		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.TableCell.Width(5).Render(selector),
			styles.TableCell.Width(35).Render(b.Name),
			styles.TableCell.Width(10).Render(stateStyle(b.State)),
		)
		rows = append(rows, row)
	}

	// Show scroll indicators if needed
	if endIdx < len(m.seLinuxBooleans) {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.MutedColor)).Render("↓ more below"))
	}

	table := lipgloss.JoinVertical(lipgloss.Left, rows...)
	help := styles.Help.Render("↑/↓: navigate • pgup/pgdn: page • enter/space: toggle • r: refresh • esc: back")

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

// viewAppArmorRules renders AppArmor profiles
func (m *Model) viewAppArmorRules(title string) string {
	subtitle := styles.Subtitle.Render(fmt.Sprintf("AppArmor Profiles (%d loaded)", len(m.appArmorProfiles)))

	if len(m.appArmorProfiles) == 0 {
		content := styles.Info.Render("No AppArmor profiles found or AppArmor is disabled.")
		help := styles.Help.Render("esc: back • r: refresh")

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

	visibleHeight := m.getSecurityVisibleHeight()
	startIdx := m.securityScrollOffset
	endIdx := startIdx + visibleHeight
	if endIdx > len(m.appArmorProfiles) {
		endIdx = len(m.appArmorProfiles)
	}

	var rows []string

	// Header
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.TableHeader.Width(5).Render(""),
		styles.TableHeader.Width(40).Render("Profile"),
		styles.TableHeader.Width(12).Render("Mode"),
	)
	rows = append(rows, header)
	rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.BorderColor)).Render(
		strings.Repeat("─", 60),
	))

	// Show scroll indicators if needed
	if startIdx > 0 {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.MutedColor)).Render("↑ more above"))
	}

	// Visible rows only
	for i := startIdx; i < endIdx; i++ {
		p := m.appArmorProfiles[i]
		modeStyle := styles.Success.Render
		if p.Mode == "complain" {
			modeStyle = styles.Warning.Render
		}

		// Selection indicator
		selector := "  "
		if i == m.securitySelectionIdx {
			selector = styles.Success.Render("> ")
		}

		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.TableCell.Width(5).Render(selector),
			styles.TableCell.Width(40).Render(p.Name),
			styles.TableCell.Width(12).Render(modeStyle(p.Mode)),
		)
		rows = append(rows, row)
	}

	// Show scroll indicators if needed
	if endIdx < len(m.appArmorProfiles) {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(styles.MutedColor)).Render("↓ more below"))
	}

	table := lipgloss.JoinVertical(lipgloss.Left, rows...)
	help := styles.Help.Render("↑/↓: navigate • pgup/pgdn: page • enter/space: toggle • d: disable • r: refresh • esc: back")

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
