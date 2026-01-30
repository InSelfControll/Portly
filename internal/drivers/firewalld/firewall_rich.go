package firewalld

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// parseRichRules parses rich rules to find IP-limited port rules
func (d *Driver) parseRichRules(ctx context.Context) ([]models.FirewallRule, error) {
	cmd := exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--list-rich-rules")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var rules []models.FirewallRule
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "source address") {
			continue
		}

		var rule models.FirewallRule
		rule.Type = models.RuleTypePortLimit

		if idx := strings.Index(line, `source address="`); idx != -1 {
			start := idx + len(`source address="`)
			end := strings.Index(line[start:], `"`)
			if end != -1 {
				rule.SourceIP = line[start : start+end]
			}
		}

		if idx := strings.Index(line, `port protocol="`); idx != -1 {
			start := idx + len(`port protocol="`)
			end := strings.Index(line[start:], `"`)
			if end != -1 {
				proto := line[start : start+end]
				if proto == "udp" {
					rule.Protocol = models.UDP
				} else {
					rule.Protocol = models.TCP
				}
			}
		}

		if idx := strings.Index(line, `port="`); idx != -1 {
			start := idx + len(`port="`)
			end := strings.Index(line[start:], `"`)
			if end != -1 {
				port, _ := strconv.Atoi(line[start : start+end])
				rule.Port = port
			}
		}

		if rule.Port > 0 && rule.SourceIP != "" {
			rule.ID = fmt.Sprintf("fw-limit-%s-%d-%s", rule.SourceIP, rule.Port, rule.Protocol)
			rules = append(rules, rule)
		}
	}

	return rules, nil
}
