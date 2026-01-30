package nftables

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

const (
	configDir  = "/etc/nftables.d"
	configFile = "/etc/nftables.conf"
)

func (d *Driver) ensureTableAndChain(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "nft", "list", "table", "inet", tableName)
	_, err := cmd.Output()

	if err != nil {
		cmd = exec.CommandContext(ctx, "nft", "add", "table", "inet", tableName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create table: %w (output: %s)", err, string(output))
		}
	}

	cmd = exec.CommandContext(ctx, "nft", "list", "chain", "inet", tableName, chainName)
	_, err = cmd.Output()

	if err != nil {
		cmd = exec.CommandContext(ctx, "nft", "add", "chain", "inet", tableName, chainName,
			"{ type nat hook prerouting priority dstnat; policy accept; }")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create chain: %w (output: %s)", err, string(output))
		}
	}

	return nil
}

func (d *Driver) findRuleHandle(ctx context.Context, rule models.NATRule) (string, error) {
	cmd := exec.CommandContext(ctx, "nft", "-a", "list", "chain", "inet", tableName, chainName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	proto := strings.ToLower(string(rule.Proto))
	expectedPattern := fmt.Sprintf("%s dport %d", proto, rule.ExternalPort)

	for _, line := range lines {
		if strings.Contains(line, expectedPattern) &&
			strings.Contains(line, fmt.Sprintf("dnat to %s:%d", rule.InternalIP, rule.InternalPort)) {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "handle" && i+1 < len(parts) {
					return parts[i+1], nil
				}
			}
		}
	}

	return "", nil
}

func (d *Driver) saveRules(ctx context.Context) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	orchestratorFile := configDir + "/orchestrator.conf"

	cmd := exec.CommandContext(ctx, "nft", "list", "table", "inet", tableName)
	output, err := cmd.Output()
	if err != nil {
		output = []byte(fmt.Sprintf("table inet %s {\n}\n", tableName))
	}

	if err := os.WriteFile(orchestratorFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return d.ensureIncludeInMainConfig(orchestratorFile)
}

func (d *Driver) ensureIncludeInMainConfig(includePath string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config := fmt.Sprintf("#!/usr/sbin/nft -f\n\nflush ruleset\n\ninclude \"%s\"\n", includePath)
		return os.WriteFile(configFile, []byte(config), 0644)
	}

	content, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	includeLine := fmt.Sprintf(`include "%s"`, includePath)
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == includeLine {
			return nil
		}
	}

	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + includeLine + "\n")
	return err
}

// extractHandle extracts handle number from nft output line
func extractHandle(line string) string {
	parts := strings.Fields(line)
	for i, part := range parts {
		if part == "handle" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
