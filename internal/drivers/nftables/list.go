package nftables

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// ListNATRules returns all applied NAT rules
func (d *Driver) ListNATRules(ctx context.Context) ([]models.NATRule, error) {
	cmd := exec.CommandContext(ctx, "nft", "-a", "list", "chain", "inet", tableName, chainName)
	output, err := cmd.Output()
	if err != nil {
		if strings.Contains(string(output), "does not exist") {
			return []models.NATRule{}, nil
		}
		return nil, fmt.Errorf("failed to list NAT rules: %w", err)
	}

	return d.parseRules(string(output))
}

// parseRules parses nft output into NATRule structs
func (d *Driver) parseRules(output string) ([]models.NATRule, error) {
	var rules []models.NATRule
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "dnat to") {
			continue
		}

		rule := d.parseNFTRule(line)
		if rule != nil {
			rules = append(rules, *rule)
		}
	}

	return rules, nil
}

func (d *Driver) parseNFTRule(line string) *models.NATRule {
	parts := strings.Fields(line)
	rule := &models.NATRule{}

	for i, part := range parts {
		if part == "handle" && i+1 < len(parts) {
			rule.ID = fmt.Sprintf("nft-%s", parts[i+1])
		}
	}

	if strings.HasPrefix(line, "tcp") {
		rule.Proto = models.TCP
	} else if strings.HasPrefix(line, "udp") {
		rule.Proto = models.UDP
	}

	for i, part := range parts {
		if part == "dport" && i+1 < len(parts) {
			if port, err := strconv.Atoi(parts[i+1]); err == nil {
				rule.ExternalPort = port
			}
		}
		if part == "to" && i+1 < len(parts) {
			target := parts[i+1]
			if idx := strings.Index(target, ":"); idx != -1 {
				rule.InternalIP = target[:idx]
				if port, err := strconv.Atoi(target[idx+1:]); err == nil {
					rule.InternalPort = port
				}
			}
		}
	}

	if rule.ExternalPort == 0 || rule.InternalIP == "" {
		return nil
	}

	return rule
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

	return nil
}
