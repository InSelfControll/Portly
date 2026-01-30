package tui

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// selinuxStatus holds SELinux state
type selinuxStatus struct {
	mode   string
	policy string
}

// getSELinuxStatus checks SELinux status
func (m *Model) getSELinuxStatus() selinuxStatus {
	cmd := exec.Command("sestatus")
	out, err := cmd.Output()
	if err != nil {
		return selinuxStatus{"unknown", "unknown"}
	}

	status := selinuxStatus{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Current mode:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				status.mode = strings.TrimSpace(parts[1])
			}
		}
		if strings.Contains(line, "Loaded policy name:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				status.policy = strings.TrimSpace(parts[1])
			}
		}
	}
	return status
}

// toggleSELinux sets SELinux to enforcing
func (m *Model) toggleSELinux() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("setenforce", "1")
		if out, err := cmd.CombinedOutput(); err != nil {
			return errMsg{fmt.Errorf("failed to set enforcing: %v (output: %s)", err, string(out))}
		}
		return successMsg{"SELinux set to enforcing mode"}
	}
}

// setSELinuxPermissive sets SELinux to permissive
func (m *Model) setSELinuxPermissive() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("setenforce", "0")
		if out, err := cmd.CombinedOutput(); err != nil {
			return errMsg{fmt.Errorf("failed to set permissive: %v (output: %s)", err, string(out))}
		}
		return successMsg{"SELinux set to permissive mode"}
	}
}

// appArmorStatus holds AppArmor state
type appArmorStatus struct {
	loaded bool
	mode   string
}

// getAppArmorStatus checks AppArmor status
func (m *Model) getAppArmorStatus() appArmorStatus {
	cmd := exec.Command("aa-status")
	out, err := cmd.Output()
	if err != nil {
		return appArmorStatus{false, "unknown"}
	}

	status := appArmorStatus{loaded: true}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "apparmor module is loaded") {
			status.loaded = true
		}
		if strings.Contains(line, "profiles are in enforce mode") {
			status.mode = "enforce"
		} else if strings.Contains(line, "profiles are in complain mode") {
			status.mode = "complain"
		}
	}
	return status
}

// toggleAppArmor toggles AppArmor
func (m *Model) toggleAppArmor() tea.Cmd {
	return func() tea.Msg {
		status := m.getAppArmorStatus()

		var cmd *exec.Cmd
		if status.loaded {
			cmd = exec.Command("aa-teardown")
		} else {
			cmd = exec.Command("systemctl", "start", "apparmor")
		}

		if cmd != nil {
			if out, err := cmd.CombinedOutput(); err != nil {
				return errMsg{fmt.Errorf("failed to toggle AppArmor: %v (output: %s)", err, string(out))}
			}
		}

		action := "started"
		if status.loaded {
			action = "stopped"
		}
		return successMsg{fmt.Sprintf("AppArmor %s", action)}
	}
}
