package pf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

func (d *Driver) enablePF(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, pfctlPath, "-s", "info")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "Enabled") {
		return nil
	}

	cmd = exec.CommandContext(ctx, pfctlPath, "-e")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable PF: %w (output: %s)", err, string(output))
	}

	return nil
}

func (d *Driver) ensureAnchor(ctx context.Context) error {
	if err := os.MkdirAll(anchorDir, 0755); err != nil {
		return fmt.Errorf("failed to create anchor dir: %w", err)
	}

	if _, err := os.Stat(anchorFile); os.IsNotExist(err) {
		header := fmt.Sprintf("# Orchestrator NAT Rules\n# Anchor: %s\n\n", anchorName)
		if err := os.WriteFile(anchorFile, []byte(header), 0644); err != nil {
			return fmt.Errorf("failed to create anchor file: %w", err)
		}
	}

	return d.ensureAnchorInPfConf(ctx)
}

func (d *Driver) ensureAnchorInPfConf(ctx context.Context) error {
	if _, err := os.Stat(pfConfFile); os.IsNotExist(err) {
		config := fmt.Sprintf(`# PF configuration\nanchor "%s"\nload anchor "%s" from "%s"\n`, anchorName, anchorName, anchorFile)
		return os.WriteFile(pfConfFile, []byte(config), 0644)
	}

	content, err := os.ReadFile(pfConfFile)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), anchorName) {
		return nil
	}

	newContent := string(content) + fmt.Sprintf(`\nanchor "%s"\nload anchor "%s" from "%s"\n`, anchorName, anchorName, anchorFile)
	return os.WriteFile(pfConfFile, []byte(newContent), 0644)
}

func (d *Driver) addRuleToAnchor(ctx context.Context, rule models.NATRule) error {
	content, err := os.ReadFile(anchorFile)
	if err != nil {
		return err
	}

	pfRule := d.buildPFRule(rule)
	newContent := string(content) + pfRule + "\n"

	return os.WriteFile(anchorFile, []byte(newContent), 0644)
}

func (d *Driver) buildPFRule(rule models.NATRule) string {
	proto := strings.ToLower(string(rule.Proto))

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# ID: %s\n", rule.ID))
	sb.WriteString(fmt.Sprintf("# Product: %s\n", rule.Product))
	if rule.Description != "" {
		sb.WriteString(fmt.Sprintf("# Description: %s\n", rule.Description))
	}
	sb.WriteString(fmt.Sprintf("rdr pass on any inet proto %s from any to any port %d -> %s port %d",
		proto, rule.ExternalPort, rule.InternalIP, rule.InternalPort))

	return sb.String()
}

func (d *Driver) loadAnchor(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, pfctlPath, "-a", anchorName, "-f", anchorFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pfctl failed: %w (output: %s)", err, string(output))
	}
	return nil
}

func (d *Driver) rewriteAnchorWithoutRule(ctx context.Context, ruleID string) error {
	content, err := os.ReadFile(anchorFile)
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

	return os.WriteFile(anchorFile, []byte(strings.Join(newLines, "\n")), 0644)
}
