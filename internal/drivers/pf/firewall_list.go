package pf

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// ClosePort removes a firewall rule
func (d *Driver) ClosePort(ctx context.Context, ruleID string) error {
	portlyAnchorFile := ""
	content, err := os.ReadFile(portlyAnchorFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	skipBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "# ID: ") {
			id := strings.TrimPrefix(line, "# ID: ")
			skipBlock = (id == ruleID)
		}

		if !skipBlock {
			newLines = append(newLines, line)
		}

		if skipBlock && !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
			skipBlock = false
		}
	}

	if err := os.WriteFile(portlyAnchorFile, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to update anchor: %w", err)
	}

	return d.loadPortlyAnchor(ctx)
}

// ListFirewallRules lists all firewall rules
func (d *Driver) ListFirewallRules(ctx context.Context) ([]models.FirewallRule, error) {
	portlyAnchorFile := ""
	content, err := os.ReadFile(portlyAnchorFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.FirewallRule{}, nil
		}
		return nil, err
	}

	return d.parseFilterRules(string(content))
}

// parseFilterRules parses pass rules from anchor file
func (d *Driver) parseFilterRules(content string) ([]models.FirewallRule, error) {
	var rules []models.FirewallRule
	var currentRule *models.FirewallRule

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "# ID: ") {
			if currentRule != nil {
				rules = append(rules, *currentRule)
			}
			currentRule = &models.FirewallRule{
				ID: strings.TrimPrefix(line, "# ID: "),
			}
			continue
		}

		if currentRule == nil {
			continue
		}

		if strings.HasPrefix(line, "# Type: ") {
			typeStr := strings.TrimPrefix(line, "# Type: ")
			currentRule.Type = models.FirewallRuleType(typeStr)
			continue
		}

		if strings.HasPrefix(line, "pass ") {
			d.parsePassRule(line, currentRule)
			rules = append(rules, *currentRule)
			currentRule = nil
		}
	}

	if currentRule != nil {
		rules = append(rules, *currentRule)
	}

	return rules, nil
}

// parsePassRule parses a PF pass rule
func (d *Driver) parsePassRule(line string, rule *models.FirewallRule) {
	parts := strings.Fields(line)

	for i, part := range parts {
		switch part {
		case "proto":
			if i+1 < len(parts) {
				proto := strings.ToUpper(parts[i+1])
				if proto == "TCP" {
					rule.Protocol = models.TCP
				} else if proto == "UDP" {
					rule.Protocol = models.UDP
				}
			}
		case "from":
			if i+1 < len(parts) {
				nextPart := parts[i+1]
				if nextPart != "any" {
					rule.SourceIP = nextPart
					rule.Type = models.RuleTypePortLimit
				}
			}
		case "port":
			if i+1 < len(parts) {
				if port, err := strconv.Atoi(parts[i+1]); err == nil {
					rule.Port = port
				}
			}
		}
	}
}
