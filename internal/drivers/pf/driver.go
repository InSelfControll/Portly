package pf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/internal/platform"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

const (
	driverName       = "pf"
	anchorName       = "com.orchestrator.nat"
	anchorDir        = "/etc/pf.anchors"
	anchorFile       = "/etc/pf.anchors/com.orchestrator.nat"
	portlyAnchorFile = "/etc/pf.anchors/com.portly.rules"
	pfConfFile       = "/etc/pf.conf"
	pfctlPath        = "/sbin/pfctl"
)

// Driver implements the Provider interface for macOS PF
type Driver struct {
	osInfo *platform.OSInfo
}

// New creates a new PF driver
func New() *Driver {
	return &Driver{}
}

// Name returns the provider name
func (d *Driver) Name() string {
	return driverName
}

// IsAvailable returns true if PF is available (macOS)
func (d *Driver) IsAvailable() bool {
	osInfo, err := platform.DetectOS()
	if err != nil {
		return false
	}
	d.osInfo = osInfo

	if !osInfo.IsDarwin() {
		return false
	}

	_, err = os.Stat(pfctlPath)
	return err == nil
}

// IsProductInstalled checks if a product is installed
func (d *Driver) IsProductInstalled(ctx context.Context, name string) (models.ProductInfo, error) {
	path, exists := platform.IsProductInstalled(name)

	info := models.ProductInfo{
		Name:        name,
		Path:        path,
		IsInstalled: exists,
	}

	if exists {
		cmd := exec.CommandContext(ctx, name, "--version")
		output, err := cmd.Output()
		if err == nil {
			info.Version = strings.TrimSpace(string(output))
		}
	}

	return info, nil
}

// GetInstalledProducts returns a list of supported products
func (d *Driver) GetInstalledProducts(ctx context.Context) ([]models.ProductInfo, error) {
	products := []string{"podman", "docker", "tailscale", "twingate"}

	var result []models.ProductInfo
	for _, p := range products {
		info, err := d.IsProductInstalled(ctx, p)
		if err == nil {
			result = append(result, info)
		}
	}

	return result, nil
}

// EnsureSecurityPolicy applies security policies
func (d *Driver) EnsureSecurityPolicy(ctx context.Context, product string, policy models.SecurityPolicy) error {
	return nil
}

// RemoveSecurityPolicy removes security policies
func (d *Driver) RemoveSecurityPolicy(ctx context.Context, product string) error {
	return nil
}

// TrustIP opens all ports for a specific source IP
func (d *Driver) TrustIP(ctx context.Context, rule models.FirewallRule) error {
	return fmt.Errorf("TrustIP not implemented for PF")
}
