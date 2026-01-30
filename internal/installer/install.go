package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// InstallProduct installs a product using the appropriate package manager
func (i *Installer) InstallProduct(ctx context.Context, config ProductConfig) error {
	pm := i.osInfo.PackageManager()

	switch pm {
	case "dnf":
		return i.installWithDNF(ctx, config)
	case "yum":
		return i.installWithYUM(ctx, config)
	case "apt":
		return i.installWithAPT(ctx, config)
	case "brew":
		return i.installWithBrew(ctx, config)
	default:
		return fmt.Errorf("unsupported package manager: %s", pm)
	}
}

func (i *Installer) runPostInstall(ctx context.Context, config ProductConfig) error {
	for _, cmdStr := range config.PostInstall {
		fmt.Printf("Running post-install: %s\n", cmdStr)

		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			continue
		}

		cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: command failed: %v\n", err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
