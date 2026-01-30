package pf

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// ListNATRules returns all applied NAT rules
func (d *Driver) ListNATRules(ctx context.Context) ([]models.NATRule, error) {
	content, err := os.ReadFile(anchorFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.NATRule{}, nil
		}
		return nil, fmt.Errorf("failed to read anchor: %w", err)
	}

	return d.parseAnchorFile(string(content))
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

func (d *Driver) parseAnchorFile(content string) ([]models.NATRule, error) {
	var rules []models.NATRule
	var currentRule *models.NATRule

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "# ID: ") {
			if currentRule != nil {
				rules = append(rules, *currentRule)
			}
			currentRule = &models.NATRule{
				ID: strings.TrimPrefix(line, "# ID: "),
			}
			continue
		}

		if currentRule == nil {
			continue
		}

		if strings.HasPrefix(line, "# Product: ") {
			currentRule.Product = strings.TrimPrefix(line, "# Product: ")
			continue
		}

		if strings.HasPrefix(line, "# Description: ") {
			currentRule.Description = strings.TrimPrefix(line, "# Description: ")
			continue
		}

		if strings.HasPrefix(line, "rdr ") {
			d.parsePFRuleLine(line, currentRule)
			rules = append(rules, *currentRule)
			currentRule = nil
		}
	}

	if currentRule != nil {
		rules = append(rules, *currentRule)
	}

	return rules, scanner.Err()
}

func (d *Driver) parsePFRuleLine(line string, rule *models.NATRule) {
	parts := strings.Fields(line)

	for i, part := range parts {
		switch part {
		case "proto":
			if i+1 < len(parts) {
				proto := strings.ToUpper(parts[i+1])
				if proto == "TCP" {
					rule.Proto = models.TCP
				} else if proto == "UDP" {
					rule.Proto = models.UDP
				}
			}
		case "port":
			if i+1 < len(parts) {
				if port, err := strconv.Atoi(parts[i+1]); err == nil {
					if rule.ExternalPort == 0 {
						rule.ExternalPort = port
					}
				}
			}
		case "->":
			if i+1 < len(parts) {
				target := parts[i+1]
				if idx := strings.Index(target, ":"); idx != -1 {
					rule.InternalIP = target[:idx]
					if port, err := strconv.Atoi(target[idx+1:]); err == nil {
						rule.InternalPort = port
					}
				}
			}
		}
	}
}
