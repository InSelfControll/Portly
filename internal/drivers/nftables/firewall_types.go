package nftables

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

const (
	filterTableName = "orchestrator_filter"
	filterChainName = "input"
)

// OpenPort opens a port in the firewall
func (d *Driver) OpenPort(ctx context.Context, rule models.FirewallRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid firewall rule: %w", err)
	}

	if err := d.ensureFilterTable(ctx); err != nil {
		return err
	}

	proto := strings.ToLower(string(rule.Protocol))
	ruleStr := fmt.Sprintf("%s dport %d accept", proto, rule.Port)

	cmd := exec.CommandContext(ctx, "nft", "add", "rule", "inet", filterTableName, filterChainName, ruleStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to open port: %w (output: %s)", err, string(output))
	}

	return d.saveRules(ctx)
}

// OpenPortForIP opens a port limited to a specific source IP
func (d *Driver) OpenPortForIP(ctx context.Context, rule models.FirewallRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid firewall rule: %w", err)
	}

	if err := d.ensureFilterTable(ctx); err != nil {
		return err
	}

	proto := strings.ToLower(string(rule.Protocol))
	ruleStr := fmt.Sprintf("ip saddr %s %s dport %d accept", rule.SourceIP, proto, rule.Port)

	cmd := exec.CommandContext(ctx, "nft", "add", "rule", "inet", filterTableName, filterChainName, ruleStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add IP-limited rule: %w (output: %s)", err, string(output))
	}

	return d.saveRules(ctx)
}

// ensureFilterTable ensures the filter table and chain exist
func (d *Driver) ensureFilterTable(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "nft", "list", "table", "inet", filterTableName)
	_, err := cmd.Output()

	if err != nil {
		cmd = exec.CommandContext(ctx, "nft", "add", "table", "inet", filterTableName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create filter table: %w (output: %s)", err, string(output))
		}
	}

	cmd = exec.CommandContext(ctx, "nft", "list", "chain", "inet", filterTableName, filterChainName)
	_, err = cmd.Output()

	if err != nil {
		cmd = exec.CommandContext(ctx, "nft", "add", "chain", "inet", filterTableName, filterChainName,
			"{ type filter hook input priority 0; policy accept; }")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create filter chain: %w (output: %s)", err, string(output))
		}
	}

	return nil
}
