package security

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/orchestrator/unified-firewall/internal/platform"
	"text/template"
)

const appArmorProfileTemplate = `# AppArmor profile for {{.Name}}
#include <tunables/global>
{{.BinaryPath}} {
  #include <abstractions/base>
  #include <abstractions/nameservice>
  capability net_bind_service,
  network inet stream,
  network inet6 stream,
  /etc/{{.Name}}/** r,
  /var/lib/{{.Name}}/** rwk,
}
`

// GenerateAppArmorProfile generates an AppArmor profile for a product
func (m *Manager) GenerateAppArmorProfile(product string, ports []int) (string, error) {
	if !m.IsAppArmorAvailable() {
		return "", fmt.Errorf("AppArmor is not available")
	}

	binaryPath, exists := platform.IsProductInstalled(product)
	if !exists {
		binaryPath = fmt.Sprintf("/usr/bin/%s", product)
	}

	tmpl, err := template.New("apparmor").Parse(appArmorProfileTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := AppArmorProfileTemplate{
		Name:        product,
		BinaryPath:  binaryPath,
		AllowPorts:  ports,
		NetworkBind: len(ports) > 0,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// LoadAppArmorProfile loads an AppArmor profile
func (m *Manager) LoadAppArmorProfile(ctx context.Context, profileName, content string) error {
	if !m.IsAppArmorAvailable() {
		return nil
	}

	profilePath := filepath.Join("/etc/apparmor.d", fmt.Sprintf("orchestrator.%s", profileName))

	if err := os.WriteFile(profilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	cmd := exec.CommandContext(ctx, "apparmor_parser", "-r", profilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to load profile: %w (output: %s)", err, string(output))
	}

	return nil
}

// UnloadAppArmorProfile unloads an AppArmor profile
func (m *Manager) UnloadAppArmorProfile(ctx context.Context, profileName string) error {
	if !m.IsAppArmorAvailable() {
		return nil
	}

	profilePath := filepath.Join("/etc/apparmor.d", fmt.Sprintf("orchestrator.%s", profileName))

	cmd := exec.CommandContext(ctx, "apparmor_parser", "-R", profilePath)
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "does not exist") {
		return fmt.Errorf("failed to unload profile: %w (output: %s)", err, string(output))
	}

	os.Remove(profilePath)
	return nil
}
