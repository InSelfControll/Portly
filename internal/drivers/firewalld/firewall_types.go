package firewalld

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// OpenPort opens a port in the firewall
func (d *Driver) OpenPort(ctx context.Context, rule models.FirewallRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid firewall rule: %w", err)
	}

	proto := strings.ToLower(string(rule.Protocol))
	portStr := fmt.Sprintf("%d/%s", rule.Port, proto)

	cmd := exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--add-port", portStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "already") {
			return nil
		}
		return fmt.Errorf("failed to open port: %w (output: %s)", err, string(output))
	}

	return d.reload(ctx)
}

// OpenPortForIP opens a port limited to a specific source IP
func (d *Driver) OpenPortForIP(ctx context.Context, rule models.FirewallRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid firewall rule: %w", err)
	}

	richRule := fmt.Sprintf(
		`rule family="ipv4" source address="%s" port protocol="%s" port="%d" accept`,
		rule.SourceIP,
		strings.ToLower(string(rule.Protocol)),
		rule.Port,
	)

	cmd := exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--add-rich-rule", richRule)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "already") {
			return nil
		}
		return fmt.Errorf("failed to add IP-limited rule: %w (output: %s)", err, string(output))
	}

	return d.reload(ctx)
}

// TrustIP opens all ports for a specific source IP
func (d *Driver) TrustIP(ctx context.Context, rule models.FirewallRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid firewall rule: %w", err)
	}

	richRule := fmt.Sprintf(
		`rule family="ipv4" source address="%s" accept`,
		rule.SourceIP,
	)

	cmd := exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--add-rich-rule", richRule)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "already") {
			return nil
		}
		return fmt.Errorf("failed to trust IP: %w (output: %s)", err, string(output))
	}

	return d.reload(ctx)
}
