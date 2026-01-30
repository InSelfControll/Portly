package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func (i *Installer) installWithAPT(ctx context.Context, config ProductConfig) error {
	fmt.Println("Updating package list...")
	cmd := exec.CommandContext(ctx, "apt-get", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: apt-get update failed: %v\n", err)
	}

	for _, dep := range config.Dependencies {
		cmd := exec.CommandContext(ctx, "apt-get", "install", "-y", dep)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: failed to install %s: %v\n", dep, err)
		}
	}

	cmd = exec.CommandContext(ctx, "apt-get", "install", "-y", config.PackageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return i.runPostInstall(ctx, config)
}
