package nftables

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// ApplyNAT applies a NAT rule using nftables
func (d *Driver) ApplyNAT(ctx context.Context, rule models.NATRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid NAT rule: %w", err)
	}

	rules, err := d.ListNATRules(ctx)
	if err != nil {
		return err
	}

	for _, r := range rules {
		if r.ExternalPort == rule.ExternalPort && r.Proto == rule.Proto {
			return fmt.Errorf("port %d/%s already mapped: %w", rule.ExternalPort, rule.Proto, errors.New("rule exists"))
		}
	}

	if err := d.ensureTableAndChain(ctx); err != nil {
		return fmt.Errorf("failed to ensure table structure: %w", err)
	}

	ruleStr := d.buildNATRule(rule)

	cmd := exec.CommandContext(ctx, "nft", "add", "rule", "inet", tableName, chainName, ruleStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add NAT rule: %w (output: %s)", err, string(output))
	}

	if err := d.saveRules(ctx); err != nil {
		return fmt.Errorf("failed to save rules: %w", err)
	}

	return nil
}

// RemoveNAT removes a NAT rule
func (d *Driver) RemoveNAT(ctx context.Context, ruleID string) error {
	rules, err := d.ListNATRules(ctx)
	if err != nil {
		return err
	}

	var targetRule *models.NATRule
	for i, r := range rules {
		if r.ID == ruleID {
			targetRule = &rules[i]
			break
		}
	}

	if targetRule == nil {
		return errors.New("rule not found")
	}

	handle, err := d.findRuleHandle(ctx, *targetRule)
	if err != nil {
		return fmt.Errorf("failed to find rule handle: %w", err)
	}

	if handle == "" {
		return fmt.Errorf("could not find rule to remove")
	}

	cmd := exec.CommandContext(ctx, "nft", "delete", "rule", "inet", tableName, chainName, "handle", handle)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove NAT rule: %w (output: %s)", err, string(output))
	}

	if err := d.saveRules(ctx); err != nil {
		return fmt.Errorf("failed to save rules: %w", err)
	}

	return nil
}

func (d *Driver) buildNATRule(rule models.NATRule) string {
	proto := strings.ToLower(string(rule.Proto))
	return fmt.Sprintf("%s dport %d dnat to %s:%d",
		proto, rule.ExternalPort, rule.InternalIP, rule.InternalPort)
}
