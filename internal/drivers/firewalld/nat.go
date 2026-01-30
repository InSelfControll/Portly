package firewalld

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// ApplyNAT applies a NAT rule using firewalld
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

	if err := d.enableIPForwarding(ctx); err != nil {
		return fmt.Errorf("failed to enable IP forwarding: %w", err)
	}

	if err := d.enableMasquerade(ctx); err != nil {
		return fmt.Errorf("failed to enable masquerade: %w", err)
	}

	richRule := fmt.Sprintf(
		`rule family="ipv4" forward-port port="%d" protocol="%s" to-port="%d" to-addr="%s" accept`,
		rule.ExternalPort,
		strings.ToLower(string(rule.Proto)),
		rule.InternalPort,
		rule.InternalIP,
	)

	cmd := exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--add-rich-rule", richRule)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add NAT rule: %w (output: %s)", err, string(output))
	}

	if err := d.reload(ctx); err != nil {
		return fmt.Errorf("failed to reload firewalld: %w", err)
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

	richRule := fmt.Sprintf(
		`rule family="ipv4" forward-port port="%d" protocol="%s" to-port="%d" to-addr="%s" accept`,
		targetRule.ExternalPort,
		strings.ToLower(string(targetRule.Proto)),
		targetRule.InternalPort,
		targetRule.InternalIP,
	)

	cmd := exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--remove-rich-rule", richRule)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove NAT rule: %w (output: %s)", err, string(output))
	}

	if err := d.reload(ctx); err != nil {
		return fmt.Errorf("failed to reload firewalld: %w", err)
	}

	return nil
}
