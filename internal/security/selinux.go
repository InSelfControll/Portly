package security

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// GetSELinuxStatus returns the current SELinux status
func (m *Manager) GetSELinuxStatus() (string, error) {
	if !m.osInfo.IsRHEL() {
		return "N/A", nil
	}

	output, err := exec.Command("sestatus").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get SELinux status: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Current mode:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				return strings.TrimSpace(strings.ToLower(parts[1])), nil
			}
		}
	}

	return "unknown", nil
}

// SetSELinuxBoolean sets an SELinux boolean persistently
func (m *Manager) SetSELinuxBoolean(ctx context.Context, name string, value bool) error {
	if !m.osInfo.IsRHEL() {
		return nil
	}

	valStr := "off"
	if value {
		valStr = "on"
	}

	cmd := exec.CommandContext(ctx, "setsebool", "-P", name, valStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set SELinux boolean %s: %w (output: %s)", name, err, string(output))
	}

	return nil
}

// AddSELinuxPort adds a port label for SELinux
func (m *Manager) AddSELinuxPort(ctx context.Context, port int, proto models.Protocol, selinuxType string) error {
	if !m.osInfo.IsRHEL() {
		return nil
	}

	if selinuxType == "" {
		selinuxType = "http_port_t"
	}

	protoStr := strings.ToLower(string(proto))
	cmd := exec.CommandContext(ctx, "semanage", "port", "-a", "-t", selinuxType, "-p", protoStr, fmt.Sprintf("%d", port))
	output, err := cmd.CombinedOutput()

	if err != nil && !strings.Contains(string(output), "already defined") {
		return fmt.Errorf("failed to add SELinux port: %w (output: %s)", err, string(output))
	}

	return nil
}

// RemoveSELinuxPort removes a port label from SELinux
func (m *Manager) RemoveSELinuxPort(ctx context.Context, port int, proto models.Protocol) error {
	if !m.osInfo.IsRHEL() {
		return nil
	}

	protoStr := strings.ToLower(string(proto))
	cmd := exec.CommandContext(ctx, "semanage", "port", "-d", "-p", protoStr, fmt.Sprintf("%d", port))
	output, err := cmd.CombinedOutput()

	if err != nil && !strings.Contains(string(output), "does not exist") {
		return fmt.Errorf("failed to remove SELinux port: %w (output: %s)", err, string(output))
	}

	return nil
}
