package firewalld

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// ListNATRules returns all applied NAT rules
func (d *Driver) ListNATRules(ctx context.Context) ([]models.NATRule, error) {
	cmd := exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--list-rich-rules")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list NAT rules: %w", err)
	}

	var rules []models.NATRule
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "forward-port") {
			continue
		}

		rule, err := d.parseRichRule(line)
		if err == nil && rule != nil {
			rules = append(rules, *rule)
		}
	}

	return rules, nil
}

// parseRichRule parses a firewalld rich rule string into a NATRule
func (d *Driver) parseRichRule(ruleStr string) (*models.NATRule, error) {
	rule := &models.NATRule{
		ID: fmt.Sprintf("fw-%d", time.Now().UnixNano()),
	}

	if port := extractValue(ruleStr, `port="`); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			rule.ExternalPort = p
		}
	}

	if proto := extractValue(ruleStr, `protocol="`); proto != "" {
		rule.Proto = models.Protocol(strings.ToUpper(proto))
	}

	if toPort := extractValue(ruleStr, `to-port="`); toPort != "" {
		if p, err := strconv.Atoi(toPort); err == nil {
			rule.InternalPort = p
		}
	}

	if toAddr := extractValue(ruleStr, `to-addr="`); toAddr != "" {
		rule.InternalIP = toAddr
	}

	if rule.ExternalPort == 0 || rule.InternalPort == 0 || rule.InternalIP == "" {
		return nil, fmt.Errorf("could not parse rule")
	}

	return rule, nil
}

// CheckConflicts checks for port conflicts
func (d *Driver) CheckConflicts(ctx context.Context, port int, proto models.Protocol) error {
	rules, err := d.ListNATRules(ctx)
	if err != nil {
		return err
	}

	for _, r := range rules {
		if r.ExternalPort == port && r.Proto == proto {
			return fmt.Errorf("port %d/%s already in use", port, proto)
		}
	}

	protoStr := strings.ToLower(string(proto))
	cmd := exec.CommandContext(ctx, "firewall-cmd", "--query-port", fmt.Sprintf("%d/%s", port, protoStr))
	output, err := cmd.CombinedOutput()

	if err == nil && strings.TrimSpace(string(output)) == "yes" {
		return fmt.Errorf("port %d/%s already open in firewalld", port, proto)
	}

	return nil
}

// extractValue extracts a value from a string with format key="value"
func extractValue(s, key string) string {
	idx := strings.Index(s, key)
	if idx == -1 {
		return ""
	}
	start := idx + len(key)
	end := strings.Index(s[start:], `"`)
	if end == -1 {
		return ""
	}
	return s[start : start+end]
}
