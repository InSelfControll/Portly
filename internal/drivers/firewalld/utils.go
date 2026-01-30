package firewalld

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func (d *Driver) enableIPForwarding(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "sysctl", "-n", "net.ipv4.ip_forward")
	output, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(output)) == "1" {
		return nil
	}

	cmd = exec.CommandContext(ctx, "sysctl", "-w", "net.ipv4.ip_forward=1")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable IP forwarding: %w (output: %s)", err, string(output))
	}

	return nil
}

func (d *Driver) enableMasquerade(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "firewall-cmd", "--get-default-zone")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get default zone: %w", err)
	}
	zone := strings.TrimSpace(string(output))

	cmd = exec.CommandContext(ctx, "firewall-cmd", "--zone", zone, "--query-masquerade")
	_, err = cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	cmd = exec.CommandContext(ctx, "firewall-cmd", "--permanent", "--zone", zone, "--add-masquerade")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable masquerade: %w (output: %s)", err, string(output))
	}

	return nil
}

func (d *Driver) reload(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "firewall-cmd", "--reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("firewall-cmd --reload failed: %w (output: %s)", err, string(output))
	}
	return nil
}
