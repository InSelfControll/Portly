package pf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// loadPortlyAnchor loads the portly PF anchor
func (d *Driver) loadPortlyAnchor(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, pfctlPath, "-a", "com.portly", "-f", portlyAnchorFile)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// OpenPort opens a port in PF
func (d *Driver) OpenPort(ctx context.Context, rule models.FirewallRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid firewall rule: %w", err)
	}

	proto := strings.ToLower(string(rule.Protocol))
	ruleStr := fmt.Sprintf("# ID: %s\n# Type: port\npass in inet proto %s to any port %d\n",
		rule.ID, proto, rule.Port)

	return d.appendToAnchor(ctx, ruleStr)
}

// OpenPortForIP opens a port limited to a specific source IP
func (d *Driver) OpenPortForIP(ctx context.Context, rule models.FirewallRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid firewall rule: %w", err)
	}

	proto := strings.ToLower(string(rule.Protocol))
	ruleStr := fmt.Sprintf("# ID: %s\n# Type: port_limit\npass in inet proto %s from %s to any port %d\n",
		rule.ID, proto, rule.SourceIP, rule.Port)

	return d.appendToAnchor(ctx, ruleStr)
}

// appendToAnchor appends a rule to the PF anchor file
func (d *Driver) appendToAnchor(ctx context.Context, ruleStr string) error {
	content, err := os.ReadFile(portlyAnchorFile)
	if err != nil {
		return err
	}

	newContent := string(content) + ruleStr
	if err := os.WriteFile(portlyAnchorFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write anchor: %w", err)
	}

	return d.loadPortlyAnchor(ctx)
}
