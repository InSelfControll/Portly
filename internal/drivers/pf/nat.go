package pf

import (
	"context"
	"errors"
	"fmt"

	"github.com/orchestrator/unified-firewall/internal/platform"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// ApplyNAT applies a NAT rule using PF
func (d *Driver) ApplyNAT(ctx context.Context, rule models.NATRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid NAT rule: %w", err)
	}

	if !platform.IsRoot() {
		return errors.New("root privileges required for PF")
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

	if err := d.enablePF(ctx); err != nil {
		return fmt.Errorf("failed to enable PF: %w", err)
	}

	if err := d.ensureAnchor(ctx); err != nil {
		return fmt.Errorf("failed to ensure anchor: %w", err)
	}

	if err := d.addRuleToAnchor(ctx, rule); err != nil {
		return fmt.Errorf("failed to add rule: %w", err)
	}

	if err := d.loadAnchor(ctx); err != nil {
		return fmt.Errorf("failed to load anchor: %w", err)
	}

	return nil
}

// RemoveNAT removes a NAT rule
func (d *Driver) RemoveNAT(ctx context.Context, ruleID string) error {
	if !platform.IsRoot() {
		return errors.New("root privileges required for PF")
	}

	rules, err := d.ListNATRules(ctx)
	if err != nil {
		return err
	}

	var found bool
	for _, r := range rules {
		if r.ID == ruleID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("rule not found")
	}

	if err := d.rewriteAnchorWithoutRule(ctx, ruleID); err != nil {
		return fmt.Errorf("failed to update anchor: %w", err)
	}

	if err := d.loadAnchor(ctx); err != nil {
		return fmt.Errorf("failed to reload anchor: %w", err)
	}

	return nil
}
