package nftables

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
	handle, err := d.findFilterRuleHandle(ctx, rule)
	if err != nil {
		return err
	}

	if handle == "" {
		return fmt.Errorf("rule not found in firewall")
	}

	cmd := exec.CommandContext(ctx, "nft", "delete", "rule", "inet", filterTableName, filterChainName, "handle", handle)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove rule: %w (output: %s)", err, string(output))
	}

	return d.saveRules(ctx)
}

// findFilterRuleHandle finds the handle for a filter rule
func (d *Driver) findFilterRuleHandle(ctx context.Context, rule models.FirewallRule) (string, error) {
	cmd := exec.CommandContext(ctx, "nft", "-a", "list", "chain", "inet", filterTableName, filterChainName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	proto := strings.ToLower(string(rule.Protocol))

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if rule.Type == models.RuleTypePortLimit {
			expected := fmt.Sprintf("ip saddr %s %s dport %d", rule.SourceIP, proto, rule.Port)
			if strings.Contains(line, expected) {
				return extractHandle(line), nil
			}
		} else {
			expected := fmt.Sprintf("%s dport %d", proto, rule.Port)
			if strings.Contains(line, expected) && !strings.Contains(line, "ip saddr") {
				return extractHandle(line), nil
			}
		}
	}

	return "", nil
}

// ListFirewallRules lists all firewall rules
func (d *Driver) ListFirewallRules(ctx context.Context) ([]models.FirewallRule, error) {
	cmd := exec.CommandContext(ctx, "nft", "-a", "list", "chain", "inet", filterTableName, filterChainName)
	output, err := cmd.Output()
	if err != nil {
		return []models.FirewallRule{}, nil
	}

	var rules []models.FirewallRule
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "dport") {
			continue
		}

		rule := d.parseFilterRule(line)
		if rule != nil {
			rules = append(rules, *rule)
		}
	}

	return rules, nil
}

// parseFilterRule parses an nft filter rule
func (d *Driver) parseFilterRule(line string) *models.FirewallRule {
	parts := strings.Fields(line)
	rule := &models.FirewallRule{}

	for i, part := range parts {
		if part == "handle" && i+1 < len(parts) {
			rule.ID = fmt.Sprintf("nft-filter-%s", parts[i+1])
		}
	}

	for i, part := range parts {
		if part == "saddr" && i+1 < len(parts) {
			rule.Type = models.RuleTypePortLimit
			rule.SourceIP = parts[i+1]
		}
	}

	if rule.Type == "" {
		rule.Type = models.RuleTypePort
	}

	if strings.HasPrefix(line, "tcp") {
		rule.Protocol = models.TCP
	} else if strings.HasPrefix(line, "udp") {
		rule.Protocol = models.UDP
	}

	for i, part := range parts {
		if part == "dport" && i+1 < len(parts) {
			if port, err := strconv.Atoi(parts[i+1]); err == nil {
				rule.Port = port
			}
		}
	}

	if rule.Port == 0 {
		return nil
	}

	return rule
}
