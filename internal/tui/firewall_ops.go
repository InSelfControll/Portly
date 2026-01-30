package tui

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// firewallStatus holds firewall state
type firewallStatus struct {
	running bool
	enabled bool
}

// getFirewallStatus checks if firewall is running
func (m *Model) getFirewallStatus() firewallStatus {
	if m.provider == nil {
		return firewallStatus{false, false}
	}

	var running bool
	switch m.provider.Name() {
	case "firewalld":
		cmd := exec.Command("systemctl", "is-active", "firewalld")
		running = cmd.Run() == nil
	case "nftables":
		cmd := exec.Command("systemctl", "is-active", "nftables")
		running = cmd.Run() == nil
	case "pf":
		cmd := exec.Command("/sbin/pfctl", "-s", "info")
		out, _ := cmd.Output()
		running = len(out) > 0
	}

	return firewallStatus{running, running}
}

// toggleFirewall starts or stops the firewall
func (m *Model) toggleFirewall(start bool) tea.Cmd {
	return func() tea.Msg {
		if m.provider == nil {
			return errMsg{fmt.Errorf("no firewall provider available")}
		}

		action := "stop"
		if start {
			action = "start"
		}

		var cmd *exec.Cmd
		switch m.provider.Name() {
		case "firewalld":
			cmd = exec.Command("systemctl", action, "firewalld")
		case "nftables":
			cmd = exec.Command("systemctl", action, "nftables")
		case "pf":
			if start {
				cmd = exec.Command("/sbin/pfctl", "-e")
			} else {
				cmd = exec.Command("/sbin/pfctl", "-d")
			}
		}

		if cmd != nil {
			if out, err := cmd.CombinedOutput(); err != nil {
				return errMsg{fmt.Errorf("failed to %s firewall: %v (output: %s)", action, err, string(out))}
			}
		}

		status := "stopped"
		if start {
			status = "started"
		}
		return successMsg{fmt.Sprintf("Firewall %s successfully", status)}
	}
}

// installFirewall installs the appropriate firewall
func (m *Model) installFirewall() tea.Cmd {
	return func() tea.Msg {
		pm := m.osInfo.PackageManager()
		var pkg string

		switch m.osInfo.Family {
		case "rhel":
			pkg = "firewalld"
		case "debian":
			pkg = "nftables"
		default:
			return errMsg{fmt.Errorf("automatic installation not supported on %s", m.osInfo.Family)}
		}

		var cmd *exec.Cmd
		switch pm {
		case "dnf":
			cmd = exec.Command("dnf", "install", "-y", pkg)
		case "yum":
			cmd = exec.Command("yum", "install", "-y", pkg)
		case "apt":
			cmd = exec.Command("apt-get", "install", "-y", pkg)
		}

		if cmd != nil {
			if out, err := cmd.CombinedOutput(); err != nil {
				return errMsg{fmt.Errorf("failed to install %s: %v (output: %s)", pkg, err, string(out))}
			}
		}

		// Enable and start the service
		enableCmd := exec.Command("systemctl", "enable", "--now", pkg)
		enableCmd.Run()

		return successMsg{fmt.Sprintf("%s installed and started", pkg)}
	}
}
