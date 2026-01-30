package installer

import (
	"github.com/orchestrator/unified-firewall/internal/platform"
)

// Installer handles product installation
type Installer struct {
	osInfo     *platform.OSInfo
	autoAccept bool
}

// ProductConfig contains installation configuration for a product
type ProductConfig struct {
	Name         string
	DisplayName  string
	PackageName  string
	BrewFormula  string
	Repositories []string
	Dependencies []string
	PostInstall  []string
}
