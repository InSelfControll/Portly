package security

import (
	"fmt"
	"os"

	"github.com/orchestrator/unified-firewall/internal/platform"
)

// NewManager creates a new security manager
func NewManager() (*Manager, error) {
	osInfo, err := platform.DetectOS()
	if err != nil {
		return nil, fmt.Errorf("failed to detect OS: %w", err)
	}
	return &Manager{osInfo: osInfo}, nil
}

// IsSELinuxEnforcing returns true if SELinux is in enforcing mode
func (m *Manager) IsSELinuxEnforcing() bool {
	status, err := m.GetSELinuxStatus()
	if err != nil {
		return false
	}
	return status == "enforcing"
}

// IsAppArmorAvailable returns true if AppArmor is available
func (m *Manager) IsAppArmorAvailable() bool {
	if !m.osInfo.IsDebian() {
		return false
	}
	return fileExists("/sys/kernel/security/apparmor")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
