package security

import (
	"github.com/orchestrator/unified-firewall/internal/platform"
)

// Manager handles security policy enforcement
type Manager struct {
	osInfo *platform.OSInfo
}

// AppArmorProfileTemplate is the template for generating AppArmor profiles
type AppArmorProfileTemplate struct {
	Name        string
	BinaryPath  string
	AllowPorts  []int
	NetworkBind bool
}
