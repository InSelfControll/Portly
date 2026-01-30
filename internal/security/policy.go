package security

import (
	"context"
	"fmt"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// ApplySecurityPolicy applies the appropriate security policy based on OS
func (m *Manager) ApplySecurityPolicy(ctx context.Context, product string, policy models.SecurityPolicy, ports []int) error {
	if m.osInfo.IsRHEL() {
		for _, boolean := range policy.SelinuxBooleans {
			if err := m.SetSELinuxBoolean(ctx, boolean, true); err != nil {
				return fmt.Errorf("failed to set boolean %s: %w", boolean, err)
			}
		}

		for _, port := range ports {
			if err := m.AddSELinuxPort(ctx, port, models.TCP, ""); err != nil {
				return fmt.Errorf("failed to add SELinux port: %w", err)
			}
		}
	}

	if m.osInfo.IsDebian() && policy.AppArmorProfile != "" {
		profile, err := m.GenerateAppArmorProfile(product, ports)
		if err != nil {
			return fmt.Errorf("failed to generate profile: %w", err)
		}

		if err := m.LoadAppArmorProfile(ctx, product, profile); err != nil {
			return fmt.Errorf("failed to load profile: %w", err)
		}
	}

	return nil
}

// RemoveSecurityPolicy removes security policies for a product
func (m *Manager) RemoveSecurityPolicy(ctx context.Context, product string, policy models.SecurityPolicy, ports []int) error {
	if m.osInfo.IsRHEL() {
		for _, port := range ports {
			m.RemoveSELinuxPort(ctx, port, models.TCP)
			m.RemoveSELinuxPort(ctx, port, models.UDP)
		}
	}

	if m.osInfo.IsDebian() {
		m.UnloadAppArmorProfile(ctx, product)
	}

	return nil
}
