package firewalld

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// ClosePort removes a firewall rule
func (d *Driver) ClosePort(ctx context.Context, ruleID string) error {
	rules, err := d.ListFirewallRules(ctx)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.ID == ruleID {
			return d.removeFirewallRule(ctx, rule)
		}
	}

	return fmt.Errorf("rule not found: %s", ruleID)
}

// removeFirewallRule removes a specific firewall rule
func (d *Driver) removeFirewallRule(ctx context.Context, rule models.FirewallRule) error {
	proto := strings.ToLower(string(rule.Protocol))

	switch rule.Type {
	case models.RuleTypePort:
		portStr := fmt.Sprintf("%d/%s", rule.Port, proto)
		cmd := exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--remove-port", portStr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to close port: %w (output: %s)", err, string(output))
		}

	case models.RuleTypePortLimit:
		richRule := fmt.Sprintf(
			`rule family="ipv4" source address="%s" port protocol="%s" port="%d" accept`,
			rule.SourceIP, proto, rule.Port,
		)
		cmd := exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--remove-rich-rule", richRule)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to remove IP-limited rule: %w (output: %s)", err, string(output))
		}
	}

	return d.reload(ctx)
}

// ListFirewallRules lists all firewall rules
func (d *Driver) ListFirewallRules(ctx context.Context) ([]models.FirewallRule, error) {
	var rules []models.FirewallRule
	seen := make(map[string]bool)

	// Try to list runtime ports first
	cmd := exec.CommandContext(ctx, "firewall-cmd", "--list-ports")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.Split(line, "/")
			if len(parts) != 2 {
				continue
			}

			port, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}

			proto := models.TCP
			if parts[1] == "udp" {
				proto = models.UDP
			}

			key := fmt.Sprintf("%d-%s", port, proto)
			if !seen[key] {
				seen[key] = true
				rules = append(rules, models.FirewallRule{
					ID:       fmt.Sprintf("fw-port-%d-%s", port, proto),
					Type:     models.RuleTypePort,
					Port:     port,
					Protocol: proto,
				})
			}
		}
	}

	// Also check permanent config
	cmd = exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--list-ports")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.Split(line, "/")
			if len(parts) != 2 {
				continue
			}

			port, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}

			proto := models.TCP
			if parts[1] == "udp" {
				proto = models.UDP
			}

			key := fmt.Sprintf("%d-%s", port, proto)
			if !seen[key] {
				seen[key] = true
				rules = append(rules, models.FirewallRule{
					ID:       fmt.Sprintf("fw-port-%d-%s", port, proto),
					Type:     models.RuleTypePort,
					Port:     port,
					Protocol: proto,
				})
			}
		}
	}

	// Parse rich rules
	richRules, err := d.parseRichRules(ctx)
	if err == nil {
		for _, r := range richRules {
			key := r.ID
			if !seen[key] {
				seen[key] = true
				rules = append(rules, r)
			}
		}
	}

	return rules, nil
}
