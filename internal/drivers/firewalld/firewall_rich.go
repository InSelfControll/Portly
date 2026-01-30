package firewalld

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// parseRichRules parses rich rules to find IP-limited rules (both port-limited and trust IP)
func (d *Driver) parseRichRules(ctx context.Context) ([]models.FirewallRule, error) {
	var allRules []models.FirewallRule
	seen := make(map[string]bool)

	// Parse both runtime and permanent rich rules
	for _, permanent := range []bool{false, true} {
		var cmd *exec.Cmd
		if permanent {
			cmd = exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--list-rich-rules")
		} else {
			cmd = exec.CommandContext(ctx, "firewall-cmd", "--list-rich-rules")
		}

		output, err := cmd.Output()
		if err != nil {
			continue
		}

		lines := strings.Split(string(output), "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || !strings.Contains(line, "source address") {
				continue
			}

			var rule models.FirewallRule
			rule.Protocol = models.TCP // default

			// Extract source IP
			if idx := strings.Index(line, `source address="`); idx != -1 {
				start := idx + len(`source address="`)
				end := strings.Index(line[start:], `"`)
				if end != -1 {
					rule.SourceIP = line[start : start+end]
				}
			}

			if rule.SourceIP == "" {
				continue
			}

			// Check if this is a port-limited rule or a trust IP rule
			hasPort := strings.Contains(line, `port protocol="`) || strings.Contains(line, `port="`)

			if hasPort {
				// Port-limited rule
				rule.Type = models.RuleTypePortLimit

				if idx := strings.Index(line, `port protocol="`); idx != -1 {
					start := idx + len(`port protocol="`)
					end := strings.Index(line[start:], `"`)
					if end != -1 {
						proto := line[start : start+end]
						if proto == "udp" {
							rule.Protocol = models.UDP
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

				if rule.Port > 0 {
					rule.ID = fmt.Sprintf("fw-limit-%s-%d-%s", rule.SourceIP, rule.Port, rule.Protocol)
				}
			} else {
				// Trust IP rule (no port specified, allows all traffic from IP)
				rule.Type = models.RuleTypeTrustIP
				rule.Port = 0
				rule.ID = fmt.Sprintf("fw-trust-%s", rule.SourceIP)
			}

			if rule.ID != "" {
				key := rule.ID
				if !seen[key] {
					seen[key] = true
					allRules = append(allRules, rule)
				}
			}
		}
	}

	return allRules, nil
}
